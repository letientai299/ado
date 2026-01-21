// Package api provides a dynamic CLI for Azure DevOps REST APIs.
//
// The api command uses reflection to discover API methods from the rest package
// at runtime, enabling tab completion and invocation without manual registration.
// When new API clients are added to internal/rest, they automatically become
// available through this command.
//
// # Architecture
//
// The package uses a two-level hierarchy mirroring the REST client structure:
//
//	Client
//	├── Git()          -> Git client
//	│   ├── PRs()      -> GitPRs (scoped to repository)
//	│   │   ├── List()
//	│   │   ├── ByID()
//	│   │   └── ...
//	│   └── RepoInfo()
//	├── Builds()       -> Builds client
//	│   └── ForProject() -> ProjectBuilds (scoped)
//	│       ├── List()
//	│       └── ...
//	└── ...
//
// The [Registry] discovers this structure via reflection and builds an endpoint
// tree that maps CLI paths (like "git.prs.list") to method invocations.
package api

import (
	"context"
	"reflect"
	"sort"
	"strings"
	"sync"
	"unicode"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/rest"
)

// Registry holds discovered API endpoints from the rest.Client.
// It uses reflection to discover methods and build a tree of endpoints
// that can be invoked via CLI paths.
//
// The registry is designed to be initialized once and reused. It caches
// the reflection results to avoid repeated introspection.
type Registry struct {
	mu        sync.RWMutex
	endpoints map[string]*Endpoint // path -> endpoint, e.g. "git.prs.list"
	tree      *node                // tree structure for completion
	built     bool
}

// node represents a node in the endpoint tree for tab completion.
// Each node can have children (for nested paths) and/or an endpoint
// (if this path is a valid API call).
type node struct {
	name     string
	children map[string]*node
	endpoint *Endpoint
}

// NewRegistry creates a new empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		endpoints: make(map[string]*Endpoint),
		tree: &node{
			children: make(map[string]*node),
		},
	}
}

// Build discovers all API endpoints from the rest.Client using reflection.
// This method is idempotent - calling it multiple times has no effect after
// the first call.
//
// The discovery process:
//  1. Get all exported methods on rest.Client that return API group clients
//     (e.g., Git(), Builds(), Pipelines())
//  2. For each API group, get methods that return scoped clients
//     (e.g., Git.PRs() returns GitPRs)
//  3. For each scoped client, discover the actual API methods
//     (e.g., GitPRs.List(), GitPRs.ByID())
//  4. Build endpoint metadata including parameter information
func (r *Registry) Build(client *rest.Client, repo config.Repository) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.built {
		return
	}
	r.built = true

	clientType := reflect.TypeOf(client)
	clientVal := reflect.ValueOf(client)

	// Iterate over all methods on *rest.Client
	for i := 0; i < clientType.NumMethod(); i++ {
		method := clientType.Method(i)

		// Skip non-exported or special methods
		if !isAPIGroupMethod(method) {
			continue
		}

		groupName := toSnakeCase(method.Name)
		groupVal := clientVal.MethodByName(method.Name).Call(nil)[0]

		r.discoverGroup(groupName, groupVal, repo)
	}
}

// discoverGroup discovers endpoints within an API group (e.g., Git, Builds).
// It looks for methods that return scoped clients and methods that are
// direct API calls.
func (r *Registry) discoverGroup(groupName string, groupVal reflect.Value, repo config.Repository) {
	groupType := groupVal.Type()

	for i := 0; i < groupType.NumMethod(); i++ {
		method := groupType.Method(i)

		if !method.IsExported() {
			continue
		}

		// Check if this method returns a scoped client (takes repo/project arg)
		if isScopedClientFactory(method) {
			r.discoverScopedClient(groupName, groupVal, method, repo)
			continue
		}

		// Check if this is a direct API method (takes context as first arg)
		if isAPIMethod(method) {
			path := groupName + "." + toSnakeCase(method.Name)
			endpoint := buildEndpoint(path, groupVal, method)
			r.register(path, endpoint)
		}
	}
}

// discoverScopedClient discovers endpoints from a scoped client factory method.
// For example, Git.PRs(repo) returns GitPRs, which has methods like List, ByID.
func (r *Registry) discoverScopedClient(
	groupName string,
	groupVal reflect.Value,
	factoryMethod reflect.Method,
	repo config.Repository,
) {
	scopeName := toSnakeCase(factoryMethod.Name)

	// Get the scoped client type to discover its methods
	scopeType := factoryMethod.Type.Out(0)

	// Create the scoped client by calling the factory
	repoVal := reflect.ValueOf(repo)
	scopeVal := groupVal.MethodByName(factoryMethod.Name).Call([]reflect.Value{repoVal})[0]

	for i := 0; i < scopeType.NumMethod(); i++ {
		method := scopeType.Method(i)

		if !method.IsExported() || !isAPIMethod(method) {
			continue
		}

		path := groupName + "." + scopeName + "." + toSnakeCase(method.Name)
		endpoint := buildEndpoint(path, scopeVal, method)
		r.register(path, endpoint)
	}
}

// register adds an endpoint to the registry and updates the tree structure.
func (r *Registry) register(path string, endpoint *Endpoint) {
	r.endpoints[path] = endpoint

	// Build tree structure
	parts := strings.Split(path, ".")
	current := r.tree
	for _, part := range parts {
		if current.children == nil {
			current.children = make(map[string]*node)
		}
		if _, ok := current.children[part]; !ok {
			current.children[part] = &node{
				name:     part,
				children: make(map[string]*node),
			}
		}
		current = current.children[part]
	}
	current.endpoint = endpoint
}

// Get returns the endpoint for the given path, or nil if not found.
func (r *Registry) Get(path string) *Endpoint {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.endpoints[path]
}

// Paths returns all registered endpoint paths, sorted alphabetically.
func (r *Registry) Paths() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	paths := make([]string, 0, len(r.endpoints))
	for path := range r.endpoints {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

// Complete returns completion suggestions for the given prefix.
// It returns both partial path completions and full endpoint paths.
func (r *Registry) Complete(prefix string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if prefix == "" {
		// Return top-level groups
		var results []string
		for name := range r.tree.children {
			results = append(results, name)
		}
		sort.Strings(results)
		return results
	}

	parts := strings.Split(prefix, ".")
	current := r.tree

	// Navigate to the current position in the tree
	for i, part := range parts[:len(parts)-1] {
		child, ok := current.children[part]
		if !ok {
			return nil
		}
		current = child
		_ = i // compiler hint
	}

	// Get the last part being typed (may be partial)
	lastPart := parts[len(parts)-1]
	basePath := strings.Join(parts[:len(parts)-1], ".")
	if basePath != "" {
		basePath += "."
	}

	var results []string
	for name, child := range current.children {
		if strings.HasPrefix(name, lastPart) {
			fullPath := basePath + name
			// Add a dot if there are more children
			if len(child.children) > 0 {
				results = append(results, fullPath)
			}
			// Also add if this is a valid endpoint
			if child.endpoint != nil {
				results = append(results, fullPath)
			}
		}
	}
	sort.Strings(results)
	return results
}

// isAPIGroupMethod checks if a method returns an API group client.
// API group methods have no arguments (besides receiver) and return a single value.
func isAPIGroupMethod(m reflect.Method) bool {
	if !m.IsExported() {
		return false
	}

	// Skip Identity - it takes context, not a group factory
	if m.Name == "Identity" {
		return false
	}

	t := m.Type
	// Method on value receiver: (receiver) -> (result)
	// NumIn() == 1 means just the receiver
	return t.NumIn() == 1 && t.NumOut() == 1
}

// isScopedClientFactory checks if a method is a factory that creates a scoped client.
// These methods take a config.Repository and return a scoped client.
func isScopedClientFactory(m reflect.Method) bool {
	t := m.Type
	// (receiver, repo) -> (scopedClient)
	if t.NumIn() != 2 || t.NumOut() != 1 {
		return false
	}

	// Check if input is config.Repository
	repoType := reflect.TypeOf(config.Repository{})
	return t.In(1) == repoType
}

// isAPIMethod checks if a method is an API call method.
// API methods take context.Context as the first argument.
func isAPIMethod(m reflect.Method) bool {
	t := m.Type
	if t.NumIn() < 2 {
		return false
	}

	// First arg (after receiver) should be context.Context
	ctxType := reflect.TypeOf((*context.Context)(nil)).Elem()
	return t.In(1).Implements(ctxType)
}

// toSnakeCase converts a PascalCase or camelCase string to snake_case.
// It uses a simple rule: insert underscore before uppercase letters only when
// the previous character is lowercase. This handles common patterns correctly:
//   - "PRs" -> "prs" (acronyms stay together)
//   - "ByID" -> "by_id"
//   - "ForProject" -> "for_project"
//   - "RepoInfo" -> "repo_info"
func toSnakeCase(s string) string {
	var result strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			// Only add underscore when transitioning from lowercase to uppercase.
			// This keeps acronyms like "PRs", "ID", "URL" together.
			if i > 0 && unicode.IsLower(runes[i-1]) {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

package api

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/rest"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/util"
	"github.com/spf13/cobra"
)

//go:embed api.md
var doc string

const (
	outputJSON = "json"
	outputYAML = "yaml"
)

// registry is the global API endpoint registry, built lazily on first use.
var registry = NewRegistry()

// Cmd returns the api command with dynamic subcommands for each API endpoint.
// The command structure is built by reflecting on the rest.Client to discover
// available APIs automatically.
func Cmd() *cobra.Command {
	output := util.NewEnumFlag(outputJSON, outputYAML).Default(outputJSON)

	cmd := &cobra.Command{
		Use:               "api <endpoint> [flags]",
		Aliases:           []string{"a"},
		Short:             "Call Azure DevOps REST APIs directly",
		Long:              doc,
		ValidArgsFunction: completeEndpoint,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAPI(cmd, args, output)
		},
	}

	flags := cmd.Flags()
	flags.VarP(output, "output", "o", "output format (json, yaml)")
	flags.Bool("list", false, "list all available API endpoints")
	flags.Bool("describe", false, "describe an endpoint's parameters")

	output.RegisterCompletion(cmd, "output")

	return cmd
}

// runAPI executes the api command.
func runAPI(cmd *cobra.Command, args []string, output *util.EnumFlag[string]) error {
	ctx := cmd.Context()
	cfg := config.From(ctx)

	// Initialize registry
	token, err := cfg.Token()
	if err != nil {
		return err
	}
	client := rest.New(token)
	registry.Build(client, cfg.Repository)

	// Handle --list flag
	if list, _ := cmd.Flags().GetBool("list"); list {
		return listEndpoints()
	}

	// Require endpoint argument
	if len(args) == 0 {
		return errors.New("endpoint path required. Use --list to see available endpoints")
	}

	endpointPath := args[0]
	endpoint := registry.Get(endpointPath)
	if endpoint == nil {
		// Suggest similar endpoints
		suggestions := registry.Complete(endpointPath)
		if len(suggestions) > 0 {
			return fmt.Errorf(
				"unknown endpoint %q. Did you mean: %s",
				endpointPath,
				strings.Join(suggestions, ", "),
			)
		}
		return fmt.Errorf(
			"unknown endpoint %q. Use --list to see available endpoints",
			endpointPath,
		)
	}

	// Handle --describe flag
	if describe, _ := cmd.Flags().GetBool("describe"); describe {
		return describeEndpoint(endpoint)
	}

	// Build arguments from positional args (key=value format)
	apiArgs, err := parseArgs(cmd, endpoint, args[1:])
	if err != nil {
		return err
	}

	// Invoke the API
	result, err := endpoint.Invoke(ctx, apiArgs)
	if err != nil {
		return fmt.Errorf("API call failed: %w", err)
	}

	// Output result
	return outputResult(result, output.Value())
}

// listEndpoints prints all available API endpoints.
func listEndpoints() error {
	paths := registry.Paths()

	// Group by first component
	groups := make(map[string][]string)
	for _, path := range paths {
		parts := strings.SplitN(path, ".", 2)
		group := parts[0]
		groups[group] = append(groups[group], path)
	}

	var sb strings.Builder
	sb.WriteString("Available API endpoints:\n\n")

	// Sort groups
	groupNames := make([]string, 0, len(groups))
	for g := range groups {
		groupNames = append(groupNames, g)
	}

	for _, group := range groupNames {
		sb.WriteString(group + ":\n")
		for _, path := range groups[group] {
			sb.WriteString(fmt.Sprintf("  %s\n", path))
		}
		sb.WriteString("\n")
	}

	fmt.Print(sb.String())
	return nil
}

// describeEndpoint prints detailed information about an endpoint.
func describeEndpoint(endpoint *Endpoint) error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Endpoint: %s\n\n", endpoint.Path))

	if len(endpoint.Params) == 0 {
		sb.WriteString("Parameters: none\n")
	} else {
		sb.WriteString("Parameters:\n")
		for _, param := range endpoint.Params {
			required := ""
			if param.Required {
				required = " (required)"
			}
			sb.WriteString(fmt.Sprintf("  --%s  %s%s\n", param.Name, param.Type, required))

			// Show struct fields
			for _, field := range param.Fields {
				fieldRequired := ""
				if field.Required {
					fieldRequired = " (required)"
				}
				sb.WriteString(
					fmt.Sprintf("    --%s.%s  %s%s\n",
						param.Name, field.Name, field.Type, fieldRequired),
				)
				if len(field.EnumValues) > 0 {
					sb.WriteString(
						fmt.Sprintf("      values: %s\n",
							strings.Join(field.EnumValues, ", ")),
					)
				}
			}

			if len(param.EnumValues) > 0 {
				sb.WriteString(
					fmt.Sprintf("    values: %s\n",
						strings.Join(param.EnumValues, ", ")),
				)
			}
		}
	}

	fmt.Print(sb.String())
	return nil
}

// parseArgs extracts API arguments from positional args in key=value format.
// Supports multiple formats:
//   - Simple: id=123
//   - Nested: list_query.top=10
//   - Shorthand: top=10 (automatically prefixed if unambiguous)
func parseArgs(
	_ *cobra.Command,
	endpoint *Endpoint,
	positionalArgs []string,
) (map[string]string, error) {
	args := make(map[string]string)

	// Parse positional arguments as key=value pairs
	for _, arg := range positionalArgs {
		key, value, found := strings.Cut(arg, "=")
		if !found {
			return nil, fmt.Errorf("invalid argument %q: expected key=value format", arg)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		// Try to resolve shorthand keys to full paths
		resolvedKey := resolveParamKey(endpoint, key)
		args[resolvedKey] = value
	}

	return args, nil
}

// resolveParamKey attempts to resolve a shorthand parameter key to its full path.
// For example, "id" might resolve to "pr_id" if that's the only match.
// If the key already contains a dot or is a direct match, it's returned as-is.
func resolveParamKey(endpoint *Endpoint, key string) string {
	// If the key already has a dot, assume it's a full path
	if strings.Contains(key, ".") {
		return key
	}

	// Check for direct parameter match
	for _, param := range endpoint.Params {
		if param.Name == key {
			return key
		}
	}

	// Check for field match and return a full path
	for _, param := range endpoint.Params {
		for _, field := range param.Fields {
			if field.Name == key {
				return param.Name + "." + field.Name
			}
		}
	}

	// No match found, return as-is
	return key
}

// outputResult formats and prints the API result.
func outputResult(result any, format string) error {
	if result == nil {
		return nil
	}

	switch format {
	case outputYAML:
		return styles.DumpYAML(result)
	default:
		return styles.DumpJSON(result)
	}
}

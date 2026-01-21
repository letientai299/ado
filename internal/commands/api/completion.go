package api

import (
	"strings"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest"
	"github.com/spf13/cobra"
)

// completeEndpoint provides tab completion for API endpoint paths and arguments.
// It builds the registry lazily and returns matching completions.
func completeEndpoint(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	// Try to build a registry for completion
	// This may fail if not in a valid repo context, which is fine
	buildRegistryForCompletion(cmd)

	// If we don't have an endpoint yet, complete endpoint paths
	if len(args) == 0 {
		return completeEndpointPath(toComplete)
	}

	// We have an endpoint, complete key=value arguments
	endpointPath := args[0]
	endpoint := registry.Get(endpointPath)
	if endpoint == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return completeArgument(endpoint, toComplete)
}

// completeEndpointPath provides completion for API endpoint paths.
func completeEndpointPath(toComplete string) ([]string, cobra.ShellCompDirective) {
	completions := registry.Complete(toComplete)

	// Add description suffix to help users
	var described []string
	for _, c := range completions {
		endpoint := registry.Get(c)
		if endpoint != nil {
			// This is a valid endpoint
			described = append(described, c+"\tAPI endpoint")
		} else {
			// This is a partial path (more children available)
			described = append(described, c+"\texpand with .")
		}
	}

	return described, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
}

// completeArgument provides completion for key=value arguments.
// It completes parameter keys and, after the = sign, parameter values.
func completeArgument(endpoint *Endpoint, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Check if we're completing a value (after =)
	if key, valuePrefix, found := strings.Cut(toComplete, "="); found {
		// Complete the value part with the partial value as a filter
		values := CompleteParamValue(endpoint, key, valuePrefix)
		if len(values) == 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		// Prepend the key= to each value
		var completions []string
		for _, v := range values {
			completions = append(completions, key+"="+v)
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	// Complete the key part
	return completeParamKey(endpoint, toComplete)
}

// completeParamKey provides completion for parameter key names.
func completeParamKey(endpoint *Endpoint, prefix string) ([]string, cobra.ShellCompDirective) {
	var completions []string
	prefix = strings.ToLower(prefix)

	for _, param := range endpoint.Params {
		// Add the main parameter name
		if strings.HasPrefix(strings.ToLower(param.Name), prefix) {
			desc := param.Type
			if param.Required {
				desc += " (required)"
			}
			completions = append(completions, param.Name+"=\t"+desc)
		}

		// Add field names (both short and full form)
		for _, field := range param.Fields {
			fullName := param.Name + "." + field.Name

			// Full form: param.field
			if strings.HasPrefix(strings.ToLower(fullName), prefix) {
				desc := field.Type
				if field.Required {
					desc += " (required)"
				}
				completions = append(completions, fullName+"=\t"+desc)
			}

			// Short form: just field name
			if strings.HasPrefix(strings.ToLower(field.Name), prefix) {
				desc := field.Type + " (shorthand for " + fullName + ")"
				completions = append(completions, field.Name+"=\t"+desc)
			}
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
}

// buildRegistryForCompletion attempts to build the registry for completion.
// It handles errors gracefully since completion shouldn't fail.
func buildRegistryForCompletion(cmd *cobra.Command) {
	ctx := cmd.Context()
	if ctx == nil {
		return
	}

	cfg := config.From(ctx)
	if cfg == nil {
		return
	}

	token, err := cfg.Token()
	if err != nil {
		return
	}

	client := rest.New(token)
	registry.Build(client, cfg.Repository)
}

// knownEnums maps type names to their known enum values.
// This provides completion for Azure DevOps enum types that can't be
// discovered via reflection alone.
//
// To add new enums:
//  1. Add the type name (as it appears in reflect.Type.Name()) as the key
//  2. Add the valid string values as the slice
var knownEnums = map[string][]string{
	// Pull request statuses
	"PullRequestStatus": {
		string(models.PullRequestStatusAbandoned),
		string(models.PullRequestStatusActive),
		string(models.PullRequestStatusAll),
		string(models.PullRequestStatusCompleted),
		string(models.PullRequestStatusNotSet),
	},

	// Build statuses
	"BuildStatus": {
		"all",
		"cancelling",
		"completed",
		"inProgress",
		"none",
		"notStarted",
		"postponed",
	},

	// Build results
	"BuildResult": {
		"canceled",
		"failed",
		"none",
		"partiallySucceeded",
		"succeeded",
	},

	// PR time range types
	"PullRequestTimeRangeType": {
		string(models.PullRequestTimeRangeTypeCreated),
		string(models.PullRequestTimeRangeTypeClosed),
	},
}

// CompleteParamValue provides completion for parameter values.
// It checks for known enums and returns valid values.
// The paramPath can be a full path (param.field) or shorthand (field).
func CompleteParamValue(
	endpoint *Endpoint,
	paramPath string,
	toComplete string,
) []string {
	// First, resolve the parameter path to find the field type
	fieldType := resolveFieldType(endpoint, paramPath)
	if fieldType == "" {
		return nil
	}

	// Check known enums by type name
	if enumVals := lookupKnownEnum(fieldType); len(enumVals) > 0 {
		return filterByPrefix(enumVals, toComplete)
	}

	// Check endpoint's own enum values
	values := endpoint.FieldEnumValues(paramPath)
	if len(values) > 0 {
		return filterByPrefix(values, toComplete)
	}

	return nil
}

// resolveFieldType finds the type string for a parameter path.
// Handles both full paths (param.field.subfield) and shorthand names (field).
func resolveFieldType(endpoint *Endpoint, path string) string {
	// First, try direct match on a parameter name
	for _, param := range endpoint.Params {
		if param.Name == path {
			return param.Type
		}
	}

	// Then try matching fields (full path or shorthand)
	for _, param := range endpoint.Params {
		for _, field := range param.Fields {
			// Full path match
			if field.Name == path {
				return field.Type
			}
			// Shorthand match (last component of the field name)
			lastDot := strings.LastIndex(field.Name, ".")
			if lastDot >= 0 {
				shortName := field.Name[lastDot+1:]
				if shortName == path {
					return field.Type
				}
			}
		}
	}

	return ""
}

// lookupKnownEnum returns enum values for a type name.
func lookupKnownEnum(typeName string) []string {
	// Extract the base type name (remove package prefix and pointer)
	typeName = strings.TrimPrefix(typeName, "*")
	if idx := strings.LastIndex(typeName, "."); idx >= 0 {
		typeName = typeName[idx+1:]
	}

	return knownEnums[typeName]
}

// filterByPrefix filters values that start with the given prefix.
func filterByPrefix(values []string, prefix string) []string {
	if prefix == "" {
		return values
	}

	var filtered []string
	prefix = strings.ToLower(prefix)
	for _, v := range values {
		if strings.HasPrefix(strings.ToLower(v), prefix) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

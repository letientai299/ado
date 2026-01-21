package api

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Endpoint represents a single API endpoint that can be invoked.
// It contains metadata about the endpoint's parameters and the
// reflection information needed to invoke it.
type Endpoint struct {
	// Path is the CLI path to this endpoint (e.g., "git.prs.list")
	Path string

	// Params describes the parameters this endpoint accepts.
	// Order matches the method signature (excluding context).
	Params []Param

	// receiver is the reflected value to call the method on
	receiver reflect.Value

	// method is the reflected method to call
	method reflect.Method
}

// Param describes a single parameter for an API endpoint.
type Param struct {
	// Name is the parameter name (derived from type or struct field)
	Name string

	// Type is the Go type name (e.g., "int32", "string", "ListOptions")
	Type string

	// Kind indicates the parameter category
	Kind ParamKind

	// Required indicates if this parameter must be provided
	Required bool

	// Fields contains struct field information if Kind is ParamKindStruct
	Fields []ParamField

	// EnumValues contains valid values if this is an enum type
	EnumValues []string

	// goType is the reflect.Type for this parameter
	goType reflect.Type
}

// ParamField describes a field within a struct parameter.
type ParamField struct {
	// Name is the field name as it appears in JSON/YAML
	Name string

	// GoName is the original Go field name
	GoName string

	// Type is the field's Go type
	Type string

	// Required indicates if this field must be provided
	Required bool

	// EnumValues contains valid values if this is an enum type
	EnumValues []string

	// goType is the reflect.Type for this field
	goType reflect.Type
}

// ParamKind categorizes parameters for CLI handling.
type ParamKind int

const (
	ParamKindPrimitive ParamKind = iota // int, string, bool, etc.
	ParamKindStruct                     // struct with fields
	ParamKindSlice                      // slice of values
	ParamKindPointer                    // pointer to primitive or struct
)

// buildEndpoint creates an Endpoint from a reflected method.
// It extracts parameter information for CLI flags and completion.
func buildEndpoint(path string, receiver reflect.Value, method reflect.Method) *Endpoint {
	endpoint := &Endpoint{
		Path:     path,
		receiver: receiver,
		method:   method,
	}

	// Skip receiver and context.Context, process remaining params
	methodType := method.Type
	for i := 2; i < methodType.NumIn(); i++ {
		paramType := methodType.In(i)
		param := buildParam(paramType, i-2)
		endpoint.Params = append(endpoint.Params, param)
	}

	return endpoint
}

// buildParam creates a Param from a reflected type.
func buildParam(t reflect.Type, index int) Param {
	param := Param{
		Name:   inferParamName(t, index),
		Type:   t.String(),
		goType: t,
	}

	switch t.Kind() {
	case reflect.Ptr:
		param.Kind = ParamKindPointer
		// Check the element type
		elem := t.Elem()
		if elem.Kind() == reflect.Struct {
			param.Fields = extractStructFields(elem)
		}
		param.EnumValues = inferEnumValues(elem)

	case reflect.Struct:
		param.Kind = ParamKindStruct
		param.Required = true
		param.Fields = extractStructFields(t)

	case reflect.Slice:
		param.Kind = ParamKindSlice

	default:
		param.Kind = ParamKindPrimitive
		param.Required = true
		param.EnumValues = inferEnumValues(t)
	}

	return param
}

// extractStructFields returns field metadata for a struct type.
// It recursively extracts fields from nested structs (via pointers),
// but limits depth to prevent infinite recursion with circular types.
func extractStructFields(t reflect.Type) []ParamField {
	visited := make(map[reflect.Type]bool)
	return extractStructFieldsWithPrefix(t, "", 0, visited)
}

// maxFieldDepth limits recursive struct field extraction to prevent
// infinite loops with circular type references.
const maxFieldDepth = 3

// extractStructFieldsWithPrefix recursively extracts fields with a name prefix.
// The depth parameter prevents infinite recursion.
func extractStructFieldsWithPrefix(
	t reflect.Type,
	prefix string,
	depth int,
	visited map[reflect.Type]bool,
) []ParamField {
	if depth > maxFieldDepth {
		return nil
	}

	// Prevent cycles
	if visited[t] {
		return nil
	}
	visited[t] = true
	defer delete(visited, t)

	var fields []ParamField

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldName := fieldJSONName(field)
		if prefix != "" {
			fieldName = prefix + "." + fieldName
		}

		pf := ParamField{
			GoName: field.Name,
			Name:   fieldName,
			Type:   field.Type.String(),
			goType: field.Type,
		}

		// Pointers are optional, non-pointers are required
		if field.Type.Kind() != reflect.Ptr {
			pf.Required = true
		}

		pf.EnumValues = inferEnumValues(field.Type)
		fields = append(fields, pf)

		// Recursively extract nested struct fields
		elemType := field.Type
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		if elemType.Kind() == reflect.Struct && !isStdlibType(elemType) {
			nestedFields := extractStructFieldsWithPrefix(elemType, fieldName, depth+1, visited)
			fields = append(fields, nestedFields...)
		}
	}

	return fields
}

// isStdlibType checks if a type is from the standard library that shouldn't
// be recursively expanded (like time.Time).
func isStdlibType(t reflect.Type) bool {
	pkg := t.PkgPath()
	return pkg == "time" || pkg == "encoding/json" || pkg == "net/url"
}

// fieldJSONName returns the JSON field name from the tag or falls back to the Go name.
func fieldJSONName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		tag = field.Tag.Get("yaml")
	}
	if tag == "" {
		return toSnakeCase(field.Name)
	}

	// Handle "name,omitempty" format
	name, _, _ := strings.Cut(tag, ",")
	if name == "" || name == "-" {
		return toSnakeCase(field.Name)
	}
	return name
}

// inferParamName derives a parameter name from its type.
// For primitive types, it uses meaningful names based on common API patterns.
// For struct/complex types, it uses the type name converted to snake_case.
func inferParamName(t reflect.Type, index int) string {
	// For primitive types, use meaningful generic names
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "id"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "id"
	case reflect.String:
		return "value"
	case reflect.Bool:
		return "enabled"
	}

	// Use the type name if available
	name := t.Name()
	if name == "" {
		// For pointers, slices, etc., use the element type
		if t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
			elem := t.Elem()
			name = elem.Name()
			// If element is also primitive, use meaningful name
			switch elem.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return "id"
			case reflect.String:
				return "value"
			}
		}
	}

	if name != "" {
		return toSnakeCase(name)
	}

	return fmt.Sprintf("arg%d", index)
}

// inferEnumValues attempts to detect enum values for known types.
// This is a best-effort approach for common Azure DevOps enums.
func inferEnumValues(t reflect.Type) []string {
	// Unwrap pointer
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Check for known enum types by name
	typeName := t.Name()
	if strings.HasSuffix(typeName, "Status") ||
		strings.HasSuffix(typeName, "Type") ||
		strings.HasSuffix(typeName, "State") {
		// These are likely enums, but we can't get values via reflection
		// The completion.go file will handle known enums explicitly
		return nil
	}

	return nil
}

// Invoke calls the endpoint with the provided arguments.
// Arguments are provided as a map of parameter name to value (as strings).
// For struct parameters, the map should contain field names with dot notation
// (e.g., "search_criteria.status" for nested fields).
//
// Returns the result as an interface{} that can be JSON/YAML encoded,
// and any error from the API call.
func (e *Endpoint) Invoke(ctx context.Context, args map[string]string) (any, error) {
	callArgs := []reflect.Value{reflect.ValueOf(ctx)}

	for _, param := range e.Params {
		argVal, err := e.buildArgument(param, args)
		if err != nil {
			return nil, fmt.Errorf("parameter %q: %w", param.Name, err)
		}
		callArgs = append(callArgs, argVal)
	}

	// Call the method
	results := e.receiver.MethodByName(e.method.Name).Call(callArgs)

	// Handle return values (typically (result, error) or just error)
	if len(results) == 0 {
		return nil, nil
	}

	// Check for error in last return value
	errIdx := len(results) - 1
	if !results[errIdx].IsNil() {
		return nil, results[errIdx].Interface().(error)
	}

	// Return the result if there is one
	if len(results) > 1 {
		return results[0].Interface(), nil
	}

	return nil, nil
}

// buildArgument constructs a reflect.Value for a parameter from string arguments.
func (e *Endpoint) buildArgument(param Param, args map[string]string) (reflect.Value, error) {
	switch param.Kind {
	case ParamKindPrimitive:
		return e.buildPrimitive(param.goType, args[param.Name])

	case ParamKindPointer:
		// Check if any value is provided for this param
		val, ok := args[param.Name]
		if !ok || val == "" {
			// Check for nested field values
			hasNested := false
			prefix := param.Name + "."
			for k := range args {
				if strings.HasPrefix(k, prefix) {
					hasNested = true
					break
				}
			}
			if !hasNested {
				// Return nil pointer
				return reflect.Zero(param.goType), nil
			}
		}

		// Build the element value and return pointer to it
		elem := param.goType.Elem()
		if elem.Kind() == reflect.Struct {
			structVal, err := e.buildStruct(elem, param.Name, param.Fields, args)
			if err != nil {
				return reflect.Value{}, err
			}
			ptr := reflect.New(elem)
			ptr.Elem().Set(structVal)
			return ptr, nil
		}

		primVal, err := e.buildPrimitive(elem, val)
		if err != nil {
			return reflect.Value{}, err
		}
		ptr := reflect.New(elem)
		ptr.Elem().Set(primVal)
		return ptr, nil

	case ParamKindStruct:
		return e.buildStruct(param.goType, param.Name, param.Fields, args)

	case ParamKindSlice:
		// For slices, expect comma-separated values or JSON array
		val := args[param.Name]
		if val == "" {
			return reflect.MakeSlice(param.goType, 0, 0), nil
		}
		return e.buildSlice(param.goType, val)

	default:
		return reflect.Zero(param.goType), nil
	}
}

// buildPrimitive converts a string to a primitive reflect.Value.
func (e *Endpoint) buildPrimitive(t reflect.Type, val string) (reflect.Value, error) {
	switch t.Kind() {
	case reflect.String:
		return reflect.ValueOf(val).Convert(t), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val == "" {
			return reflect.Zero(t), nil
		}
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("invalid integer: %w", err)
		}
		return reflect.ValueOf(n).Convert(t), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val == "" {
			return reflect.Zero(t), nil
		}
		n, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("invalid unsigned integer: %w", err)
		}
		return reflect.ValueOf(n).Convert(t), nil

	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("invalid boolean: %w", err)
		}
		return reflect.ValueOf(b), nil

	case reflect.Float32, reflect.Float64:
		if val == "" {
			return reflect.Zero(t), nil
		}
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("invalid float: %w", err)
		}
		return reflect.ValueOf(f).Convert(t), nil

	default:
		return reflect.Zero(t), fmt.Errorf("unsupported primitive type: %s", t.Kind())
	}
}

// buildStruct constructs a struct value from arguments.
// Arguments should use dot notation for nested fields (e.g., "param.field").
func (e *Endpoint) buildStruct(
	t reflect.Type,
	prefix string,
	fields []ParamField,
	args map[string]string,
) (reflect.Value, error) {
	structVal := reflect.New(t).Elem()

	for _, field := range fields {
		argKey := prefix + "." + field.Name
		val, ok := args[argKey]
		if !ok {
			// Also try without prefix for top-level struct
			val, ok = args[field.Name]
		}

		if !ok || val == "" {
			if field.Required {
				// Required field missing - that's okay for API calls,
				// the server will validate
			}
			continue
		}

		fieldVal := structVal.FieldByName(field.GoName)
		if !fieldVal.IsValid() || !fieldVal.CanSet() {
			continue
		}

		var newVal reflect.Value
		var err error

		if field.goType.Kind() == reflect.Ptr {
			// Handle pointer fields
			elem := field.goType.Elem()
			primVal, e := e.buildPrimitive(elem, val)
			if e != nil {
				err = e
			} else {
				ptr := reflect.New(elem)
				ptr.Elem().Set(primVal)
				newVal = ptr
			}
		} else {
			newVal, err = e.buildPrimitive(field.goType, val)
		}

		if err != nil {
			return reflect.Value{}, fmt.Errorf("field %q: %w", field.Name, err)
		}

		fieldVal.Set(newVal)
	}

	return structVal, nil
}

// buildSlice constructs a slice from a comma-separated string or JSON array.
func (e *Endpoint) buildSlice(t reflect.Type, val string) (reflect.Value, error) {
	elemType := t.Elem()

	// Try JSON array first
	if strings.HasPrefix(val, "[") {
		slice := reflect.MakeSlice(t, 0, 0)
		var items []json.RawMessage
		if err := json.Unmarshal([]byte(val), &items); err != nil {
			return reflect.Value{}, fmt.Errorf("invalid JSON array: %w", err)
		}

		for _, item := range items {
			elem := reflect.New(elemType).Elem()
			if err := json.Unmarshal(item, elem.Addr().Interface()); err != nil {
				return reflect.Value{}, fmt.Errorf("invalid array element: %w", err)
			}
			slice = reflect.Append(slice, elem)
		}
		return slice, nil
	}

	// Comma-separated values for primitives
	parts := strings.Split(val, ",")
	slice := reflect.MakeSlice(t, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		elem, err := e.buildPrimitive(elemType, part)
		if err != nil {
			return reflect.Value{}, err
		}
		slice = reflect.Append(slice, elem)
	}
	return slice, nil
}

// ParamNames returns the names of all parameters and their fields.
// This is used for flag completion.
func (e *Endpoint) ParamNames() []string {
	var names []string

	for _, param := range e.Params {
		if param.Kind == ParamKindPrimitive {
			names = append(names, param.Name)
		} else if param.Kind == ParamKindStruct || param.Kind == ParamKindPointer {
			for _, field := range param.Fields {
				names = append(names, param.Name+"."+field.Name)
			}
			// Also add the param name itself for JSON input
			names = append(names, param.Name)
		}
	}

	sort.Strings(names)
	return names
}

// FieldEnumValues returns enum values for a field, if known.
// The fieldPath uses dot notation (e.g., "search_criteria.status").
func (e *Endpoint) FieldEnumValues(fieldPath string) []string {
	for _, param := range e.Params {
		if param.Name == fieldPath {
			return param.EnumValues
		}
		for _, field := range param.Fields {
			if param.Name+"."+field.Name == fieldPath {
				return field.EnumValues
			}
		}
	}
	return nil
}

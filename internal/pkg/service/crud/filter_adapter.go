package crud

import (
	"reflect"
	"regexp"
	"strings"
)

// BuildFilter converts a GraphQL filter input struct into a FilterGroup.
// It uses reflection and struct tags to map fields to filter conditions.
//
// Struct tags format:
// `filter:"field:db_column_name;operator:eq"`
//   - `field` (optional): The database column name. Defaults to the struct field name converted to snake_case.
//   - `operator` (optional): The filter operator (e.g., "eq", "like", "gt"). Defaults to "eq".
func BuildFilter(input interface{}) *FilterGroup {
	if input == nil {
		return nil
	}

	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	typ := val.Type()
	filters := []interface{}{}

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldTyp := typ.Field(i)

		// Skip nil pointer fields, as they indicate the filter is not being applied.
		if fieldVal.Kind() == reflect.Ptr && fieldVal.IsNil() {
			continue
		}

		// Dereference pointer to get the actual value for the filter.
		if fieldVal.Kind() == reflect.Ptr {
			fieldVal = fieldVal.Elem()
		}

		// Get tag information
		tag := fieldTyp.Tag.Get("filter")
		if tag == "" {
			continue
		}

		tagParts := strings.Split(tag, ";")
		filterConfig := make(map[string]string)
		for _, part := range tagParts {
			kv := strings.SplitN(part, ":", 2)
			if len(kv) == 2 {
				filterConfig[kv[0]] = kv[1]
			}
		}

		dbField, ok := filterConfig["field"]
		if !ok {
			dbField = toSnakeCase(fieldTyp.Name)
		}

		operator, ok := filterConfig["operator"]
		if !ok {
			operator = "eq" // Default operator
		}

		filters = append(filters, Filter{
			Field:    dbField,
			Operator: FilterOperator(operator),
			Value:    fieldVal.Interface(),
		})
	}

	if len(filters) == 0 {
		return nil
	}

	// Default to AND for all filters at the top level.
	// You could make this configurable via a struct tag on the struct itself if needed.
	return &FilterGroup{
		Operator: LogicalAnd,
		Filters:  filters,
	}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// toSnakeCase converts a CamelCase string to snake_case.
func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

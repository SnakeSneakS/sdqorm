package parser

import (
	"fmt"
	"reflect"
)

// handleFuncOnSpecificTag runs function ${f} about specific field ${tag} on ${target}
func handleFuncOnSpecificTag(
	specificTag string, target interface{},
	f func(field reflect.Value, typeField reflect.StructField, tag string) error,
) error {
	// Get the reflection value of the result
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Pointer || val.IsNil() {
		return fmt.Errorf("target must be a pointer to a struct")
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	// Iterate over the fields of the struct
	found := false
	if val.NumField() == 0 {
		found = true
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag.Get(specificTag)

		if tag == "" {
			continue
		}

		found = true
		if err := f(field, typeField, tag); err != nil {
			return err
		}
	}
	if !found {
		return fmt.Errorf("tag %s is required, but not found for target %v", specificTag, target)
	}

	return nil
}

// Check for required keys
//func checkRequiredKeys(tagParts []string, required []TagPartKey) error {
//	for _, partKey := range required {
//		found := false
//		for _, part := range tagParts {
//			if strings.HasPrefix(part, fmt.Sprintf("%s:", partKey)) {
//				found = true
//				break
//			}
//		}
//		if !found {
//			return fmt.Errorf("key %s is required for usage of tag %s", partKey, SdqormTagName)
//		}
//	}
//	return nil
//}

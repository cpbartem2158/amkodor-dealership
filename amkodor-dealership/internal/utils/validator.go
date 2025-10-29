package utils

import (
	"reflect"
	"strings"
)

// ValidateStruct проверяет структуру на валидность
func ValidateStruct(s interface{}) []string {
	var errors []string
	
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		
		// Проверяем тег "required"
		if required := fieldType.Tag.Get("required"); required == "true" {
			if field.Kind() == reflect.String && field.String() == "" {
				errors = append(errors, fieldType.Name+" is required")
			}
			if field.Kind() == reflect.Ptr && field.IsNil() {
				errors = append(errors, fieldType.Name+" is required")
			}
		}
		
		// Проверяем тег "email"
		if email := fieldType.Tag.Get("email"); email == "true" {
			if field.Kind() == reflect.String && field.String() != "" {
				if !strings.Contains(field.String(), "@") {
					errors = append(errors, fieldType.Name+" must be a valid email")
				}
			}
		}
	}
	
	return errors
}
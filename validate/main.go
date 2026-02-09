package validate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var v *validator.Validate

func init() {
	v = validator.New(validator.WithRequiredStructEnabled())
}

// Validate validates the struct or map based on the provided fields and rules.
func Validate(fields map[string]string, in any) error {
	for field, rule := range fields {
		var value any
		r := reflect.ValueOf(in)

		// Dereference pointer if necessary
		isPointer := false
		for r.Kind() == reflect.Ptr {
			if r.IsNil() {
				break
			}
			isPointer = true
			r = r.Elem()
		}

		if r.IsValid() {
			switch r.Kind() {
			case reflect.Struct:
				fieldValue := r.FieldByName(field)
				if !fieldValue.IsValid() {
					return fmt.Errorf("field %s does not exist in struct", field)
				}
				value = fieldValue.Interface()
				isPointer = fieldValue.Kind() == reflect.Ptr
			case reflect.Map:
				mapValue := r.MapIndex(reflect.ValueOf(field))
				if !mapValue.IsValid() {
					return fmt.Errorf("field %s does not exist in map", field)
				}
				value = mapValue.Interface()
				isPointer = mapValue.Kind() == reflect.Pointer
			default:
				value = in
			}
		} else {
			value = in
		}

		// Workaround for required + numeric 0 on non-pointers
		if !isPointer && value != nil {
			rv := reflect.ValueOf(value)
			isZero := false
			switch rv.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				isZero = rv.Int() == 0
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				isZero = rv.Uint() == 0
			case reflect.Float32, reflect.Float64:
				isZero = rv.Float() == 0
			}

			if isZero && strings.Contains(rule, "required") && (strings.Contains(rule, "gte=0") || strings.Contains(rule, "min=0")) {
				rule = strings.Replace(rule, "required,", "", 1)
				rule = strings.Replace(rule, ",required", "", 1)
				rule = strings.Replace(rule, "required", "", 1)
			}
		}

		// Validate using the provided rule
		if e := v.Var(value, rule); e != nil {
			errMsg := ""
			if validationErrors, ok := e.(validator.ValidationErrors); ok {
				for _, fieldError := range validationErrors {
					if errMsg != "" {
						errMsg += ", "
					}
					errMsg += fieldError.Tag()
				}
			} else {
				errMsg = e.Error()
			}
			return fmt.Errorf("validation failed for field '%s' to rules '%s'", field, errMsg)
		}
	}

	return nil
}

package validate

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var v *validator.Validate

func init() {
	v = validator.New(validator.WithRequiredStructEnabled())
}

// fields should be a map of field name and rule,
// in should be a struct to be validated
// example: 
// fields := map[string]string{
//     "Name": "required,min=3,max=100",
//     "Age":  "required,numeric,min=18,max=100",
// }
// Validate validates the struct based on the provided fields and rules.
func Validate(fields map[string]string, in any) error {
	for field, rule := range fields {
		var value any
		r := reflect.ValueOf(in)

		// Dereference pointer if necessary
		if r.Kind() == reflect.Ptr {
			r = r.Elem()
		}

		if r.IsValid() && r.Kind() == reflect.Struct {
			// Get the field by name and check if it's valid
			fieldValue := r.FieldByName(field)
			if !fieldValue.IsValid() {
				return fmt.Errorf("field %s does not exist", field)
			}
			value = fieldValue.Interface()
		} else {
			value = in
		}

		// Validate using the provided rule
		if e := v.Var(value, rule); e != nil {
			errMsg := ""
			for _, fieldError := range e.(validator.ValidationErrors) {
				if errMsg != "" {
					errMsg += ", "
				}
				errMsg += fieldError.Tag()
			}
			return fmt.Errorf("validation failed for field '%s' to rules '%s'", field, errMsg)
		}
	}

	return nil
}

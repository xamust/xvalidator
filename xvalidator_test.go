package xvalidator

import (
	"gopkg.in/go-playground/validator.v9"
	"strconv"
	"testing"
)

func Test_xValidator_ValidateStructWithOnlyCustomTag(t *testing.T) {
	val := validator.New()
	in := []InputTagsData{
		{"inn",
			"INN must be numeric and contains only 12 digits",
			func(fl validator.FieldLevel) bool {
				inn := fl.Field().String()
				if _, err := strconv.Atoi(inn); err != nil || len(inn) != 12 {
					return false
				}
				return true
			},
		},
	}
	v := &xValidator{translatorInit(val), in, val}
	var testStruct struct {
		INN string `validate:"required,inn"`
	}
	tests := []struct {
		name    string
		inn     string
		wantErr bool
	}{
		{
			name:    "correct INN (custom tag)",
			inn:     "111111111111",
			wantErr: false,
		},
		{
			name:    "incorrect INN (custom tag)",
			inn:     "11111111111a",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStruct.INN = tt.inn
			if err := v.ValidateStruct(testStruct); (err != nil) != tt.wantErr {
				t.Errorf("ValidateStruct() error:\n%s,\nwantErr: %v", err.Error(), tt.wantErr)
			}
		})
	}
}

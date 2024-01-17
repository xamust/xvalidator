package xvalidator

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func Test_xValidator_ValidateStructWithOnlyCustomTag(t *testing.T) {
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
	v := NewXValidator(in...)
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

func Test_xValidator_ValidateStructWithoutCustomTag(t *testing.T) {
	v := NewXValidator()
	var testStruct struct {
		Email string `validate:"required,email"`
	}
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "correct Email",
			email:   "1@1.ru",
			wantErr: false,
		},
		{
			name:    "incorrect Email",
			email:   "1111@111111a",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStruct.Email = tt.email
			if err := v.ValidateStruct(testStruct); (err != nil) != tt.wantErr {
				t.Errorf("ValidateStruct() error:\n%s,\nwantErr: %v", err.Error(), tt.wantErr)
			}
		})
	}
}

func Test_xValidator_ValidateVar(t *testing.T) {
	v := NewXValidator()
	tests := []struct {
		name    string
		valData []InputValData
		wantErr bool
	}{
		{
			name: "correct email",
			valData: []InputValData{
				{
					Key:     "email",
					ValData: "1@1.ru",
				},
			},
			wantErr: false,
		},
		{
			name: "correct email, incorrect passwd len",
			valData: []InputValData{
				{
					Key:     "required,len=9",
					ValData: "1111",
				},
				{
					Key:     "email",
					ValData: "1@1.ru",
				},
			},
			wantErr: true,
		},
		{
			name: "incorrect email",
			valData: []InputValData{
				{
					Key:     "email",
					ValData: "1111@111111a",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.ValidateVar(tt.valData...); (err != nil) != tt.wantErr {
				t.Errorf("ValidateStruct() error:\n%s,\nwantErr: %v", err.Error(), tt.wantErr)
			}
		})
	}
}

func Test_xValidator_Echo(t *testing.T) {
	type testStruct struct {
		Name string `json:"name" validate:"required"`
		INN  string `json:"inn" validate:"required,inn"`
	}

	tests := []struct {
		name     string
		userJSON string
		wantErr  bool
	}{
		{
			name:     "correct custom tag data",
			userJSON: `{"name":"Alice","inn":"111111111111"}`,
			wantErr:  false,
		},
		{
			name:     "incorrect custom tag data",
			userJSON: `{"name":"Alice","inn":"jon@labstack.com"}`,
			wantErr:  true,
		},
	}

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

	innTestStruct := new(testStruct)
	e := echo.New()
	e.Validator = NewXValidator(in...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.userJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if err := c.Bind(innTestStruct); err != nil {
				t.Error(err)
			}
			if err := c.Validate(innTestStruct); (err != nil) != tt.wantErr {
				t.Error(err)
			}
		})
	}
}

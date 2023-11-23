package xvalidator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	"log"
)

type XValidator interface {
	// ValidateStruct validation structs by tags
	ValidateStruct(in interface{}) error

	// ValidateVar validation var by tag
	ValidateVar(field interface{}, tag string) error
}

type xValidator struct {
	trans ut.Translator
	in    []InputTagsData
	*validator.Validate
}

type InputTagsData struct {
	Key            string
	ErrDescription string
	CustomValidation
}

type CustomValidation func(fl validator.FieldLevel) bool

// NewXValidator init new validator with custom tags
//
// Example:
//
//	 var testStruct struct {
//			INN string `validate:"required,inn"`
//		}
//		in := InputTagsData{
//			{"inn",
//				"INN must be numeric and contains only 12 digits",
//				func(fl validator.FieldLevel) bool {
//					inn := fl.Field().String()
//					if  _,err := strconv.Atoi(inn); err != nil || len(inn) != 12 {
//						return false
//					}
//					return true
//				},
//			},
//		}
//		v := NewXValidator(in)
//	 ...
//	 v.ValidateStruct(testStruct)
func NewXValidator(tags ...InputTagsData) XValidator {
	val := validator.New(validator.WithRequiredStructEnabled())
	return &xValidator{
		translatorInit(val),
		tags,
		val,
	}
}

func (v *xValidator) ValidateStruct(tags interface{}) error {
	if len(v.in) > 0 || v.in != nil {
		if err := v.registryCustomTags(); err != nil {
			return err
		}
	}
	return v.translateError(v.Struct(tags))
}

func (v *xValidator) ValidateVar(field interface{}, tag string) error {
	return v.translateError(v.Validate.Var(field, tag))
}

func translatorInit(val *validator.Validate) ut.Translator {
	uni := ut.New(en.New())
	trans, _ := uni.GetTranslator("en")
	if err := enTrans.RegisterDefaultTranslations(val, trans); err != nil {
		log.Fatalf("failed to initialize validator: %s", err)
	}
	return trans
}

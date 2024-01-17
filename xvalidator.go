package xvalidator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	"log"
)

type XValidator interface {
	// ValidateVar validation var by tag
	ValidateVar(in ...InputValData) error

	// Validate validation for echo
	Validate(in interface{}) error
}

type xValidator struct {
	trans    ut.Translator
	in       []InputTagsData
	validate *validator.Validate
}

type InputTagsData struct {
	Key            string
	ErrDescription string
	CustomValidation
}

type InputValData struct {
	Key     string
	ValData interface{}
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

// Validate все кастомные теги грузятся при инициализации валидатора,
// далее в метод передаем структуру
func (v *xValidator) Validate(in interface{}) error {
	if len(v.in) > 0 || v.in != nil {
		if err := v.registryCustomTags(); err != nil {
			return err
		}
	}
	return v.translateError(v.validate.Struct(in))
}

// ValidateVar обертка для дефолтного валидатора,
// в метод передаем структуру(ы) с тегом и данными
func (v *xValidator) ValidateVar(in ...InputValData) error {
	if len(v.in) > 0 || v.in != nil {
		if err := v.registryCustomTags(); err != nil {
			return err
		}
	}
	for _, data := range in {
		if err := v.translateError(v.validate.Var(data.ValData, data.Key)); err != nil {
			return err
		}
	}
	return nil
}

func translatorInit(val *validator.Validate) ut.Translator {
	uni := ut.New(en.New())
	trans, _ := uni.GetTranslator("en")
	if err := enTrans.RegisterDefaultTranslations(val, trans); err != nil {
		log.Fatalf("failed to initialize validator: %s", err)
	}
	return trans
}

package xvalidator

import (
	"errors"
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func (v *xValidator) translateError(err error) error {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	errs := make([]error, 0, len(validatorErrs))
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(v.trans))
		errs = append(errs, translatedErr)
	}
	return errors.Join(errs...)
}

func (v *xValidator) addCustomTranslation(tag string, errMessage string) {
	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, errMessage, false)
	}
	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		tag := fe.Tag()

		t, err := ut.T(tag, fe.Field(), param)
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}
	_ = v.validate.RegisterTranslation(tag, v.trans, registerFn, transFn)
}

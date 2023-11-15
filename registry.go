package xValidator

import (
	"gopkg.in/go-playground/validator.v9"
)

// registryCustomTags custom tags section
func (v *xValidator) registryCustomTags() error {
	for _, s := range v.in {
		if err := v.registryTag(s.Key, s.ErrDescription, s.CustomValidation); err != nil {
			return err
		}
	}
	return nil
}

// registryTag tag registry
func (v *xValidator) registryTag(tagName, errDescr string, fl func(fl validator.FieldLevel) bool) error {
	// register custom tag
	if err := v.RegisterValidation(tagName, fl); err != nil {
		return err
	}
	// register custom err translation
	v.addCustomTranslation(tagName, errDescr)
	return nil
}

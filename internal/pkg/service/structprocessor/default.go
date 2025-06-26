package structprocessor

import (
	"context"
	"errors"
	"fmt"

	validator "github.com/go-playground/validator/v10"
)

func (m *StructProcessorServiceModule) DefaultTagValidationErrorHandler(ctx context.Context, languageId string, validationErrors validator.ValidationErrors) error {
	if len(validationErrors) > 0 {
		for _, r := range validationErrors {
			fieldName := m.localizer.Localize(languageId, r.Field(), nil)
			if r.Tag() == "required" {
				return errors.New(m.localizer.Localize(
					languageId,
					"ErrorFieldRequired",
					map[string]interface{}{
						"FieldName": fieldName,
					},
				))
			}
			if r.Tag() == "number" {
				return errors.New(m.localizer.Localize(
					languageId,
					"ErrorStringNotValidNumber",
					map[string]interface{}{
						"StringName": fieldName,
					},
				))
			}
			if r.Tag() == "min" {
				return errors.New(m.localizer.Localize(
					languageId,
					"ErrorMinStringLength",
					map[string]interface{}{
						"StringName": fieldName,
						"MinLength":  r.Param(),
					},
				))
			}
			if r.Tag() == "max" {
				return errors.New(m.localizer.Localize(
					languageId,
					"ErrorMaxStringLength",
					map[string]interface{}{
						"StringName": fieldName,
						"MaxLength":  r.Param(),
					},
				))
			}
			if r.Tag() == "email" {
				return errors.New(m.localizer.Localize(languageId, "ErrorInvalidEmail", nil))
			}
			if r.Tag() == "phone" {
				return errors.New(m.localizer.Localize(languageId, "ErrorInvalidPhoneNumber", nil))
			}
			if r.Tag() == "longitude" {
				return errors.New(m.localizer.Localize(languageId, "ErrorInvalidLongLat", nil))
			}
			if r.Tag() == "latitude" {
				return errors.New(m.localizer.Localize(languageId, "ErrorInvalidLongLat", nil))
			}
			if r.Tag() == "eqfield" {
				return errors.New(m.localizer.Localize(
					languageId,
					"ErrorFieldNotEqual",
					map[string]interface{}{
						"FieldName":    fieldName,
						"EqualToField": r.Param(),
					},
				))
			}
		}

		return fmt.Errorf("tag validation failed. field name: %s. tag arg:%s", validationErrors[0].Field(), validationErrors[0].Tag())
	}

	return nil
}

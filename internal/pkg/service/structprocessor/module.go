package structprocessor

import (
	"context"
	"errors"
	"reflect"

	"clodeo.tech/public/go-universe/pkg/localization"
	pkgTagTransform "clodeo.tech/public/go-universe/pkg/tag/component/transform"
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type TagValidationErrorHandler = func(ctx context.Context, languageId string, validationErrors validator.ValidationErrors) error

type StructProcessorService interface {
	TransformByTag(obj interface{}) error
	ValidateByTag(ctx context.Context, obj interface{}) error
	TransformAndValidateByTag(ctx context.Context, obj interface{}) error
}

type StructProcessorServiceModule struct {
	localizer              localization.Localizer
	validatorValidate      *validator.Validate
	validationErrorHandler TagValidationErrorHandler
	transformFunc          func(obj interface{}) error
}

type StructProcessorServiceModuleOpts struct {
	Localizer              localization.Localizer
	ValidatorValidate      *validator.Validate
	ValidationErrorHandler TagValidationErrorHandler
	TransformFunc          func(obj interface{}) error
}

func NewStructProcessorService(opts StructProcessorServiceModuleOpts) StructProcessorService {
	if opts.ValidatorValidate == nil {
		opts.ValidatorValidate = validator.New()
	}
	if opts.TransformFunc == nil {
		opts.TransformFunc = pkgTagTransform.Transform
	}

	opts.ValidatorValidate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if uuid, ok := field.Interface().(uuid.UUID); ok {
			return uuid.String()
		}
		return nil
	}, uuid.UUID{})

	return &StructProcessorServiceModule{
		localizer:              opts.Localizer,
		transformFunc:          opts.TransformFunc,
		validatorValidate:      opts.ValidatorValidate,
		validationErrorHandler: opts.ValidationErrorHandler,
	}
}

func (m *StructProcessorServiceModule) TransformByTag(obj interface{}) error {
	return m.transformFunc(obj)
}

func (m *StructProcessorServiceModule) ValidateByTag(ctx context.Context, obj interface{}) error {
	langId := "id"
	err := m.validatorValidate.Struct(obj)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return err
		}

		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			if m.validationErrorHandler != nil {
				return m.validationErrorHandler(ctx, langId, validateErrs)
			}

			return m.DefaultTagValidationErrorHandler(ctx, langId, validateErrs)
		}
	}

	return nil
}

func (m *StructProcessorServiceModule) TransformAndValidateByTag(ctx context.Context, obj interface{}) error {
	err := m.TransformByTag(obj)
	if err != nil {
		return err
	}

	err = m.ValidateByTag(ctx, obj)
	if err != nil {
		return err
	}

	return nil
}

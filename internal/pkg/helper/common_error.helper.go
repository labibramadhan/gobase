package helper

import (
	"errors"

	"clodeo.tech/public/go-universe/pkg/localization"
)

func GetErrorDataNotFoundWithParam(localizer localization.Localizer, langId string, recordType string) error {
	lRecordType := localizer.Localize(langId, recordType, nil)
	return errors.New(localizer.Localize(langId, "ErrorRecordNotFoundWithParam",
		map[string]interface{}{
			"FieldName": lRecordType,
		},
	))
}

func GetErrorDateFieldValueIsFutureDate(localizer localization.Localizer, langId string, fieldName string) error {
	lFieldName := localizer.Localize(langId, fieldName, nil)
	return errors.New(localizer.Localize(langId, "ErrorInvalidDateValueUsingFutureDate",
		map[string]interface{}{
			"FieldName": lFieldName,
		},
	))
}

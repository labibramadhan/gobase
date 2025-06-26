package helper

import (
	"errors"
	"fmt"

	"clodeo.tech/public/go-universe/pkg/localization"
)

func UpdateDeleteWithRowsAffectedWrapper(localizer localization.Localizer, langId string, callFunc func() (int64, error)) error {
	rowsAffected, err := callFunc()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New(localizer.Localize(langId, "ErrorRecordNotFoundOrHasBeenChanged", nil))
	}

	return nil
}

func UpdateDeleteWithRowsAffectedWrapperWithoutLocalizer(callFunc func() (int64, error)) error {
	rowsAffected, err := callFunc()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s", "Data tidak ditemukan atau sudah diubah oleh pengguna lain")
	}

	return nil
}

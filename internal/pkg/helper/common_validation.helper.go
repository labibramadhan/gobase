package helper

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"clodeo.tech/public/go-universe/pkg/localization"
	"clodeo.tech/public/go-universe/pkg/util"
)

/*
Check if date value is valid (not a future date)
Skip if date is zero date
*/
func ValidateFieldDateValueWithFutureDate(localizer localization.Localizer, langId string, fieldName string, d time.Time, timezoneOffset int) error {
	if d.IsZero() {
		return nil
	}

	aDate := util.RemoveTime(d)

	currentLocalTime := time.Now().UTC().Add(time.Minute * time.Duration(timezoneOffset))
	currentDate := util.RemoveTime(currentLocalTime)

	if aDate.After(currentDate) {
		return GetErrorDateFieldValueIsFutureDate(localizer, langId, fieldName)
	}

	return nil
}

/*
Check if date of birth of date is greater or equal to minimum age
Min age in years
Skip if date is zero date
*/
func ValidateDateOfBirth(localizer localization.Localizer, langId string, d time.Time, minAge int, timezoneOffset int) error {
	if d.IsZero() {
		return nil
	}

	aDate := util.RemoveTime(d)

	currentLocalTime := time.Now().UTC().Add(time.Minute * time.Duration(timezoneOffset))
	currentDate := util.RemoveTime(currentLocalTime)

	var maxDate time.Time
	if minAge > 0 {
		maxDate = currentDate.AddDate(minAge, 0, 0)
	} else {
		maxDate = currentDate
	}

	if aDate.After(maxDate) {
		if minAge > 0 {
			return errors.New(localizer.Localize(langId, "ErrorDateOfBirthBelowMinAge",
				map[string]interface{}{
					"MinAge": minAge,
				},
			))
		}

		return GetErrorDateFieldValueIsFutureDate(localizer, langId, "DateOfBirth")
	}

	return nil
}

func IsValidCoordinates(coordinates string) bool {
	coords := strings.Split(coordinates, ",")
	if len(coords) != 2 {
		return false
	}
	coord1, err1 := strconv.ParseFloat(strings.TrimSpace(coords[0]), 64)
	coord2, err2 := strconv.ParseFloat(strings.TrimSpace(coords[1]), 64)
	if err1 != nil || err2 != nil || coord1 < -90 || coord1 > 90 || coord2 < -180 || coord2 > 180 {
		return false
	}
	return true
}

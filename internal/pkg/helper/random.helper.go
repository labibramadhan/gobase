package helper

import (
	"fmt"
	"time"

	pkgUtil "clodeo.tech/public/go-universe/pkg/util"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GenerateInvoiceNumber(date time.Time) (string, error) {
	dateLayout := "01" // MM (month)
	randomNumber, err := pkgUtil.GenerateSecureRandomNumberString(12)
	if err != nil {
		return "", err
	}

	// INV{mm}-{12 digit}
	return fmt.Sprintf("INV%s-%s", date.Format(dateLayout), randomNumber), nil
}

func GenerateOrderNumber(date time.Time) (string, error) {
	dateLayout := "01" // MM (month)
	randomNumber, err := pkgUtil.GenerateSecureRandomNumberString(12)
	if err != nil {
		return "", err
	}

	// ORDER{mm}-{12 digit}
	return fmt.Sprintf("ORDER%s-%s", date.Format(dateLayout), randomNumber), nil
}

func GenerateTransactionNumber(prefix string, date time.Time, counter int) (string, error) {
	dateLayout := "0601"
	return fmt.Sprintf("%s%s-%s", prefix, date.Format(dateLayout), fmt.Sprintf("%05d", counter)), nil
}

func GenerateRandomString(prefix ...string) (string, error) {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const totalLength = 15
	const defaultPrefix = ""

	// Determine the prefix to use
	usedPrefix := defaultPrefix
	if len(prefix) > 0 {
		usedPrefix = prefix[0]
	}

	// Calculate the length of the random part
	randomPartLength := totalLength - len(usedPrefix)
	if randomPartLength <= 0 {
		return "", fmt.Errorf("prefix is too long")
	}

	// Generate the random part
	randomPart, err := gonanoid.Generate(alphabet, randomPartLength)
	if err != nil {
		return "", err
	}

	// Combine the prefix with the random part
	return usedPrefix + randomPart, nil
}

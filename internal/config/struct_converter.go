package config

import (
	"time"

	pkgStructConverter "clodeo.tech/public/go-universe/pkg/struct_converter"
	"github.com/google/uuid"
)

func InitStructConverter() error {
	err := pkgStructConverter.RegisterNillableTypeAndValue(uuid.Nil)
	if err != nil {
		return err
	}

	err = pkgStructConverter.RegisterNillableTypeAndValue(time.Time{})
	if err != nil {
		return err
	}

	return nil
}

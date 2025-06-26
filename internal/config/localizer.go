package config

import (
	"fmt"

	pkgLocale "clodeo.tech/public/go-universe/pkg/localization"
)

func InitLocalizer(basePath string) (localizer pkgLocale.Localizer, err error) {
	localizer = pkgLocale.NewLocalizer()
	localizerRes := map[string]string{
		"en": fmt.Sprintf("%s/resource/locale/en.toml", basePath),
		"id": fmt.Sprintf("%s/resource/locale/id.toml", basePath),
	}

	if err := localizer.Init(localizerRes); err != nil {
		return nil, err
	}

	return localizer, nil
}

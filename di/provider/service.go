package provider

import (
	"clodeo.tech/public/go-universe/pkg/localization"
	"github.com/google/wire"

	"gobase/internal/pkg/service/structprocessor"
)

var ServiceSet = wire.NewSet(
	ProvideServiceStructProcessorService,
)

func ProvideServiceStructProcessorService(localizer localization.Localizer) structprocessor.StructProcessorService {
	return structprocessor.NewStructProcessorService(structprocessor.StructProcessorServiceModuleOpts{
		Localizer: localizer,
	})
}

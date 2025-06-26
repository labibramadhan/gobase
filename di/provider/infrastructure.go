package provider

import (
	"context"
	"path/filepath"
	"runtime"

	"clodeo.tech/public/go-universe/pkg/env"
	"clodeo.tech/public/go-universe/pkg/localization"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"

	"gobase/config"
	iconfig "gobase/internal/config"
	masterdataentity "gobase/internal/db/masterdata/entity"
	"gobase/internal/pkg/helper/excel"
	"gobase/internal/pkg/helper/excel/excelize"

	_ "github.com/lib/pq" // used for sql queries
)

var InfrastructureSet = wire.NewSet(
	ProvideInfrastructureLocalizer,
	ProvideInfrastructureBun,
	ProvideInfrastructureExcelManager,
)

func ProvideInfrastructureLocalizer() localization.Localizer {
	_, fn, _, _ := runtime.Caller(1)
	basePath := filepath.Dir(filepath.Dir(filepath.Dir(fn)))

	localizer, err := iconfig.InitLocalizer(basePath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initiate localizer")
	}

	return localizer
}

func ProvideInfrastructureBun(cfg *config.MainConfig) *bun.DB {
	var db *bun.DB

	switch cfg.Rdbms.App.Driver {
	case "postgres":
		config, err := pgx.ParseConfig(cfg.Rdbms.App.DSN)
		if err != nil {
			panic(err)
		}
		config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

		sqldb := stdlib.OpenDB(*config)
		db = bun.NewDB(sqldb, pgdialect.New())
	default:
		log.Fatal().Msgf("failed to initiate RDBMS : %s", cfg.Rdbms.App.Driver)
	}

	// enable sql logging
	if env.GetEnvironmentName() == "local" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	//nolint:gosec,errcheck
	if env.GetEnvironmentName() == "local" {
		db.NewDropTable().Model(&masterdataentity.Product{}).Exec(context.Background())
		db.NewCreateTable().Model(&masterdataentity.Product{}).Exec(context.Background())
		db.NewDropTable().Model(&masterdataentity.ProductVariant{}).Exec(context.Background())
		db.NewCreateTable().Model(&masterdataentity.ProductVariant{}).Exec(context.Background())
		db.NewDropTable().Model(&masterdataentity.ProductAttribute{}).Exec(context.Background())
		db.NewCreateTable().Model(&masterdataentity.ProductAttribute{}).Exec(context.Background())
		db.NewDropTable().Model(&masterdataentity.RelProductVariantProductAttribute{}).Exec(context.Background())
		db.NewCreateTable().Model(&masterdataentity.RelProductVariantProductAttribute{}).Exec(context.Background())
	}

	return db
}

func ProvideInfrastructureExcelManager() excel.Excel {
	return excelize.NewExcel()
}

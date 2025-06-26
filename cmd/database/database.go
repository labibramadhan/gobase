package database

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"

	outbox "clodeo.tech/public/go-outbox/event_outbox"
	pkgConfig "clodeo.tech/public/go-universe/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog/log"
	dsnParser "github.com/sidmal/dsn-parser"
	"github.com/spf13/cobra"

	"gobase/config"

	_ "github.com/golang-migrate/migrate/v4/database/postgres" // import postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var MigrationCommand = &cobra.Command{
	Use:   "db",
	Short: "Database Migration",
	Run: func(c *cobra.Command, args []string) {
		DatabaseMigration()
	},
}

var (
	flags      = flag.NewFlagSet("db", flag.ExitOnError)
	configPath = flags.String("config", "config/file", "Config URL dir i.e. config/file")
	dir        = flags.String("dir", "migration", "directory with migration files")
)

var (
	usageCommands = `
Commands:
  up [N]?              Migrate all or N up migrations app db
  outbox [N]?          Migrate outbox or N up migrations outbox db
  goto [V]             Migrate the app DB to a specific version
  down [N]?            Down all or N down migrations app db

For more features, use https://github.com/golang-migrate/migrate/tree/master/cmd/migrate`
)

func initLogger() {
	// initiate logger
}

//nolint:gocognit,gocyclo,revive
func DatabaseMigration() {
	initLogger()

	cfg := &config.MainConfig{}
	err := pkgConfig.ReadConfig(cfg, *configPath, "config")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	flags.Usage = usage
	err = flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	args := flags.Args()

	if len(args) == 0 {
		flags.Usage()
		return
	}

	var (
		m       *migrate.Migrate
		connstr string
	)

	connstr = cfg.DBMigration.App.DSN

	m, err = migrate.New(
		fmt.Sprintf("file://internal/db/masterdata/%s", *dir),
		connstr,
	)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	switch args[0] {
	case "up":
		if len(args) == 2 {
			step, err := strconv.Atoi(args[1])
			if err != nil {
				log.Fatal().Err(err).Msg(err.Error())
			}
			err = m.Steps(step)
			if err != nil {
				log.Fatal().Err(err).Msg(err.Error())
			}
		} else {
			err := m.Up()
			if err != nil {
				log.Fatal().Err(err).Msg(err.Error())
			}
		}
		return
	case "outbox":
		outboxConnstr := cfg.DBMigration.Outbox.DSN
		parsedDsn, err := dsnParser.New(outboxConnstr)
		if err != nil {
			log.Fatal().Err(err).Msg(err.Error())
		}
		err = outbox.RunDBMigration(context.Background(), outboxConnstr, parsedDsn.Database)
		if err != nil {
			log.Fatal().Err(err).Msg(err.Error())
		}
		return
	case "down":
		if len(args) == 2 {
			step, err := strconv.Atoi(args[1])
			if err != nil {
				log.Fatal().Err(err).Msg(err.Error())
			}
			err = m.Steps(step * -1)
			if err != nil {
				log.Fatal().Err(err).Msg(err.Error())
			}
		} else {
			err := m.Down()
			if err != nil {
				log.Fatal().Err(err).Msg(err.Error())
			}
		}
		return
	case "goto":
		if len(args) == 2 {
			step, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				log.Fatal().Err(err).Msg(err.Error())
			}
			err = m.Migrate(uint(step))
			if err != nil {
				log.Fatal().Err(err).Msg(err.Error())
			}
		} else {
			usage()
		}
		return
	}

	err = m.Up()
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			log.Fatal().Err(sourceErr).Msg(sourceErr.Error())
		}
		if dbErr != nil {
			log.Fatal().Err(dbErr).Msg(dbErr.Error())
		}
	}()
}

func usage() {
	_, err := fmt.Println(usageCommands)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}
}

package config

type MainConfig struct {
	ServiceName string            `fig:"serviceName"`
	General     GeneralConfig     `fig:"general"`
	Server      ServerConfig      `fig:"server"`
	DBMigration DBMigrationConfig `fig:"dbmigrate"`
	Rdbms       RdbmsConfig       `fig:"rdbms"`
	Tracer      TracerConfig      `fig:"tracer"`
}

type (
	GeneralConfig struct {
		TZ string `json:"tz"`
	}

	ServerConfig struct {
		BaseUrl string `fig:"baseUrl"`
		Grpc    struct {
			Port int `fig:"port"`
		}

		Rest ServerConfigRest `fig:"rest"`
	}

	ServerConfigRest struct {
		ListenAddress      string `fig:"listenAddress"`
		Port               int    `fig:"port"`
		BodyLimitMB        int    `fig:"bodyLimitMB"`
		ReadTimeoutSecond  int    `fig:"readTimeoutSecond"`
		WriteTimeoutSecond int    `fig:"writeTimeoutSecond"`
	}

	DBMigrationConfig struct {
		App    MigrationConfig `fig:"app"`
		Outbox MigrationConfig `fig:"outbox"`
	}

	MigrationConfig struct {
		Driver string `fig:"driver"`
		DSN    string `fig:"dsn"`
	}

	RdbmsConfig struct {
		App    DBConfig `fig:"app"`
		Outbox DBConfig `fig:"outbox"`
	}

	DBConfig struct {
		Driver          string `fig:"driver"`
		DSN             string `fig:"dsn"`
		MaxOpenConns    int    `fig:"maxOpenConns"`
		MaxIdleConns    int    `fig:"maxIdleConns"`
		ConnMaxLifetime int    `fig:"connMaxLifetime"`
		Retry           int    `fig:"retry"`
	}

	TracerConfig struct {
		Enabled  bool         `fig:"enabled"`
		Provider string       `fig:"provider"`
		Jaeger   JaegerConfig `fig:"jaeger"`
	}

	JaegerConfig struct {
		CollectorUrl string `fig:"collectorUrl"`
	}
)

package godatasource

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	golog "log"
	"os"
	"time"
)

const (
	DefaultDatasourceName = "DEFAULT"
)

const (
	DialectPostgres = "postgres"
)

var (
	datasources = make(map[string]*gorm.DB)
)

func InitDatasources(filepath string) {
	var loggingWriter = initLogger()
	var sourceProfiles = loadDatasource(filepath)
	for sourceName, profile := range sourceProfiles.Profiles {
		var datasource *gorm.DB
		var err error
		switch profile.Gorm.Dialect {
		case DialectPostgres:
			datasource, err = gorm.Open(postgres.New(postgres.Config{DSN: fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
				profile.Datasource.Host, profile.Datasource.Port, profile.Datasource.Dbname,
				profile.Datasource.User, profile.Datasource.Password,
			)}), &gorm.Config{
				Logger:                 loggingWriter,
			})
		default:
			log.Fatal().Err(err).Msgf("Not support dialect %s", profile.Gorm.Dialect)
		}

		if err != nil {
			log.Fatal().Err(err).Msgf("Got error when opening connection to datasource %s", sourceName)
		}

		datasources[sourceName] = datasource

		log.Info().Str("datasource", sourceName).
			Str("gorm_dialect", profile.Gorm.Dialect).
			Msg("Database connection created")
	}
}

func initLogger() logger.Interface {
	if loggingWriter == nil {
		log.Info().Msg("No external logging writer provided. Use std Go log")
		loggingWriter = golog.New(os.Stdout, "\r", golog.LstdFlags|golog.Lshortfile|golog.Lmicroseconds)
	}
	return logger.New(
		loggingWriter,
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)
}

func loadDatasource(filepath string) DataSourceConfig {
	log.Info().Msgf("Reading datasource profiles...")
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error loading datasource profiles file")
	}

	var config DataSourceConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error parsing datasource profiles file")
	}

	return config
}

func Connection(sourceName string) *gorm.DB {
	if datasource, ok := datasources[sourceName]; ok {
		return datasource
	}
	log.Error().Msgf("Not found connection of datasource %s", sourceName)
	return nil
}

func DefaultConnection() *gorm.DB {
	return Connection(DefaultDatasourceName)
}

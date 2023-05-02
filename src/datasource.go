package datasource

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	db  *gorm.DB
	err error
)

func InitDatasource(filepath string) {
	var sourceConfig, dialect, autoMigrate = loadDatasource(filepath)
	db, err = gorm.Open(dialect, fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		sourceConfig.Host, sourceConfig.Port, sourceConfig.Dbname,
		sourceConfig.User, sourceConfig.Password,
	))
	if err != nil {
		log.Fatal().Err(err).Msg("Got error when opening connection")
	}
	log.Info().Str("gorm_dialect", dialect).Msg("Database connection created.")

	if autoMigrate {
		log.Warn().Msg("Not supported Migrate feature.")
	}
}

func loadDatasource(filepath string) (SourceConfig, string, bool) {
	log.Info().Msgf("Read datasource properties from %s", filepath)
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error loading file %s", filepath)
	}

	var config map[string]interface{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error parsing file %s", filepath)
	}

	var sourceConfig = SourceConfig{
		Host:     config["datasource"].(map[string]interface{})["host"].(string),
		Port:     config["datasource"].(map[string]interface{})["port"].(int),
		User:     config["datasource"].(map[string]interface{})["user"].(string),
		Dbname:   config["datasource"].(map[string]interface{})["dbname"].(string),
		Password: config["datasource"].(map[string]interface{})["password"].(string),
	}
	var autoMigrate = config["gorm"].(map[string]interface{})["auto_migrate"].(bool)
	return sourceConfig, config["gorm"].(map[string]interface{})["dialect"].(string), autoMigrate
}

func Connection() *gorm.DB {
	return db
}

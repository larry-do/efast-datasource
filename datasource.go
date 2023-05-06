package godatasource

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	DefaultDatasourceName = "DEFAULT"
)

var (
	datasources = make(map[string]*gorm.DB)
)

func InitDatasources(filepath string) {
	var sourceProfiles = loadDatasource(filepath)
	for sourceName, profile := range sourceProfiles {
		datasource, err := gorm.Open(profile.Dialect, fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			profile.Host, profile.Port, profile.Dbname,
			profile.User, profile.Password,
		))
		if err != nil {
			log.Fatal().Err(err).Msgf("Got error when opening connection to datasource %s.", sourceName)
		}

		datasource.LogMode(profile.PrintLog)

		datasources[sourceName] = datasource

		log.Info().Str("datasource", sourceName).
			Str("gorm_dialect", profile.Dialect).
			Msg("Database connection created.")
	}
}

func loadDatasource(filepath string) map[string]SourceProfile {
	log.Info().Msgf("Read datasource profiles from %s.", filepath)
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error loading file %s.", filepath)
	}

	var config map[string]interface{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error parsing file %s.", filepath)
	}

	sourceProfiles := make(map[string]SourceProfile)

	for sourceName, profile := range config["profiles"].(map[string]interface{}) {
		connectionProfile := profile.(map[string]interface{})["datasource"]
		gormConfig := profile.(map[string]interface{})["gorm"]

		var sourceProfile = SourceProfile{
			Host:     connectionProfile.(map[string]interface{})["host"].(string),
			Port:     connectionProfile.(map[string]interface{})["port"].(int),
			User:     connectionProfile.(map[string]interface{})["user"].(string),
			Dbname:   connectionProfile.(map[string]interface{})["dbname"].(string),
			Password: connectionProfile.(map[string]interface{})["password"].(string),
			Dialect:  gormConfig.(map[string]interface{})["dialect"].(string),
			PrintLog: gormConfig.(map[string]interface{})["print_log"].(bool),
		}

		sourceProfiles[sourceName] = sourceProfile
	}

	return sourceProfiles
}

func Connection(sourceName string) *gorm.DB {
	return datasources[sourceName]
}

func DefaultConnection() *gorm.DB {
	return Connection(DefaultDatasourceName)
}

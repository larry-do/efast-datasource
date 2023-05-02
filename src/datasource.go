package datasource

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

var (
	db  *gorm.DB
	err error
)

func InitDatasource() {
	InitDatasourceWithFile("./resources/datasource.properties")
}

func InitDatasourceWithFile(filename string) {
	var sourceConfig, dialect, autoMigrate = loadDatasource(filename)

	db, err = gorm.Open(dialect, fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		sourceConfig.Host, sourceConfig.Port, sourceConfig.Dbname,
		sourceConfig.User, sourceConfig.Password,
	))
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Str("gorm_dialect", dialect).Msg("Database connection created.")

	if autoMigrate {
		log.Warn().Msg("Not supported Migrate feature.")
		/*migrateToDatabase()*/
	}
}

/*func migrateToDatabase() {
	db.AutoMigrate(model.Channel{})
	db.AutoMigrate(model.Audio{})
	db.AutoMigrate(model.Program{})
	db.Model(model.Program{}).AddForeignKey("channel_id", "channel(id)", "RESTRICT", "RESTRICT")
	db.Model(model.Program{}).AddForeignKey("audio_id", "audio(id)", "RESTRICT", "RESTRICT")

	log.Println("Migrated database.")
}*/

func loadDatasource(filename string) (SourceConfig, string, bool) {
	log.Printf("Read datasource properties from %s\n", filename)
	var err = godotenv.Load(filename)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error loading file %s", filename)
	}
	var datasource = SourceConfig{}
	datasource.Host = os.Getenv("datasource.host")
	datasource.Port, _ = strconv.Atoi(os.Getenv("datasource.port"))
	datasource.Dbname = os.Getenv("datasource.dbname")
	datasource.User = os.Getenv("database.user")
	datasource.Password = os.Getenv("database.password")
	var autoMigrate, _ = strconv.ParseBool(os.Getenv("gorm.auto_migrate"))
	return datasource, os.Getenv("gorm.dialect"), autoMigrate
}

func Connection() *gorm.DB {
	return db
}

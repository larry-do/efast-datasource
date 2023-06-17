package godatasource

type DataSourceConfig struct {
	Profiles map[string]Profile `yaml:"profiles"`
}

type Profile struct {
	Datasource Datasource `yaml:"datasource"`
	Gorm       GormProps  `yaml:"gorm"`
}

type Datasource struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Dbname   string `yaml:"dbname"`
	Password string `yaml:"password"`
}

type GormProps struct {
	Dialect                string `yaml:"dialect"`
	PrintLog               bool   `yaml:"print_log"`
	SkipDefaultTransaction bool   `yaml:"skip_default_transaction"`
}

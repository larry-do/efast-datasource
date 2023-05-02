package main

import (
	datasource "efast-datasource/src"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	datasource.InitDatasource("./resources/datasource.yaml")
}

package common_dao

import (
	"godatasource"
	"gorm.io/gorm"
)

func Save(object any) (tx *gorm.DB) {
	return godatasource.DefaultConnection().Create(object)
}

func Update(object any) (tx *gorm.DB) {
	return godatasource.DefaultConnection().Save(object)
}

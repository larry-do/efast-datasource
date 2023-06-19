package dao

import (
	"gorm.io/gorm"
)

type IDao[E interface{}] interface {
	Save(E) (tx *gorm.DB)
	Update(E) (tx *gorm.DB)
}

type AbstractDao[E interface{}] struct {
	IDao[E]
}

func (dao *AbstractDao[E]) Update(tx *gorm.DB, entity interface{}) *gorm.DB {
	return tx.Save(entity)
}

func (dao *AbstractDao[E]) Save(tx *gorm.DB, entity interface{}) *gorm.DB {
	return tx.Create(entity)
}

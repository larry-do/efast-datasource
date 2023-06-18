package dao

import (
	"godatasource"
	"gorm.io/gorm"
)

type IDao[E interface{}] interface {
	Save(E) (tx *gorm.DB)
	Update(E) (tx *gorm.DB)
}

type TxDao[E interface{}] struct {
	IDao[E]
	tx *gorm.DB
}

func (dao *TxDao[E]) datasource() string {
	return godatasource.DefaultDatasourceName
}

func (dao *TxDao[E]) WithTx(tx *gorm.DB) *TxDao[E] {
	dao.tx = tx
	return dao
}

func (dao *TxDao[E]) Tx() (tx *gorm.DB) {
	if dao.tx == nil {
		dao.tx = godatasource.Connection(dao.datasource())
	}
	return dao.tx
}

type AbstractDao[E interface{}] struct {
	TxDao[E]
}

func (dao *AbstractDao[E]) Update(entity interface{}) (tx *gorm.DB) {
	return dao.Tx().Save(entity)
}

func (dao *AbstractDao[E]) Save(entity interface{}) (tx *gorm.DB) {
	return dao.Tx().Create(entity)
}

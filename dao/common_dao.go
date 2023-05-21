package common_dao

import "godatasource"

func Save(object any) {
	godatasource.DefaultConnection().Create(object)
}

func Update(object any) {
	godatasource.DefaultConnection().Save(object)
}

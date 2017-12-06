package v1beta1

import (
	"server/pkg/storage/mysql"
)

//Insert insert storage to db
func (storage *Storage) Insert() error {
	return mysql.GetDB().Create(storage).Error
}

//GetByID get storage by id
func (storage *Storage) GetByID() (*Storage, error) {
	err := mysql.GetDB().First(storage, storage.ID).Error
	return storage, err
}

//GetAll get all storage
func (storage *Storage) GetAll() ([]*Storage, error) {
	var storages []*Storage
	err := mysql.GetDB().Find(&storages).Error
	return storages, err
}

//DeleteByName delete storage by name
func (storage *Storage) DeleteByName() error {
	return mysql.GetDB().Delete(Storage{}, "name=?", storage.Name).Error
}

//GetByNamespace get  storage by namespace
func (storage *Storage) GetByNamespace() ([]*Storage, error) {
	var storages []*Storage
	err := mysql.GetDB().Find(&storages, "namespace=?", storage.Namespace).Error
	return storages, err
}

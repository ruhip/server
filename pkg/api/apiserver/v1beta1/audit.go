package apiserver

import (
	"server/pkg/storage/mysql"
)

// Insert insert audit
func (audit *Audit) Insert() error {
	return mysql.GetDB().Create(audit).Error
}

// GetAll get all audit
func (audit *Audit) GetAll() ([]*Audit, error) {
	var audits []*Audit
	err := mysql.GetDB().Find(&audits).Error
	return audits, err
}

// Delete delete audit
func (audit *Audit) Delete() error {
	return mysql.GetDB().Delete(audit).Error
}

// Get get all audit
func (audit *Audit) Get() (*Audit, error) {
	var tk *Audit
	err := mysql.GetDB().Where("id=?", audit.ID).Find(tk).Error
	return tk, err
}

// Update update audit
func (audit *Audit) Update() error {
	return mysql.GetDB().Update(audit).Error
}

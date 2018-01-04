package v1beta1

import (
	"server/pkg/storage/mysql"
)

//GetAll get all app
func (pc *PodLifeCycle) GetAll() ([]*PodLifeCycle, error) {
	var pcs []*PodLifeCycle
	err := mysql.GetDB().Find(&pcs).Error
	return pcs, err
}

// GetByID get  PodLifeCycle  by clusterID
func (pc *PodLifeCycle) GetByID(ID string) (*PodLifeCycle, error) {
	err := mysql.GetDB().Where("cluster_id=?", pc.ClusterID).Find(pc).Error
	return pc, err
}

// Insert insert PodLifeCycle to db
func (pc *PodLifeCycle) Insert() (*PodLifeCycle, error) {
	err := mysql.GetDB().Create(pc).Error
	return pc, err
}

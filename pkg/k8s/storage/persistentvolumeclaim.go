// Copyright © 2017 huang jia <449264675@qq.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//PersistentVolumeClaimInterface has methods to work with PersistentVolumeClaim resources.
type PersistentVolumeClaimInterface interface {
	CreatePersistentVolumeClaim(persistentVolumeClaim *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error)
	UpdatePersistentVolumeClaim(persistentVolumeClaim *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error)
	DeletePersistentVolumeClaim(name, namespace string) error
	GetPersistentVolumeClaim(name, namespace string) (*v1.PersistentVolumeClaim, error)
	ListPersistentVolumeClaim(labels, namespace string) ([]v1.PersistentVolumeClaim, error)
}

//persistentVolumeClaims implements PersistentVolumeClaimInterface.
type persistentVolumeClaims struct {
	*kubernetes.Clientset
}

//NewPersistentVolumeClaims return a PersistentVolumeClaim.
func NewPersistentVolumeClaims(client *kubernetes.Clientset) PersistentVolumeClaimInterface {
	return &persistentVolumeClaims{Clientset: client}
}

func (client *persistentVolumeClaims) CreatePersistentVolumeClaim(persistentVolumeClaim *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	return client.CoreV1().PersistentVolumeClaims(persistentVolumeClaim.Namespace).Create(persistentVolumeClaim)
}

func (client *persistentVolumeClaims) UpdatePersistentVolumeClaim(persistentVolumeClaim *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	return client.CoreV1().PersistentVolumeClaims(persistentVolumeClaim.Namespace).Update(persistentVolumeClaim)
}

func (client *persistentVolumeClaims) DeletePersistentVolumeClaim(name, namespace string) error {
	return client.CoreV1().PersistentVolumeClaims(namespace).Delete(name, &metav1.DeleteOptions{})
}

func (client *persistentVolumeClaims) GetPersistentVolumeClaim(name, namespace string) (*v1.PersistentVolumeClaim, error) {
	return client.CoreV1().PersistentVolumeClaims(namespace).Get(name, metav1.GetOptions{})
}

func (client *persistentVolumeClaims) ListPersistentVolumeClaim(labels string, namespace string) ([]v1.PersistentVolumeClaim, error) {
	listOption := metav1.ListOptions{}
	if labels != "" {
		listOption.LabelSelector = labels
	}
	list, err := client.CoreV1().PersistentVolumeClaims(namespace).List(listOption)
	if err != nil {
		return []v1.PersistentVolumeClaim{}, err
	}
	return list.Items, err
}
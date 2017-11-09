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

package metric

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

//EventInterface has methods to work with Event resources.
type EventInterface interface {
	GetEvents(namespace string) ([]v1.Event, error)
}

//events implements HPAInterface.
type events struct {
	*kubernetes.Clientset
}

//Newevents return a events.
func NewEvents(client *kubernetes.Clientset) EventInterface {
	return &events{Clientset: client}
}

func (client *events) GetEvents(namespace string) ([]v1.Event, error) {
	list, err := client.CoreV1().Events(namespace).List(metav1.ListOptions{
		LabelSelector: labels.Everything().String(),
		FieldSelector: fields.Everything().String(),
	})
	if err != nil {
		return []v1.Event{}, err
	}
	return list.Items, nil
}
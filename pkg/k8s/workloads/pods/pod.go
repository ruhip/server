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

package pods

import (
	"io/ioutil"

	"server/pkg/utils/log"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

//PodInterface has methods to work with Pod resources.
type PodInterface interface {
	DeletePod(pod v1.Pod) error
	GetPod(name, namesapce string) (*v1.Pod, error)
	ListPods(namespace string) ([]v1.Pod, error)
	ListPodByDeploymentName(name, namespace string) ([]v1.Pod, error)
	GetPodLogs(name, namespace string, logOptions *v1.PodLogOptions) (string, error)
}

//pods implements PodInterface.
type pods struct {
	*kubernetes.Clientset
}

//NewPods return pods.
func NewPods(client *kubernetes.Clientset) PodInterface {
	return &pods{Clientset: client}
}

func (client *pods) DeletePod(pod v1.Pod) error {
DELETE_POD:
	deletePropagationForeground := new(metav1.DeletionPropagation)
	*deletePropagationForeground = metav1.DeletePropagationForeground
	if err := client.CoreV1().Pods(pod.Namespace).Delete(pod.Name, &metav1.DeleteOptions{PropagationPolicy: deletePropagationForeground}); err != nil {
		if errors.IsConflict(err) {
			goto DELETE_POD
		}
		return err
	}
	return nil
}

func (client *pods) GetPod(name, namesapce string) (*v1.Pod, error) {
	return client.CoreV1().Pods(namesapce).Get(name, metav1.GetOptions{})
}

func (client *pods) ListPods(namespace string) ([]v1.Pod, error) {
	list, err := client.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return []v1.Pod{}, err
	}
	return list.Items, nil
}

func (client *pods) ListPodByDeploymentName(name, namespace string) ([]v1.Pod, error) {
	list, err := client.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: "minipaas.io/name=" + name})
	if err != nil {
		return []v1.Pod{}, err
	}
	return list.Items, nil
}

func (client *pods) GetPodLogs(name, namespace string, logOptions *v1.PodLogOptions) (string, error) {
	req := client.CoreV1().RESTClient().Get().
		Namespace(namespace).
		Name(name).
		Resource("pods").
		SubResource("log").
		VersionedParams(logOptions, scheme.ParameterCodec)

	readCloser, err := req.Stream()
	if err != nil {
		return err.Error(), nil
	}

	defer func() {
		if err = readCloser.Close(); err != nil {
			log.Error("close readstream err:%v", err)
		}
	}()

	result, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

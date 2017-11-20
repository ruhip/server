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

package client

import (
	"os"
	"path/filepath"

	api "server/pkg/api/apiserver/v1beta1"
	"server/pkg/k8s/configuration"
	"server/pkg/k8s/metric"
	"server/pkg/k8s/service"
	"server/pkg/k8s/storage"
	"server/pkg/k8s/workloads/controllers"
	"server/pkg/k8s/workloads/pods"
	"server/pkg/utils/log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var fakes = make(map[string]*clientset)

//Clientset contains the clients for groups. Each group has exactly one
type clientset struct {
	*kubernetes.Clientset
}

//Init init all k8s cluster client
func Init() {
	clusters, err := new(api.Cluster).GetAll()
	if err != nil {
		log.Critical("query cluster err: %v", err)
	}
	config := filepath.Join(homeDir(), "config")
	file, err := os.Create(config)
	if err != nil {
		log.Critical("create k8s config file err: %v", err)
	}

	for k := range clusters {
		if err = file.Truncate(0); err != nil {
			log.Critical("clean k8s tmp config file err: %v", err)
		}
		if _, err = file.Write([]byte(clusters[k].ConfigContent)); err != nil {
			log.Critical("write k8s tmp config file err: %v", err)
		}
		fake, err := newClientset(config)
		if err != nil {
			log.Critical("create k8s client err: %v, who's clusterID is %q", err, clusters[k].ClusterID)
		}
		fakes[clusters[k].ClusterID] = fake
	}
}

func newClientset(config string) (*clientset, error) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		return nil, err
	}
	cs, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	return &clientset{Clientset: cs}, nil
}

//GetClientset return clientset
func GetClientset(clusterID string) *clientset {
	if fake, ok := fakes[clusterID]; ok {
		return fake
	}
	return nil
}

// AddClientset add clientset by clusterID
func AddClientset(cluster *api.Cluster) error {
	config := filepath.Join(homeDir(), "config")
	file, err := os.Create(config)
	if err != nil {
		log.Critical("create k8s config file err: %v", err)
		return err
	}
	if err = file.Truncate(0); err != nil {
		log.Critical("clean k8s tmp config file err: %v", err)
		return err
	}
	if _, err = file.Write([]byte(cluster.ConfigContent)); err != nil {
		log.Critical("write k8s tmp config file err: %v", err)
		return err
	}
	fake, err := newClientset(config)
	if err != nil {
		log.Critical("create k8s client err: %v, who's clusterID is %q", err, cluster.ClusterID)
		return err
	}
	fakes[cluster.ClusterID] = fake
	return nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func (c *clientset) Deployments() controllers.DeploymentInterface {
	return controllers.NewDeployments(c.Clientset)
}

func (c *clientset) HPAs() controllers.HPAInterface {
	return controllers.NewHpas(c.Clientset)
}

func (c *clientset) Statefulsets() controllers.StatefulsetInterface {
	return controllers.NewStatefulsets(c.Clientset)
}

func (c *clientset) Pods() pods.PodInterface {
	return pods.NewPods(c.Clientset)
}

func (c *clientset) Services() service.ServiceInterface {
	return service.NewServices(c.Clientset)
}

func (c *clientset) ConfigMaps() configuration.ConfigMapInterface {
	return configuration.NewconfigMaps(c.Clientset)
}

func (c *clientset) Events() metric.EventInterface {
	return metric.NewEvents(c.Clientset)
}

func (c *clientset) PersistentVolumeClaims() storage.PersistentVolumeClaimInterface {
	return storage.NewPersistentVolumeClaims(c.Clientset)
}

func (c *clientset) StorageClasses() storage.StorageClassInterface {
	return storage.NewStorageClasses(c.Clientset)
}

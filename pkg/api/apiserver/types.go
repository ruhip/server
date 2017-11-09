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

package apiserver

import (
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Cluster the k8s cluster info
type Cluster struct {
	ClusterID     string `json:"cluster_id" gorm:"primary_key"`
	Name          string `json:"name"`
	Describe      string `json:"describe"`
	ConfigContent string `json:"config_content"`
}

//App the application that use deploy
type App struct {
	ID            uint       `json:"id"`
	CreatedAt     time.Time  `json:"createAt"`
	Name          string     `json:"name"`
	UserName      string     `json:"nameSpace"`
	Description   string     `json:"description"`
	AppStatus     int        `json:"appStatus"`
	ServiceCount  int        `json:"serviceCount"`
	InstanceCount int        `json:"intanceCount"`
	External      string     `json:"external"`
	Items         []*Service `json:"services" gorm:"-"`
}

//Service the app's service
type Service struct {
	AppName       string           `json:"appName"`
	Name          string           `json:"name"`
	Image         string           `json:"image"`
	InstanceCount int              `json:"instanceCount" `
	Status        int              `json:"status"`
	Type          int              `json:"type"` //0 stateless 1 statefulset
	NodeName      string           `json:"nodeName"`
	External      string           `json:"external"`
	LoadbalanceIP string           `json:"loadbalanceIP"`
	CPU           string           `json:"cpu"`
	Memory        string           `json:"memory"`
	Ports         []v1.ServicePort `json:"ports"`
	Envs          []v1.EnvVar      `json:"envs"`
	Cmds          []string         `json:"cmds"`
	Volumes       []Volume         `json:"volumes"`
	Storage       *Storage         `json:"storage"`
}

//Volume service volume
type Volume struct {
	Type          int      `json:"type"` //0 挂载整个配置组（目录的方式挂载）1 挂载配置组中的一个key（挂载配置组中的单个文件）
	MountPath     string   `json:"mountPath"`
	ConfigMapName string   `json:"configMapName"`
	ConfigMapKey  []string `json:"configMapKey"`
}

//Storage rdb storage
type Storage struct {
	Type       string `json:"Type"`
	Size       string `json:"Size"`
	AccessMode string `json:"accessMode"`
}

//Config config group
type Config struct {
	Name string            `json:"name"`
	Data map[string]string `json:"data"`
}

//Event the resource event
type Event struct {
	Reason        string      `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`
	Message       string      `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`
	LastTimestamp metav1.Time `json:"lastTimestamp,omitempty" protobuf:"bytes,7,opt,name=lastTimestamp"`
	Type          string      `json:"type,omitempty" protobuf:"bytes,9,opt,name=type"`
}

//Process container's process
type Process struct {
	User        string  `json:"user"`
	PID         int64   `json:"pid"`
	ParentPID   int64   `json:"parent_pid"`
	StartTime   string  `json:"start_time"`
	PercentCPU  float64 `json:"percent_cpu"`
	PercentMEM  float64 `json:"percent_mem"`
	RSS         int64   `json:"rss"`
	VirtualSize int64   `json:"virtual_size"`
	Status      string  `json:"status"`
	RunningTime string  `json:"running_time"`
	CgroupPath  string  `json:"cgroup_path"`
	Cmd         string  `json:"cmd"`
}

//Item Cadvisor api result item
type Item struct {
	Name string `json:"name"`
}

//CadvisorResult request Cadvisor api result
type CadvisorResult struct {
	Subcontainers []*Item `json:"subcontainers"`
}

//HPA k8s rsource hpa's dto
type HPA struct {
	Name                           string `json:"name"`
	RefObjectName                  string `json:"refObject"`
	MinReplicas                    *int32 `json:"minReplicas"`
	MaxReplicas                    int32  `json:"maxReplicas"`
	TargetCPUUtilizationPercentage *int32 `json:"targetCPUUtilizationPercentage"`
}

//CephRBD ceph rbd's info
type CephRBD struct {
	Provisioner          string `json:"provisioner"`
	Name                 string `json:"name"`
	Monitors             string `json:"monitors"`
	AdminID              string `json:"adminId"`
	AdminSecretName      string `json:"adminSecretName"`
	AdminSecretNamespace string `json:"adminSecretNamespace"`
	Pool                 string `json:"pool"`
	UserID               string `json:"userId"`
	UserSecretName       string `json:"userSecretName"`
	FsType               string `json:"fsType"`
	ImageFormat          string `json:"imageFormat"`
	ImageFeatures        string `json:"imageFeatures"`
}

// TickScaleTask ticker time scale task
type TickScaleTask struct {
	ID                             uint      `json:"id"`
	Name                           string    `json:"name"`
	Namespace                      string    `json:"namespace"`
	Spec                           string    `json:"spec"`
	Desired                        int32     `json:"desired"`
	ClusterID                      string    `json:"clusterID"`
	Status                         uint      `json:"status"` //0 off 1 on
	CreationTime                   time.Time `json:"creationTime"`
	RefObjectName                  string    `json:"refObject"`
	MinReplicas                    *int32    `json:"minReplicas"`
	MaxReplicas                    int32     `json:"maxReplicas"`
	TargetCPUUtilizationPercentage *int32    `json:"targetCPUUtilizationPercentage"`
}

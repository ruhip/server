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

package cluster

import (
	"fmt"
	"net/http"
	"strings"

	k8s "server/pkg/k8s/client"
	"server/pkg/k8s/util/node"
	podUtil "server/pkg/k8s/util/pod"
	"server/pkg/utils/httpx"
	"server/pkg/utils/validate"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gorilla/mux"
)

// RegisterClusterAPI register the api of app
func RegisterClusterAPI(router *mux.Router) {
	httpx.RegisterHttpHandler(router, "/cluster", "POST", AddCluster)
	httpx.RegisterHttpHandler(router, "/clusterinfo", "GET", GetCluterInfo)
	httpx.RegisterHttpHandler(router, "/clusterinfo/nodes", "GET", GetCluserNodes)
	httpx.RegisterHttpHandler(router, "/clusterinfo/nodes/{name}", "GET", GetNodeDetail)
}

// GetCluterInfo get the cluster's info, include clusterNodes,clusterCpu,clusterMemory,clusterContainer
func GetCluterInfo(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	fake := k8s.GetClientset(clusterID)
	nodeList, err := fake.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return httpx.StatusInternalServerError, err
	}

	clusterNodes := struct {
		Total     int
		Scheduler int
		Heathy    int
	}{
		Total:     len(nodeList.Items),
		Scheduler: len(nodeList.Items),
		Heathy:    len(nodeList.Items),
	}

	clusterCpu := struct {
		Total            int64
		Allocated        int64
		AllocatedPersent int64 //需要前端自己去算百分比
	}{}

	clusterMemory := struct {
		Total            int64
		Allocated        int64
		AllocatedPersent int64 //需要前端自己去算百分比
	}{}

	clusterContainer := struct {
		Total     int64
		Operation int64
		Error     int64
	}{}

	for _, item := range nodeList.Items {
		if !node.IsNodeSchedule(&item) {
			clusterNodes.Scheduler--
		}
		if !node.IsNodeReady(&item) {
			clusterNodes.Heathy--
		}
		//cpu sum
		nodeCpuCapacity, _ := item.Status.Capacity.Cpu().AsInt64()
		clusterCpu.Total += nodeCpuCapacity
		nodeCpuAllocatable, _ := item.Status.Allocatable.Cpu().AsInt64()
		clusterCpu.Allocated = nodeCpuCapacity - nodeCpuAllocatable

		//memory sum
		nodeMemCapacity, _ := item.Status.Capacity.Memory().AsInt64()
		clusterMemory.Total += nodeMemCapacity
		nodeMemAllocatable, _ := item.Status.Allocatable.Memory().AsInt64()
		clusterMemory.Allocated = nodeMemCapacity - nodeMemAllocatable
	}

	podList, err := fake.CoreV1().
		Pods("").
		List(metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system"})
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	clusterContainer.Total = int64(len(podList.Items))
	for _, pod := range podList.Items {
		if !podUtil.IsPodReady(&pod) && (pod.Status.Phase == v1.PodFailed || pod.Status.Phase == v1.PodUnknown) {
			clusterContainer.Error++
		} else {
			clusterContainer.Operation++
		}
	}

	return httpx.StatusOK, map[string]interface{}{
		"clusterNodes":     clusterNodes,
		"clusterCpu":       clusterCpu,
		"clusterMemory":    clusterMemory,
		"clusterContainer": clusterContainer,
	}
}

//GetCluserNodes get cluster's nodes
func GetCluserNodes(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	type Node struct {
		HostName        string      `json:"hostName"`
		Internal        string      `json:"internal"`
		Status          bool        `json:"status"`
		MasterOrSlave   string      `json:"matserOrslave"`
		ContainerCnt    int         `json:"containerCnt"`
		CpuUsage        int         `json:"cpuUsage"`
		CpuAllocated    int64       `json:"cpuAllocated"`
		MemoryUsage     int         `json:"memoryUsage"`
		MemoryAllocated int64       `json:"memoryAllocated"`
		Schedulable     bool        `json:"schedulable"`
		DiskPressure    bool        `json:"diskPressure"`
		MemoryPressure  bool        `json:"memoryPressure"`
		CreateTime      metav1.Time `json:"createT_at"`
	}
	nodes := []*Node{}
	fake := k8s.GetClientset(clusterID)
	nodeList, err := fake.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	for _, item := range nodeList.Items {
		n := &Node{}
		n.HostName = node.GetHostName(&item)
		n.Internal = node.GetInternalIP(&item)
		n.Schedulable = node.IsNodeSchedule(&item)
		n.Status = node.IsNodeReady(&item)
		n.DiskPressure = node.IsDiskPressure(&item)
		n.MemoryPressure = node.IsMemoryPressure(&item)
		n.CreateTime = item.ObjectMeta.CreationTimestamp
		podList, err := fake.CoreV1().Pods("").List(metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system,spec.nodeName=" + item.Name})
		if err != nil {
			return httpx.StatusInternalServerError, err
		}
		n.ContainerCnt = len(podList.Items)

		for _, pod := range podList.Items {
			//is contain's bug
			n.MemoryAllocated += pod.Spec.Containers[0].Resources.Requests.Memory().Value()
			n.CpuAllocated += pod.Spec.Containers[0].Resources.Requests.Cpu().Value()
		}
		//TODO
		// n.CpuUsage=
		// n.MemoryUsage=

		//assert node is master or slave ？
		componetsPodList, err := fake.CoreV1().Pods("").List(metav1.ListOptions{FieldSelector: "metadata.namespace=kube-system,spec.nodeName=" + item.Name})
		if err != nil {
			return httpx.StatusInternalServerError, err
		}
		for _, pod := range componetsPodList.Items {
			//is contain's bug
			n.MemoryAllocated += pod.Spec.Containers[0].Resources.Requests.Memory().Value()
			n.CpuAllocated += pod.Spec.Containers[0].Resources.Requests.Cpu().Value()
			if strings.Contains(pod.Spec.Containers[0].Image, "kube-apiserver") ||
				strings.Contains(pod.Spec.Containers[0].Image, "kube-scheduler") ||
				strings.Contains(pod.Spec.Containers[0].Image, "kube-controller-manager") {
				n.MasterOrSlave = "master"
			} else {
				continue
			}
		}
		if n.MasterOrSlave != "master" {
			n.MasterOrSlave = "slave"
		}

		nodes = append(nodes, n)

	}

	return httpx.StatusOK, map[string]interface{}{"nodes": nodes}
}

// GetNodeDetail get node's cpu , memory and containers
//step:
//1. caculate the node's cpu , memory and containers capacity
//2. get pod of the node
func GetNodeDetail(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	nodeCpu := struct {
		Total            int64 `json:"total"`
		Allocated        int64 `json:"allocated"`
		AllocatedPersent int64 `json:"allocatedPersent"` //需要前端自己去算百分比
	}{}

	nodeMemory := struct {
		Total            int64 `json:"total"`
		Allocated        int64 `json:"allocated"`
		AllocatedPersent int64 `json:"allocatedPersent"` //需要前端自己去算百分比
	}{}

	nodeContainer := struct {
		Total     int64 `json:""`
		Allocated int64 `json:""`
	}{}

	type Container struct {
		Name      string      `json:"name"`
		Status    int         `json:"status"`
		Namespace string      `json:"namespace"`
		AppName   string      `json:"appName"`
		Image     string      `json:"image"`
		Interval  string      `json:"interval"`
		CreateAt  metav1.Time `json:"create_at"`
	}
	containres := []*Container{}
	nodeName := mux.Vars(req)["name"]
	fake := k8s.GetClientset(clusterID)
	node, err := fake.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	//cpu sum
	nodeCpuCapacity, _ := node.Status.Capacity.Cpu().AsInt64()
	nodeCpu.Total += nodeCpuCapacity
	nodeCpuAllocatable, _ := node.Status.Allocatable.Cpu().AsInt64()
	nodeCpu.Allocated = nodeCpuCapacity - nodeCpuAllocatable
	//memory sum
	nodeMemCapacity, _ := node.Status.Capacity.Memory().AsInt64()
	nodeMemory.Total += nodeMemCapacity
	nodeMemAllocatable, _ := node.Status.Allocatable.Memory().AsInt64()
	nodeMemory.Allocated = nodeMemCapacity - nodeMemAllocatable
	nodeContainerCapacity, _ := node.Status.Capacity.Pods().AsInt64()
	nodeContainer.Total = nodeContainerCapacity

	podList, err := fake.CoreV1().Pods("").List(metav1.ListOptions{FieldSelector: "metadata.namespace!=kube-system,spec.nodeName=" + nodeName})
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	nodeContainer.Allocated = int64(len(podList.Items))

	for _, pod := range podList.Items {
		container := &Container{}
		container.Name = pod.Name
		container.Namespace = pod.Namespace
		container.Image = pod.Spec.Containers[0].Image
		container.Interval = fmt.Sprintf("%s:%v", pod.Status.PodIP, pod.Spec.Containers[0].Ports[0].ContainerPort)
		container.CreateAt = pod.ObjectMeta.CreationTimestamp
		containres = append(containres, container)
	}
	return httpx.StatusOK, map[string]interface{}{
		"nodeCpu":       nodeCpu,
		"nodeMemory":    nodeMemory,
		"nodeContainer": nodeContainer,
		"containres":    containres,
	}
}

// AddCluster add cluster
func AddCluster(req *http.Request) (string, interface{}) {
	cluster, err := validate.ValidateCluster(req)
	result, err := cluster.Insert()
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	if err = k8s.AddClientset(result); err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, "add cluster success"
}

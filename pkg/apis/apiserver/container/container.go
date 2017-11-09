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

package container

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	api "server/pkg/api/apiserver"
	"server/pkg/apis/apiserver"
	"server/pkg/configz"
	"server/pkg/utils/httpx"
	"server/pkg/utils/log"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gorilla/mux"
)

// maximum number of lines loaded from the apiserver
var lineReadLimit int64 = 5000

// maximum number of bytes loaded from the apiserver
var byteReadLimit int64 = 500000

//RegisterContainerAPI register the api of pod
func RegisterContainerAPI(router *mux.Router) {
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/containers/{name}/events", "GET", GetContainerEvents)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/containers/{name}/logs", "GET", GetContainerLog)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/containers/{name}/process", "GET", GetContainerProcess)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/containers", "GET", GetContainers)
}

//GetContainerEvents get container's events
func GetContainerEvents(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	name := params["name"]
	events, err := apiserver.GetPodEvents(name, namespace, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, events
}

//GetContainerLog get container's log
func GetContainerLog(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	name := params["name"]
	sinceTimeSTR := req.FormValue("sinceTime")

	pod, err := apiserver.GetPod(name, namespace, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}

	nowTime := time.Now()
	sinceTime := metav1.NewTime(nowTime)
	switch sinceTimeSTR {
	case "1h":
		sinceTime = metav1.NewTime(nowTime.Add(-time.Hour * 1))
	case "6h":
		sinceTime = metav1.NewTime(nowTime.Add(-time.Hour * 6))
	case "1d":
		sinceTime = metav1.NewTime(nowTime.AddDate(0, 0, -1))
	case "1w":
		sinceTime = metav1.NewTime(nowTime.AddDate(0, 0, -7))
	case "1m":
		sinceTime = metav1.NewTime(nowTime.AddDate(0, -1, 0))
	default:
		sinceTime = metav1.NewTime(nowTime.Add(-time.Hour * 1))
	}

	containerName := pod.Spec.Containers[0].Name
	logOptions := &v1.PodLogOptions{
		Container:  containerName,
		Follow:     false,
		Previous:   false,
		SinceTime:  &sinceTime,
		Timestamps: true,
		LimitBytes: &byteReadLimit,
		TailLines:  &lineReadLimit,
	}

	result, err := apiserver.GetPodLogs(name, namespace, clusterID, logOptions)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, map[string]interface{}{"logs": result}
}

//GetContainerProcess get container's process
func GetContainerProcess(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	name := params["name"]
	pod, err := apiserver.GetPod(name, namespace, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	containerID := strings.Replace(pod.Status.ContainerStatuses[0].ContainerID, "docker://", "", -1)
	tranport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}
	url := configz.GetString("apiserver", "cadvisor", "http://127.0.0.1:4194")
	client := &http.Client{Transport: tranport}
	url = fmt.Sprintf(url, pod.Status.HostIP)
	res, err := client.Get(url + `/api/v1.0/containers/kubepods.slice`)
	if err != nil {
		return httpx.StatusInternalServerError, fmt.Sprintf("get process of container [%s] err:%v", pod.Name, err.Error())
	}

	var (
		cadvisorResult = api.CadvisorResult{}
		processUrl     = ""
	)

	if err = json.NewDecoder(res.Body).Decode(&cadvisorResult); err != nil {
		return httpx.StatusInternalServerError, err
	}

	for k := range cadvisorResult.Subcontainers {
		cadvisorResult1 := api.CadvisorResult{}
		subContainer := cadvisorResult.Subcontainers[k].Name
		res, err := client.Get(url + "/api/v1.0/containers" + subContainer)
		if err = json.NewDecoder(res.Body).Decode(&cadvisorResult1); err != nil {
			return httpx.StatusInternalServerError, err
		}
		for k := range cadvisorResult1.Subcontainers {
			subContainer := cadvisorResult1.Subcontainers[k].Name
			if strings.Contains(subContainer, containerID) {
				processUrl = "/api/v2.0/ps" + subContainer
				goto PROCESS
			}
		}
	}

PROCESS:
	result, err := client.Get(url + processUrl)
	if err != nil {
		log.Debug(err)
		return httpx.StatusInternalServerError, fmt.Sprintf("get process of container [%s] err:%v", pod.Name, err.Error())
	}
	processes := []*api.Process{}
	if err = json.NewDecoder(result.Body).Decode(&processes); err != nil {
		return httpx.StatusInternalServerError, err
	}

	return httpx.StatusOK, map[string]interface{}{"processes": processes}
}

//GetContainers get containers
func GetContainers(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	pods, err := apiserver.GetPods(namespace, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, pods
}

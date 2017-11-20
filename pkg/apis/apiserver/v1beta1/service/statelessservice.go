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

package service

import (
	"fmt"
	"net/http"

	"server/pkg/apis/apiserver/v1beta1"
	"server/pkg/utils/httpx"
	"server/pkg/utils/log"
	"server/pkg/utils/parseUtil"
	"server/pkg/utils/validate"

	"github.com/gorilla/mux"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

//RegisterServiceAPI register the api of service
func RegisterStatelessServiceAPI(router *mux.Router) {
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/services", "GET", ListService)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/services", "POST", DeployService)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/services/{serviceName}", "DELETE", DeleteService)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/services/{serviceName}/{verb}", "PATCH", RollOrExpansionOrScaleService)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/services/{serviceName}/{verb}", "PUT", StopOrStartOrRedeployService)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/services/{serviceName}/events", "GET", GetServiceEvents)
}

//DeployService deploy service
//router /api/v1/namespaces/{namespace}/services
//method POST
func DeployService(req *http.Request) (string, interface{}) {
	svc, err := validate.ValidateService(req)
	if err != nil {
		return httpx.StatusBadRequest, err
	}
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	service := svc.TOK8SService(namespace)
	deployment := svc.TOK8SDeployment(namespace)
	result, err := apiserver.DelpoyService(clusterID, service, deployment)
	if err != nil {
		log.Error("deploy app where named %q err: %v", svc.Name, err)
		return httpx.StatusInternalServerError, fmt.Errorf("deploy app where named %q err: ", err)
	}
	return httpx.StatusOK, result
}

//DeleteService deploy service
//router /api/v1/namespaces/{namespace}/services/{serviceName}
//method DELETE
func DeleteService(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	name := params["serviceName"]
	namespace := params["namespace"]
	clusterID := params["clusterID"]
	if err := apiserver.DeleteService(name, namespace, clusterID); err != nil {
		log.Error("delete service where named %q err: %v", name, err)
		return httpx.StatusInternalServerError, fmt.Errorf("delete service where named %q err: %v", name, err)
	}
	return httpx.StatusOK, "success"
}

//StopOrStartOrRedeployService start or stop or redeploy service
func StopOrStartOrRedeployService(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	verb := params["verb"]
	name := params["serviceName"]
	namespace := params["namespace"]
	clusterID := params["clusterID"]
	svc, deploy, exist := apiserver.ServiceExist(name, namespace, clusterID)
	if !exist {
		return httpx.StatusNotFound, fmt.Errorf("service %v not found", name)
	}
	if err := apiserver.StartOrStopOrRedployService(svc, deploy, verb, clusterID); err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, fmt.Sprintf("%v service %v success", verb, name)
}

//ListService list service
func ListService(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	services, deployments, err := apiserver.ListService(namespace, clusterID)
	if err != nil {
		log.Error("get service who's namespace is %q err: %v", namespace, err)
		return httpx.StatusInternalServerError, err
	}
	result := map[string]interface{}{
		"services":   services,
		"containers": deployments,
	}
	return httpx.StatusOK, result
}

//GetServiceEvents get service events
func GetServiceEvents(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	name := params["serviceName"]
	events, err := apiserver.GetServiceEvents(name, namespace, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, events
}

//RollOrExpansionOrScaleService roll or expansion or scale service
func RollOrExpansionOrScaleService(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	verb := params["verb"]
	if verb == "roll" {
		return rollUpdateService(req)
	}
	if verb == "expansion" {
		return expansionService(req)
	}
	if verb == "scale" {
		return scaleService(req)
	}
	return httpx.StatusNotFound, fmt.Sprintf("%v not support current operation", verb)
}

func rollUpdateService(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	name := params["serviceName"]
	image := req.FormValue("image")
	deploy, err := apiserver.GetDeployment(name, namespace, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	deploy.Spec.Template.Spec.Containers[0].Image = image
	result, err := apiserver.UpdateDeployment(deploy, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, result
}

func expansionService(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	name := params["serviceName"]
	cpu := params["cpu"]
	memory := params["memory"]
	deploy, err := apiserver.GetDeployment(name, namespace, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	deploy.Spec.Template.Spec.Containers[0].Resources = v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse(cpu),    //TODO 根据前端传入的值做资源限制
			v1.ResourceMemory: resource.MustParse(memory), //TODO 根据前端传入的值做资源限制
		},
		Requests: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse(cpu),
			v1.ResourceMemory: resource.MustParse(memory),
		},
	}
	result, err := apiserver.UpdateDeployment(deploy, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, result
}

func scaleService(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	name := params["serviceName"]
	replicas := params["replicas"]
	deploy, err := apiserver.GetDeployment(name, namespace, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	deploy.Spec.Replicas = parseUtil.StringToInt32Pointer(replicas)
	result, err := apiserver.UpdateDeployment(deploy, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, result
}

func AddServicePorts(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	name := params["name"]
	servicePorts, err := validate.ValidatePorts(req)
	if err != nil {
		return httpx.StatusBadRequest, err
	}
	containerPorts := []v1.ContainerPort{}
	for _, p := range servicePorts {
		containerPorts = append(containerPorts, v1.ContainerPort{ContainerPort: int32(p.TargetPort.IntVal)})
	}
	deploy, err := apiserver.GetDeployment(name, namespace, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	if deploy == nil {
		return httpx.StatusNotFound, fmt.Sprintf("serivce %v not found", name)
	}
	deploy.Spec.Template.Spec.Containers[0].Ports = containerPorts

	svc, err := apiserver.GetK8SService(name, namespace, clusterID)
	svc.Spec.Ports = servicePorts
	resultsvc, resultdp, err := apiserver.UpdateService(svc, deploy, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, map[string]interface{}{"service": resultsvc, "deployment": resultdp}
}

func AddServiceEnvs(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	name := params["name"]
	servicePorts, err := validate.ValidatePorts(req)
	if err != nil {
		return httpx.StatusBadRequest, err
	}
	containerPorts := []v1.ContainerPort{}
	for _, p := range servicePorts {
		containerPorts = append(containerPorts, v1.ContainerPort{ContainerPort: int32(p.TargetPort.IntVal)})
	}
	deploy, err := apiserver.GetDeployment(name, namespace, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	if deploy == nil {
		return httpx.StatusNotFound, fmt.Sprintf("serivce %v not found", name)
	}
	deploy.Spec.Template.Spec.Containers[0].Ports = containerPorts
	result, err := apiserver.UpdateDeployment(deploy, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, result
}

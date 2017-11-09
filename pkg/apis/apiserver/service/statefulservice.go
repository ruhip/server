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

	"server/pkg/apis/apiserver"
	"server/pkg/utils/httpx"
	"server/pkg/utils/log"
	"server/pkg/utils/validate"

	"github.com/gorilla/mux"
)

//RegisterServiceAPI register the api of service
func RegisterStatefulServiceAPI(router *mux.Router) {
	// httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/statefulservices", "GET", ListService)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/statefulservices", "POST", DeployStatefulService)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/statefulservices/{serviceName}", "DELETE", DeleteStatefulService)
	// httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/statefulservices/{serviceName}/{verb}", "PATCH", RollOrExpansionOrScaleService)
	// httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/statefulservices/{serviceName}/{verb}", "PUT", StopOrStartOrRedeployService)
	// httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/statefulservices/{serviceName}/events", "GET", GetServiceEvents)
}

// DeployStatefulService deploy stateful service
func DeployStatefulService(req *http.Request) (string, interface{}) {
	svc, err := validate.ValidateService(req)
	if err != nil {
		return httpx.StatusBadRequest, err
	}
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]

	service := svc.TOK8SService(namespace)
	headlessService := svc.TOK8SHeadlessService(namespace)
	pvc := svc.TOPersistentVolumeClaim(namespace) // don't create pvc by manue,because the statefulset will auto create it
	statefulset := svc.TOK8SStatefulset(namespace, headlessService.Name, *pvc)
	result, err := apiserver.DeploySatefulService(service, headlessService, statefulset, clusterID)
	if err != nil {
		log.Error("deploy statefulset service  %q err: %v", svc.Name, err)
		return httpx.StatusInternalServerError, fmt.Errorf("deploy statefulset service  %q err: %v", svc.Name, err)
	}
	return httpx.StatusOK, result
}

// DeleteStatefulService delete stateful service
func DeleteStatefulService(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	serviceName := params["serviceName"]
	if err := apiserver.DeleteStatefulService(serviceName, namespace, clusterID); err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, fmt.Sprintf("delete stateful service %v success", serviceName)
}

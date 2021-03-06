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

package app

import (
	"fmt"
	"net/http"

	models "server/pkg/api/apiserver/v1beta1"
	"server/pkg/apis/apiserver/v1beta1"
	"server/pkg/utils/httpx"
	"server/pkg/utils/log"
	"server/pkg/utils/validate"

	"github.com/gorilla/mux"
)

// RegisterAppAPI register the api of app
func RegisterAppAPI(router *mux.Router) {
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/apps", "GET", ListApp)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/apps", "POST", DeployApp)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/apps/{appName}", "DELETE", DeleteApp)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/apps/{appName}/{verb}", "PUT", StopOrStartOrRedeployApp)
}

// DeployApp deploy app
// router /api/v1/clusters/{clusterID}/namespaces/{namespace}/apps
// method POST
func DeployApp(req *http.Request) (string, interface{}) {
	app, err := validate.ValidateApp(req)
	if err != nil {
		return httpx.StatusBadRequest, err
	}
	svc := app.Items[0]
	svc.AppName = app.Name
	clusterID := mux.Vars(req)["clusterID"]
	namespace := mux.Vars(req)["namespace"]
	service := svc.TOK8SService(namespace)
	deployment := svc.TOK8SDeployment(namespace)
	result, err := apiserver.DelpoyService(clusterID, service, deployment)
	if err != nil {
		log.Error("deploy app where named %q err: ", err)
		return httpx.StatusInternalServerError, fmt.Errorf("deploy app where named %q err: ", err)
	}

	app.ServiceCount = svc.InstanceCount
	app.InstanceCount = int(*(deployment.Spec.Replicas))
	if err = app.Insert(); err != nil {
		log.Error("record to db err: %v", err)
	}
	return httpx.StatusOK, result
}

// DeleteApp delete app
// router /api/v1/clusters/{clusterID}/namespaces/{namespace}/apps/{name}
// method DELETE
func DeleteApp(req *http.Request) (string, interface{}) {
	name := mux.Vars(req)["appName"]
	namespace := mux.Vars(req)["namespace"]
	clusterID := mux.Vars(req)["clusterID"]
	if err := apiserver.DeleteServiceByAppName(name, namespace, clusterID); err != nil {
		return httpx.StatusInternalServerError, err
	}
	app := &models.App{}
	app.UserName = namespace
	app.Name = name
	if err := app.DeleteByNameAndNamespace(); err != nil {
		log.Error("delete %v's app %v from db err: %v", namespace, name, err)
	}
	return httpx.StatusOK, fmt.Sprintf("delete namespace %v's app named %v success", namespace, name)
}

//GetApps get apps
func GetApps(request *http.Request) (string, interface{}) {
	return "", ""
}

// GetApp get app
func GetApp(request *http.Request) (string, interface{}) {
	return "", ""
}

// StopOrStartOrRedeployApp start stop or redploy app
func StopOrStartOrRedeployApp(req *http.Request) (string, interface{}) {
	verb := mux.Vars(req)["verb"]
	name := mux.Vars(req)["appName"]
	namespace := mux.Vars(req)["namespace"]
	clusterID := mux.Vars(req)["clusterID"]
	svcs, deploys, err := apiserver.ListServiceByAppName(name, namespace, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	errs := apiserver.StartOrStopOrRedployApp(svcs, deploys, verb, clusterID)
	if len(errs) != 0 {
		return httpx.StatusInternalServerError, errs
	}
	return httpx.StatusOK, fmt.Sprintf("%v app %v success", verb, name)
}

// ListApp list app by namespace
func ListApp(req *http.Request) (string, interface{}) {
	namespace := mux.Vars(req)["namespace"]
	app := &models.App{}
	app.UserName = namespace
	result, err := app.GetByNamespace()
	if err != nil {
		log.Error("get app where namespace %q err: %v", namespace, err)
		return httpx.StatusInternalServerError, fmt.Errorf("get app where namespace %q err: %v", namespace, err)
	}
	return httpx.StatusOK, map[string]interface{}{"apps": result}
}

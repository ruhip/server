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

package configmap

import (
	"fmt"
	"net/http"

	"server/pkg/apis/apiserver"
	"server/pkg/utils/httpx"
	"server/pkg/utils/validate"

	"github.com/gorilla/mux"
)

//RegisterConfigAPI register the api of ConfigMap
func RegisterConfigAPI(router *mux.Router) {
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/configs", "GET", ListConfig)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/configs", "POST", CreateConfig)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/configs/{name}", "POST", AddConfigData)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/configs/{name}", "DELETE", DeleteConfig)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/configs/{name}", "PUT", UpdateConfig)
}

//CreateConfig crete configMap
func CreateConfig(req *http.Request) (string, interface{}) {
	config, err := validate.ValidateConfig(req)
	if err != nil {
		return httpx.StatusBadRequest, err
	}
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	configMap := config.TOK8SConfigMap(namespace)
	result, err := apiserver.CreateConfigMap(clusterID, configMap)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, result
}

//AddConfigData add config data
func AddConfigData(req *http.Request) (string, interface{}) {
	data, err := validate.ValidateConfigData(req)
	if err != nil {
		return httpx.StatusBadRequest, err
	}
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	name := params["name"]
	configMap, err := apiserver.GetConfigMapByName(name, namespace, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	for k, v := range data {
		if configMap.Data != nil {
			configMap.Data[k] = v
		} else {
			configMap.Data = data
		}
	}
	result, err := apiserver.UpdateConfigMap(clusterID, configMap)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, result
}

//UpdateConfig update configMap
func UpdateConfig(req *http.Request) (string, interface{}) {
	config, err := validate.ValidateConfig(req)
	if err != nil {
		return httpx.StatusBadRequest, err
	}
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	name := params["name"]
	configMap, err := apiserver.GetConfigMapByName(name, namespace, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	for k, v := range config.Data {
		if configMap.Data != nil {
			configMap.Data[k] = v
		} else {
			configMap.Data = config.Data
		}
	}
	result, err := apiserver.UpdateConfigMap(clusterID, configMap)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, result
}

//DeleteConfig delete configMap
func DeleteConfig(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	name := params["name"]
	err := apiserver.DeleteConfigMap(name, namespace, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, fmt.Sprintf("delete configMap %v success", name)
}

//ListConfig get all configMap
func ListConfig(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	result, err := apiserver.ListConfigMap(namespace, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, result
}

// Copyright © 2018 huang jia <449264675@qq.com>
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

package storage

import (
	"fmt"
	"net/http"

	models "server/pkg/api/apiserver/v1beta1"
	apiserver "server/pkg/apis/apiserver/v1beta1"
	"server/pkg/utils/httpx"
	"server/pkg/utils/log"
	"server/pkg/utils/validate"

	"github.com/gorilla/mux"
)

//RegisterPersistentVolumeClaimAPI register the api of PersistentVolumeClaim
func RegisterPersistentVolumeClaimAPI(router *mux.Router) {
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/persistentvolumeclaims", "GET", ListPersistentVolumeClaim)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/persistentvolumeclaims", "POST", CreatePersistentVolumeClaim)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/persistentvolumeclaims/{name}", "DELETE", DeletePersistentVolumeClaim)
}

// CreatePersistentVolumeClaim create PersistentVolumeClaim
func CreatePersistentVolumeClaim(req *http.Request) (string, interface{}) {
	// step 1: validate request data
	storage, err := validate.ValidateStorage(req)
	if err != nil {
		return httpx.StatusBadRequest, err
	}
	// step 2: create k8s persistentvoumeclaim resource
	clusterID := mux.Vars(req)["clusterID"]
	namespace := mux.Vars(req)["namespace"]
	pvc := storage.TOK8SPersistentVolumeClaim(namespace)
	result, err := apiserver.CreatePersistentVolumeClaim(pvc, clusterID)
	if err != nil {
		log.Error("create persistentvolumeclaims named %q err: %v", storage.Name, err)
		return httpx.StatusInternalServerError, fmt.Errorf("create persistentvolumeclaims named %q err: %v", storage.Name, err)
	}
	// step 3: record pvc to db
	storage.Namespace = namespace
	if err = storage.Insert(); err != nil {
		log.Error("insert storage record to db err: %v", err)
	}
	return httpx.StatusOK, result
}

// DeletePersistentVolumeClaim delete PersistentVolumeClaim
func DeletePersistentVolumeClaim(req *http.Request) (string, interface{}) {
	// step 1: delete k8s persistentvoumeclaim resource
	clusterID := mux.Vars(req)["clusterID"]
	name := mux.Vars(req)["name"]
	namespace := mux.Vars(req)["namespace"]
	if err := apiserver.DeletePersistentVolumeClaim(name, namespace, clusterID); err != nil {
		log.Error("delete persistentvolumeclaims named %q err: %v", name, err)
		return httpx.StatusInternalServerError, fmt.Errorf("delete persistentvolumeclaims named %q err: %v", name, err)
	}
	// step 2: delete db record
	storage := &models.Storage{Name: name}
	if err := storage.DeleteByName(); err != nil {
		log.Error("delete storage record %v from db err: %v", name, err)
	}
	return httpx.StatusOK, fmt.Sprintf("delete persistentvolumeclaims %v success", name)
}

// ListPersistentVolumeClaim list PersistentVolumeClaim by namespace
func ListPersistentVolumeClaim(req *http.Request) (string, interface{}) {
	namespace := mux.Vars(req)["namespace"]
	storage := &models.Storage{Name: namespace}
	storage.Namespace = namespace
	result, err := storage.GetByNamespace()
	if err != nil {
		log.Error("list namespace %q's storage  err: %v", namespace, err)
		return httpx.StatusInternalServerError, fmt.Errorf("list namespace %q's storage  err: %v", namespace, err)
	}
	return httpx.StatusOK, map[string]interface{}{"storages": result}
}

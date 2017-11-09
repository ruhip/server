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

package storageclass

import (
	"fmt"
	"net/http"

	"server/pkg/apis/apiserver"
	"server/pkg/utils/httpx"
	"server/pkg/utils/log"
	"server/pkg/utils/validate"

	"github.com/gorilla/mux"
)

//RegisterStorageClassAPI register the api of storageclass
func RegisterStorageClassAPI(router *mux.Router) {
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/storageclasses", "POST", CreateStorageClass)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/storageclasses/{name}", "DELETE", DeleteStorageClass)
}

//CreateStorageClass create storageclass
func CreateStorageClass(req *http.Request) (string, interface{}) {
	cephrbd, err := validate.ValidateCephRBD(req)
	if err != nil {
		return httpx.StatusBadRequest, err
	}
	clusterID := mux.Vars(req)["clusterID"]
	storageclass := cephrbd.TOStorageClass()
	result, err := apiserver.CreateStorageClass(storageclass, clusterID)
	if err != nil {
		log.Error("create storageclass named %q err: %v", cephrbd.Name, err)
		return httpx.StatusInternalServerError, fmt.Errorf("create storageclass named %q err: %v", cephrbd.Name, err)
	}
	return httpx.StatusOK, result
}

//DeleteStorageClass delete storage class
func DeleteStorageClass(req *http.Request) (string, interface{}) {
	clusterID := mux.Vars(req)["clusterID"]
	name := mux.Vars(req)["name"]
	if err := apiserver.DeleteStorageClass(name, clusterID); err != nil {
		log.Error("delete storageclass named %q err: %v", cephrbd.Name, err)
		return httpx.StatusInternalServerError, fmt.Errorf("delete storageclass named %q err: %v", cephrbd.Name, err)
	}
	return httpx.StatusOK, fmt.Sprintf("delete storageclass %v success", name)
}

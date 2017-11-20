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

package metric

import (
	"net/http"

	"server/pkg/k8s/metric"
	"server/pkg/utils/httpx"

	"github.com/gorilla/mux"
)

//RegisterMetricAPI register the api of metric
func RegisterMetricAPI(router *mux.Router) {
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/containers/{name}/{metric}/{type}", "GET", GetMetrics)
}

//GetMetrics return container's metric
func GetMetrics(request *http.Request) (string, interface{}) {
	namespace := mux.Vars(request)["namespace"]
	podName := mux.Vars(request)["name"]
	metricsName := mux.Vars(request)["metric"] + "/" + mux.Vars(request)["type"]
	metrics, err := metric.GetPodMetrics(namespace, podName, metricsName)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, map[string]interface{}{"metrics": metrics}
}

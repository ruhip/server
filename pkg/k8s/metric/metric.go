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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"server/pkg/utils/log"

	"server/pkg/configz"
)

/*
**metricName参考：**
	[
	"network/tx",
	"network/tx_errors_rate",
	"memory/working_set",
	"network/tx_errors",
	"cpu/limit",
	"memory/major_page_faults",
	"memory/page_faults_rate",
	"cpu/request",
	"network/rx_rate",
	"cpu/usage_rate",
	"memory/limit",
	"memory/usage",
	"memory/cache",
	"network/rx_errors",
	"network/rx_errors_rate",
	"network/tx_rate",
	"memory/major_page_faults_rate",
	"cpu/usage",
	"network/rx",
	"memory/rss",
	"memory/page_faults",
	"memory/request",
	"uptime"
	]
*/

//GetPodMetrics get pod metric
func GetPodMetrics(namespace, podName, metric_name string) (map[string]interface{}, error) {
	path := fmt.Sprintf("%s/api/v1/model/namespaces/%s/pods/%s/metrics/%s", configz.GetString("apiserver", "heapsterEndpoint", "127.0.0.1:30003"), namespace, podName, metric_name)
	log.Info(path)
	heapsterHost := configz.GetString("apiserver", "heapsterEndpoint", "http://127.0.0.1:30003")
	log.Info("Creating remote Heapster client for %s", heapsterHost)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	v := map[string]interface{}{}
	json.Unmarshal(data, &v)
	return v, nil
}

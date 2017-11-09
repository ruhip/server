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

	models "server/pkg/api/apiserver"
	"server/pkg/apis/apiserver/app"
	"server/pkg/apis/apiserver/configmap"
	"server/pkg/apis/apiserver/container"
	"server/pkg/apis/apiserver/hpa"
	"server/pkg/apis/apiserver/metric"
	"server/pkg/apis/apiserver/service"
	"server/pkg/componentconfig"
	"server/pkg/configz"
	"server/pkg/k8s/client"
	"server/pkg/storage/mysql"
	"server/pkg/utils/log"

	"github.com/gorilla/mux"
)

//APIServer apiserver component config
type APIServer struct {
	*componentconfig.ApiserverConfig
}

//NewAPIServer return Apiserver
func NewAPIServer(config string) *APIServer {
	configz.Init(config)
	mysql.Init()
	models.Init()
	client.Init()
	hpa.Init()
	return &APIServer{
		ApiserverConfig: &componentconfig.ApiserverConfig{
			HTTPAddr: configz.GetString("apiserver", "httpAddr", "0.0.0.0"),
			HTTPPort: configz.MustInt("apiserver", "httpPort", 9090),
			RPCAddr:  configz.GetString("apiserver", "rpcAddr", "0.0.0.0"),
			RPCPort:  configz.MustInt("apiserver", "rpcPort", 7070),
		},
	}
}

//Run start apiserver component
func Run(server *APIServer) error {
	root := mux.NewRouter()
	api := root.PathPrefix("/api/v1/clusters/{clusterID}").Subrouter()
	installAPIGroup(api)
	http.Handle("/", root)
	log.Info("starting apiserver and listen on : %v", fmt.Sprintf("%v:%v", server.HTTPAddr, server.HTTPPort))
	go configz.Heatload()

	return http.ListenAndServe(fmt.Sprintf("%v:%v", server.HTTPAddr, server.HTTPPort), nil)
}

func installAPIGroup(router *mux.Router) {
	app.RegisterAppAPI(router)
	service.RegisterStatelessServiceAPI(router)
	service.RegisterStatefulServiceAPI(router)
	configmap.RegisterConfigAPI(router)
	container.RegisterContainerAPI(router)
	hpa.RegisterHPAAPI(router)
	metric.RegisterMetricAPI(router)
}

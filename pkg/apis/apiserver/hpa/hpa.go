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

package hpa

import (
	"fmt"
	"net/http"
	"strconv"

	models "server/pkg/api/apiserver"
	"server/pkg/apis/apiserver"
	"server/pkg/utils/httpx"
	"server/pkg/utils/log"
	"server/pkg/utils/validate"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"gopkg.in/robfig/cron.v2"
)

//RegisterAppAPI register the api of app
func RegisterHPAAPI(router *mux.Router) {
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/hpas", "POST", CreateHPA)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/hpas/{name}", "DELETE", DeleteHPA)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/hpas/{name}", "PUT", UpdateHPA)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/hpas/tickScaleTask/{id}", "DELETE", DeleteTickScaleTask)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/hpas/tickScaleTask/{id}/off", "PUT", TickScaleTaskOff)
	httpx.RegisterHttpHandler(router, "/namespaces/{namespace}/hpas/tickScaleTask/{id}/on", "PATCH", TickScaleTaskOn)
}

var (
	crontab = cron.New()
)

func Init() {
	crontab.Start()
	notifyTickScaleTask()
}

//notifyTickScaleTask when the process is start or restart , notify the task of service
//1. query the db to find the task of  all service
//2. add all task to the tast queue
func notifyTickScaleTask() {
	tasks, err := new(models.TickScaleTask).GetAll()
	if err != nil {
		glog.Errorf("loading service's task err: %v", err)
		return
	}

	for _, task := range tasks {
		if task.Status == 1 {
			addTickScaleTask(task)
		}
	}
}

func addTickScaleTask(task *models.TickScaleTask) error {
	hpa := &models.HPA{
		Name:                           task.Name,
		RefObjectName:                  task.RefObjectName,
		MaxReplicas:                    task.MaxReplicas,
		MinReplicas:                    task.MinReplicas,
		TargetCPUUtilizationPercentage: task.TargetCPUUtilizationPercentage,
	}
	k8sHPA := hpa.TOK8SHPA(task.Namespace)
	if _, err := crontab.AddFunc(task.Spec, func() {
		apiserver.CreateHPA(task.ClusterID, k8sHPA)
	}, int(task.ID)); err != nil {
		log.Error("add task  %v err: %v ", task.Name, err)
		return err
	}
	glog.Infoln("add task  %v successed", task.Name)
	return nil
}
func verbTickScaleTask(id int, verb string) error {
	tickScaleTask := new(models.TickScaleTask)
	tickScaleTask.ID = uint(id)
	task, err := tickScaleTask.Get()
	if err != nil {
		log.Error("query task from db err: %v ", err)
		return err
	}
	if err = apiserver.DeleteHPA(task.Name, task.Namespace, task.ClusterID); err != nil {
		log.Error("delete hpa %v's task err: %v ", task.Name, err)
		return err
	}
	if verb == "delete" {
		if err := tickScaleTask.Delete(); err != nil {
			log.Error("delete tickScaleTask err: %v", err)
			return err
		}
	} else {
		if verb == "on" {
			task.Status = 1
		}
		if verb == "off" {
			task.Status = 0
		}
		if err = task.Update(); err != nil {
			log.Error("update  tickScaleTask status to %v err: %v", task.Status, err)
			return err
		}
	}
	crontab.Remove(cron.EntryID(id))
	return nil
}

//CreateHPA create k8s HPA
func CreateHPA(req *http.Request) (string, interface{}) {
	hpa, err := validate.ValidateHPA(req)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	k8sHPA := hpa.TOK8SHPA(namespace)
	result, err := apiserver.CreateHPA(clusterID, k8sHPA)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, result
}

//DeleteHPA delete k8s HPA
func DeleteHPA(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	name := params["name"]
	if err := apiserver.DeleteHPA(name, namespace, clusterID); err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, fmt.Sprintf("delete hpa %v success", name)
}

//UpdateHPA update k8s HPA
func UpdateHPA(req *http.Request) (string, interface{}) {
	hpa, err := validate.ValidateHPA(req)
	if err != nil {
		return httpx.StatusBadRequest, err
	}
	params := mux.Vars(req)
	clusterID := params["clusterID"]
	namespace := params["namespace"]
	name := params["name"]
	k8sHPA, err := apiserver.GetHPA(name, namespace, clusterID)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	k8sHPA.Spec.MaxReplicas = hpa.MaxReplicas
	k8sHPA.Spec.MinReplicas = hpa.MinReplicas
	k8sHPA.Spec.TargetCPUUtilizationPercentage = hpa.TargetCPUUtilizationPercentage
	result, err := apiserver.UpdateHPA(clusterID, k8sHPA)
	if err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, result
}

// CreateTickScaleTask create TickScaleTask
func CreateTickScaleTask(req *http.Request) (string, interface{}) {
	tickScaleTask, err := validate.ValidateTickScaleTask(req)
	if err != nil {
		return httpx.StatusBadRequest, err
	}
	err = tickScaleTask.Insert()
	if err != nil {
		log.Error("insert tickScaleTask err: %v", err)
		return httpx.StatusInternalServerError, err
	}
	if err = addTickScaleTask(tickScaleTask); err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, tickScaleTask
}

// DeleteTickScaleTask delete TickScaleTask
func DeleteTickScaleTask(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	idstr := params["id"]
	id, _ := strconv.Atoi(idstr)
	if err := verbTickScaleTask(id, "delete"); err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, fmt.Sprintf("delete tickScaleTask %v success", id)
}

// TickScaleTaskOn TickScaleTask on
func TickScaleTaskOn(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	idstr := params["id"]
	id, _ := strconv.Atoi(idstr)
	if err := verbTickScaleTask(id, "on"); err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, nil
}

// TickScaleTaskOff TickScaleTask off
func TickScaleTaskOff(req *http.Request) (string, interface{}) {
	params := mux.Vars(req)
	idstr := params["id"]
	id, _ := strconv.Atoi(idstr)
	if err := verbTickScaleTask(id, "off"); err != nil {
		return httpx.StatusInternalServerError, err
	}
	return httpx.StatusOK, nil
}

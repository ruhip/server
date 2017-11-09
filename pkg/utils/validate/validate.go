package validate

import (
	"encoding/json"
	"net/http"

	"server/pkg/api/apiserver"

	"k8s.io/api/core/v1"
)

func ValidateApp(req *http.Request) (*apiserver.App, error) {
	app := &apiserver.App{}
	if err := json.NewDecoder(req.Body).Decode(app); err != nil {
		return nil, err
	}
	return app, nil
}

func ValidateService(req *http.Request) (*apiserver.Service, error) {
	svc := &apiserver.Service{}
	if err := json.NewDecoder(req.Body).Decode(svc); err != nil {
		return nil, err
	}
	return svc, nil
}

func ValidateConfig(req *http.Request) (*apiserver.Config, error) {
	config := &apiserver.Config{}
	if err := json.NewDecoder(req.Body).Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

func ValidateConfigData(req *http.Request) (map[string]string, error) {
	data := map[string]string{}
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func ValidateHPA(req *http.Request) (*apiserver.HPA, error) {
	hpa := &apiserver.HPA{}
	if err := json.NewDecoder(req.Body).Decode(hpa); err != nil {
		return nil, err
	}
	return hpa, nil
}

func ValidatePorts(req *http.Request) ([]v1.ServicePort, error) {
	ports := []v1.ServicePort{}
	err := json.NewDecoder(req.Body).Decode(&ports)
	return ports, err
}

func ValidateEnvs(req *http.Request) ([]v1.EnvVar, error) {
	envs := []v1.EnvVar{}
	err := json.NewDecoder(req.Body).Decode(&envs)
	return envs, err
}

func ValidateCephRBD(req *http.Request) (*apiserver.CephRBD, error) {
	rbd := &apiserver.CephRBD{}
	err := json.NewDecoder(req.Body).Decode(rbd)
	return rbd, err
}

func ValidateTickScaleTask(req *http.Request) (*apiserver.TickScaleTask, error) {
	task := &apiserver.TickScaleTask{}
	err := json.NewDecoder(req.Body).Decode(task)
	return task, err
}

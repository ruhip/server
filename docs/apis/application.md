**创建应用**

- URL: 127.0.0.1:9090/api/v1/clusters/{clusterID}/namespaces/{namespace}/apps
 
- METHOD: POST

**请求：**
- 127.0.0.1:9090/api/v1/clusters/c2a523625ba526a7/namespaces/huangjia/apps

```
{
    "nmae": "nginx",
    "nameSpace": "huangjia",
    "description": "test app",
    "services": [
        {
            "name": "demo",
            "image": "nginx:latest",
            "instanceCount": 1,
            "type": 0,
            "cpu": "256m",
            "memory": "256Mi",
            "ports": [
                {
                    "protocol": "TCP",
                    "port": 80,
                    "targetPort": 80
                }
            ]
        }
    ]
}
```

**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": [
        {
            "metadata": {
                "name": "demo",
                "namespace": "huangjia",
                "selfLink": "/apis/apps/v1beta1/namespaces/huangjia/deployments/demo",
                "uid": "1ed02a32-d98b-11e7-babc-525410406eaa",
                "resourceVersion": "8627961",
                "generation": 1,
                "creationTimestamp": "2017-12-05T07:08:39Z",
                "labels": {
                    "minipaas.io/appName": "",
                    "minipaas.io/name": "demo",
                    "minipaas.io/serviceName": "demo",
                    "replicas": "1"
                }
            },
            "spec": {
                "replicas": 1,
                "selector": {
                    "matchLabels": {
                        "minipaas.io/appName": "",
                        "minipaas.io/name": "demo",
                        "minipaas.io/serviceName": "demo",
                        "replicas": "1"
                    }
                },
                "template": {
                    "metadata": {
                        "creationTimestamp": null,
                        "labels": {
                            "minipaas.io/appName": "",
                            "minipaas.io/name": "demo",
                            "minipaas.io/serviceName": "demo",
                            "replicas": "1"
                        }
                    },
                    "spec": {
                        "containers": [
                            {
                                "name": "demo",
                                "image": "nginx:latest",
                                "ports": [
                                    {
                                        "containerPort": 80,
                                        "protocol": "TCP"
                                    }
                                ],
                                "resources": {
                                    "limits": {
                                        "cpu": "256m",
                                        "memory": "256Mi"
                                    },
                                    "requests": {
                                        "cpu": "256m",
                                        "memory": "256Mi"
                                    }
                                },
                                "terminationMessagePath": "/dev/termination-log",
                                "terminationMessagePolicy": "File",
                                "imagePullPolicy": "IfNotPresent"
                            }
                        ],
                        "restartPolicy": "Always",
                        "terminationGracePeriodSeconds": 30,
                        "dnsPolicy": "ClusterFirst",
                        "securityContext": {},
                        "schedulerName": "default-scheduler"
                    }
                },
                "strategy": {
                    "type": "RollingUpdate",
                    "rollingUpdate": {
                        "maxUnavailable": "25%",
                        "maxSurge": "25%"
                    }
                },
                "revisionHistoryLimit": 2,
                "progressDeadlineSeconds": 600
            },
            "status": {}
        },
        {
            "metadata": {
                "name": "demo",
                "namespace": "huangjia",
                "selfLink": "/api/v1/namespaces/huangjia/services/demo",
                "uid": "1eccc2ba-d98b-11e7-babc-525410406eaa",
                "resourceVersion": "8627958",
                "creationTimestamp": "2017-12-05T07:08:39Z",
                "labels": {
                    "minipaas.io/appName": "",
                    "minipaas.io/name": "demo",
                    "minipaas.io/serviceName": "demo"
                }
            },
            "spec": {
                "ports": [
                    {
                        "protocol": "TCP",
                        "port": 80,
                        "targetPort": 80,
                        "nodePort": 31189
                    }
                ],
                "selector": {
                    "minipaas.io/appName": "",
                    "minipaas.io/name": "demo",
                    "minipaas.io/serviceName": "demo"
                },
                "clusterIP": "10.100.29.209",
                "type": "NodePort",
                "sessionAffinity": "None",
                "externalTrafficPolicy": "Cluster"
            },
            "status": {
                "loadBalancer": {}
            }
        }
    ]
}
```


**删除应用**

- URL: 127.0.0.1:9090/api/v1/clusters/{clusterID}/namespaces/{namespace}/apps/{name}
 
- METHOD: DELETE

**请求：**
- 127.0.0.1:9090/api/v1/clusters/c2a523625ba526a7/namespaces/huangjia/apps/nginx


**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": "delete namespace huangjia's app named nginx success"
}
```


**删除应用**

- URL: 127.0.0.1:9090/api/v1/clusters/{clusterID}/namespaces/{namespace}/apps
 
- METHOD: GET

**请求：**
- 127.0.0.1:9090/api/v1/clusters/c2a523625ba526a7/namespaces/huangjia/apps


**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": {
        "apps": [
            {
                "id": 2,
                "createAt": "2017-12-05T19:25:18+08:00",
                "name": "nginx",
                "nameSpace": "huangjia",
                "description": "test app",
                "serviceCount": 1,
                "intanceCount": 1
            }
        ]
    }
}
```
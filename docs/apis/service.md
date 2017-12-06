**创建服务**

- URL: 127.0.0.1:9090/api/v1/clusters/{clusterID}/namespaces/{namespace}/apps/{appname}/services
 
- METHOD: POST

**请求：**
- 127.0.0.1:9090/api/v1/clusters/c2a523625ba526a7/namespaces/huangjia/apps/nginx/services

```
{
    "appName": "nginx",
    "name": "demo-1",
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
```

**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": [
        {
            "metadata": {
                "name": "demo-1",
                "namespace": "huangjia",
                "selfLink": "/apis/apps/v1beta1/namespaces/huangjia/deployments/demo-1",
                "uid": "7f6208be-d9b0-11e7-babc-525410406eaa",
                "resourceVersion": "8648843",
                "generation": 1,
                "creationTimestamp": "2017-12-05T11:36:13Z",
                "labels": {
                    "minipaas.io/appName": "nginx",
                    "minipaas.io/name": "demo-1",
                    "minipaas.io/serviceName": "demo-1",
                    "replicas": "1"
                }
            },
            "spec": {
                "replicas": 1,
                "selector": {
                    "matchLabels": {
                        "minipaas.io/appName": "nginx",
                        "minipaas.io/name": "demo-1",
                        "minipaas.io/serviceName": "demo-1",
                        "replicas": "1"
                    }
                },
                "template": {
                    "metadata": {
                        "creationTimestamp": null,
                        "labels": {
                            "minipaas.io/appName": "nginx",
                            "minipaas.io/name": "demo-1",
                            "minipaas.io/serviceName": "demo-1",
                            "replicas": "1"
                        }
                    },
                    "spec": {
                        "containers": [
                            {
                                "name": "demo-1",
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
                "name": "demo-1",
                "namespace": "huangjia",
                "selfLink": "/api/v1/namespaces/huangjia/services/demo-1",
                "uid": "7f5e56d0-d9b0-11e7-babc-525410406eaa",
                "resourceVersion": "8648841",
                "creationTimestamp": "2017-12-05T11:36:13Z",
                "labels": {
                    "minipaas.io/appName": "nginx",
                    "minipaas.io/name": "demo-1",
                    "minipaas.io/serviceName": "demo-1"
                }
            },
            "spec": {
                "ports": [
                    {
                        "protocol": "TCP",
                        "port": 80,
                        "targetPort": 80,
                        "nodePort": 30115
                    }
                ],
                "selector": {
                    "minipaas.io/appName": "nginx",
                    "minipaas.io/name": "demo-1",
                    "minipaas.io/serviceName": "demo-1"
                },
                "clusterIP": "10.99.72.17",
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


**删除服务**

- URL: 127.0.0.1:9090/api/v1/clusters/{clusterID}/namespaces/{namespace}/services/{name}
 
- METHOD: DELETE

**请求：**
- 127.0.0.1:9090/api/v1/clusters/c2a523625ba526a7/namespaces/huangjia/services/demo


**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": "delete namespace huangjia's service demo success"
}
```


**查询所有服务**

- URL: 127.0.0.1:9090/api/v1/clusters/{clusterID}/namespaces/{namespace}/services
 
- METHOD: GET

**请求：**
- 127.0.0.1:9090/api/v1/clusters/c2a523625ba526a7/namespaces/huangjia/services


**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": {
        "containers": [
            {
                "metadata": {
                    "name": "demo",
                    "namespace": "huangjia",
                    "selfLink": "/apis/apps/v1beta1/namespaces/huangjia/deployments/demo",
                    "uid": "dedb7fdf-d9b2-11e7-babc-525410406eaa",
                    "resourceVersion": "8650245",
                    "generation": 1,
                    "creationTimestamp": "2017-12-05T11:53:12Z",
                    "labels": {
                        "minipaas.io/appName": "",
                        "minipaas.io/name": "demo",
                        "minipaas.io/serviceName": "demo",
                        "replicas": "1"
                    },
                    "annotations": {
                        "deployment.kubernetes.io/revision": "1"
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
                "status": {
                    "observedGeneration": 1,
                    "replicas": 1,
                    "updatedReplicas": 1,
                    "readyReplicas": 1,
                    "availableReplicas": 1,
                    "conditions": [
                        {
                            "type": "Available",
                            "status": "True",
                            "lastUpdateTime": "2017-12-05T11:53:14Z",
                            "lastTransitionTime": "2017-12-05T11:53:14Z",
                            "reason": "MinimumReplicasAvailable",
                            "message": "Deployment has minimum availability."
                        },
                        {
                            "type": "Progressing",
                            "status": "True",
                            "lastUpdateTime": "2017-12-05T11:53:14Z",
                            "lastTransitionTime": "2017-12-05T11:53:12Z",
                            "reason": "NewReplicaSetAvailable",
                            "message": "ReplicaSet \"demo-2927577683\" has successfully progressed."
                        }
                    ]
                }
            }
        ],
        "services": [
            {
                "metadata": {
                    "name": "demo",
                    "namespace": "huangjia",
                    "selfLink": "/api/v1/namespaces/huangjia/services/demo",
                    "uid": "ded7f145-d9b2-11e7-babc-525410406eaa",
                    "resourceVersion": "8650218",
                    "creationTimestamp": "2017-12-05T11:53:12Z",
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
                            "nodePort": 30873
                        }
                    ],
                    "selector": {
                        "minipaas.io/appName": "",
                        "minipaas.io/name": "demo",
                        "minipaas.io/serviceName": "demo"
                    },
                    "clusterIP": "10.100.162.43",
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
}
```


**查询app所有服务**

- URL: 127.0.0.1:9090/api/v1/clusters/{clusterID}/namespaces/{namespace}/apps/{appname}/services
 
- METHOD: GET

**请求：**
- 127.0.0.1:9090/api/v1/clusters/c2a523625ba526a7/namespaces/huangjia/apps/nginx/services


**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": {
        "containers": [
            {
                "metadata": {
                    "name": "demo",
                    "namespace": "huangjia",
                    "selfLink": "/apis/apps/v1beta1/namespaces/huangjia/deployments/demo",
                    "uid": "dedb7fdf-d9b2-11e7-babc-525410406eaa",
                    "resourceVersion": "8650245",
                    "generation": 1,
                    "creationTimestamp": "2017-12-05T11:53:12Z",
                    "labels": {
                        "minipaas.io/appName": "",
                        "minipaas.io/name": "demo",
                        "minipaas.io/serviceName": "demo",
                        "replicas": "1"
                    },
                    "annotations": {
                        "deployment.kubernetes.io/revision": "1"
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
                "status": {
                    "observedGeneration": 1,
                    "replicas": 1,
                    "updatedReplicas": 1,
                    "readyReplicas": 1,
                    "availableReplicas": 1,
                    "conditions": [
                        {
                            "type": "Available",
                            "status": "True",
                            "lastUpdateTime": "2017-12-05T11:53:14Z",
                            "lastTransitionTime": "2017-12-05T11:53:14Z",
                            "reason": "MinimumReplicasAvailable",
                            "message": "Deployment has minimum availability."
                        },
                        {
                            "type": "Progressing",
                            "status": "True",
                            "lastUpdateTime": "2017-12-05T11:53:14Z",
                            "lastTransitionTime": "2017-12-05T11:53:12Z",
                            "reason": "NewReplicaSetAvailable",
                            "message": "ReplicaSet \"demo-2927577683\" has successfully progressed."
                        }
                    ]
                }
            }
        ],
        "services": []
    }
}
```
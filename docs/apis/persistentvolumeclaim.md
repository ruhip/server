**创建persistentvolumeclaim**

- URL: 127.0.0.1:9090/api/v1/clusters/{clusterID}/namespaces/{namespace}persistentvolumeclaims
 
- METHOD: POST

**请求：**
- 127.0.0.1:9090/api/v1/clusters/c2a523625ba526a7/namespaces/huangjia/persistentvolumeclaims

```
{
    "name":"huangjia",
    "size":"256Mi",
    "accessModes":"ReadWriteOnce"
}

```

**参数说明：**

- name 存储卷名称
- size 存储卷大小
- accessModes 存储卷属性， 默认为 ReadWriteOnce

**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": {
        "metadata": {
            "name": "huangjia",
            "namespace": "huangjia",
            "selfLink": "/api/v1/namespaces/huangjia/persistentvolumeclaims/huangjia",
            "uid": "7eb73c78-d967-11e7-babc-525410406eaa",
            "resourceVersion": "8608155",
            "creationTimestamp": "2017-12-05T02:53:38Z"
        },
        "spec": {
            "accessModes": [
                "ReadWriteOnce"
            ],
            "selector": {},
            "resources": {
                "limits": {
                    "storage": "256Mi"
                },
                "requests": {
                    "storage": "256Mi"
                }
            },
            "storageClassName": "minipaas.io.storageclass"
        },
        "status": {
            "phase": "Pending"
        }
    }
}
```


**删除persistentvolumeclaim**

- URL: 127.0.0.1:9090/api/v1/clusters/{clusterID}/namespaces/{namespace}persistentvolumeclaims/{name}
 
- METHOD: DELETE

**请求：**
- 127.0.0.1:9090/api/v1/clusters/c2a523625ba526a7/namespaces/huangjia/persistentvolumeclaims/huangjia


**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": "delete persistentvolumeclaim minipaas.io.storageclass success"
}
```


**查询persistentvolumeclaim**

- URL: 127.0.0.1:9090/api/v1/clusters/{clusterID}/namespaces/{namespace}persistentvolumeclaims
 
- METHOD: GET

**请求：**
- 127.0.0.1:9090/api/v1/clusters/c2a523625ba526a7/namespaces/huangjia/persistentvolumeclaims


**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": {
        "storages": [
            {
                "id": 14,
                "name": "huangjia2",
                "type": "rbd",
                "size": "256Mi",
                "accessModes": "ReadWriteOnce",
                "namespace": "huangjia",
                "used": false,
                "serviceName": "",
                "mountPath": ""
            }
        ]
    }
}
```
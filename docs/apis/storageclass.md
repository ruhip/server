**创建storagclass**

- URL: 127.0.0.1:9090/api/v1/clusters/{clusterID}/storageclasses
 
- METHOD: POST

**请求：**
- 127.0.0.1:9090/api/v1/clusters/c2a523625ba526a7/storageclasses

```
{
    "adminId": "admin",
    "adminSecretName": "ceph-secret",
    "adminSecretNamespace": "kube-system",
    "monitors": "10.39.0.114:6789,10.39.0.115:6789,10.39.0.116:6789",
    "pool": "tenx-pool",
    "userId": "admin",
    "userSecretName": "ceph-user-secret",
    "provisioner": "kubernetes.io/rbd"
}

```

**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": {
        "metadata": {
            "name": "minipaas.io.storageclass",
            "selfLink": "/apis/storage.k8s.io/v1/storageclasses/minipaas.io.storageclass",
            "uid": "b4a854f5-d8e3-11e7-babc-525410406eaa",
            "resourceVersion": "8538921",
            "creationTimestamp": "2017-12-04T11:10:15Z"
        },
        "provisioner": "kubernetes.io/rbd",
        "parameters": {
            "adminId": "admin",
            "adminSecretName": "ceph-secret",
            "adminSecretNamespace": "kube-system",
            "monitors": "10.39.0.114:6789,10.39.0.115:6789,10.39.0.116:6789",
            "pool": "tenx-pool",
            "userId": "admin",
            "userSecretName": "ceph-user-secret"
        }
    }
}
```


**删除storagclass**

- URL: 127.0.0.1:9090/api/v1/clusters/{clusterID}/storageclasses/{name}
 
- METHOD: DELETE

**请求：**
- 127.0.0.1:9090/api/v1/clusters/c2a523625ba526a7/storageclasses/minipaas.io.storageclass


**响应**

```
{
    "apiversion": "v1",
    "status": "200",
    "data": "delete storageclass minipaas.io.storageclass success"
}
```
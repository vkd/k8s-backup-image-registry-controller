k8s-backup-image-registry-controller
=

The purpose of this controller is to secure images (Deployment, DeamonSet) to a backup registry.

The controller watches images of Deployments and DeamonSets, copies images to a backup registry and reconfigures object's to use them.

How-to
-

Build and push docker image:
```sh
make docker-build docker-push IMG=<image_name>
```

Change image name (`docker-configuration`):
```
./config/default/deployment.yaml:

...
apiVersion: apps/v1
kind: Deployment
...
        containers:
        - name: manager
            image: <image_name>
            imagePullPolicy: Always
...

```

Change docker credentials in the k8s secret (`docker-configuration`):
```
./config/default/deployment.yaml:

...
apiVersion: v1
kind: Secret
metadata:
  name: docker-configuration
  namespace: backup-image-registry
stringData:
  config.json: |
    {
      "auths": {
        "https://index.docker.io/v1/": {
          "username": "<username>",
          "password": "<password>"
        }
      }
    }
...

```

Change configuration of controller:
```
./config/default/deployment.yaml:

...
apiVersion: v1
kind: ConfigMap
metadata:
  name: k8s-backup-image-registry-controller-configuration
  namespace: backup-image-registry
data:
  NAMESPACES_EXCLUDE_LIST: "kube-system<,namespace1,...>"
  REGISTRY_PREFIX: "<registry.domain/username>"
...

```

Create a controller deployment instance and the namespace (`backup-image-registry`):
```sh
kubectl apply -f ./config/default/
```



Alternative implementation
-

By using the admission webhooks. So far `kubebuilder` doesn't help much, probably it make sense to have a look later.


TODO
-

- [x] check if a new image version is reuploaded
- [ ] check admission webhooks
- [x] fix leader election
- [x] set correct service account
- [ ] do not upload image if already exists
- [ ] if error on update - return original image
- [x] configurable timeout
- [ ] to be able use custom namespace
- [ ] more convinient image's name parser
- [ ] use Golang's generics
- [ ] "kustomization" k8s manifests

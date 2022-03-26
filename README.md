k8s-backup-image-registry-controller
=

The purpose of this controller is to secure images (Deployment, DeamonSet) to a backup registry.

The controller watches images of Deployments and DeamonSets, copies images to a backup registry and reconfigures object's images to use them.

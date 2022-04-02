/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DeploymentObjectBuilder creates a Deployment object
type DeploymentObjectBuilder struct{}

func (r DeploymentObjectBuilder) newObject() *appsv1.Deployment { return &appsv1.Deployment{} }
func (r DeploymentObjectBuilder) NewObject() client.Object      { return r.newObject() }

func (r DeploymentObjectBuilder) Get(ctx context.Context, client client.Client, key types.NamespacedName) (ReconciledObject, error) {
	var d = &DeploymentObject{
		obj: r.newObject(),
	}
	err := client.Get(ctx, key, d.obj)
	return d, err
}

type DeploymentObject struct {
	obj *appsv1.Deployment
}

func (d *DeploymentObject) IsReady() bool {
	return d.obj.Status.Replicas != 0 && d.obj.Status.Replicas == d.obj.Status.ReadyReplicas
}

func (d *DeploymentObject) PodTemplate() v1.PodTemplateSpec { return d.obj.Spec.Template }

func (d *DeploymentObject) Update(ctx context.Context, client client.Client, template v1.PodTemplateSpec) error {
	return client.Update(ctx, d.obj)
}

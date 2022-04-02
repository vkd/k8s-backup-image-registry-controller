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

// DaemonSetObjectBuilder creates a DaemonSet object
type DaemonSetObjectBuilder struct{}

func (r DaemonSetObjectBuilder) newObject() *appsv1.DaemonSet { return &appsv1.DaemonSet{} }
func (r DaemonSetObjectBuilder) NewObject() client.Object     { return r.newObject() }

func (r DaemonSetObjectBuilder) Get(ctx context.Context, client client.Client, key types.NamespacedName) (ReconciledObject, error) {
	var d = &DaemonSetObject{
		obj: r.newObject(),
	}
	err := client.Get(ctx, key, d.obj)
	return d, err
}

type DaemonSetObject struct {
	obj *appsv1.DaemonSet
}

func (d *DaemonSetObject) IsReady() bool {
	return d.obj.Status.CurrentNumberScheduled != 0 && d.obj.Status.CurrentNumberScheduled == d.obj.Status.DesiredNumberScheduled
}

func (d *DaemonSetObject) PodTemplate() v1.PodTemplateSpec { return d.obj.Spec.Template }

func (d *DaemonSetObject) Update(ctx context.Context, client client.Client, template v1.PodTemplateSpec) error {
	return client.Update(ctx, d.obj)
}

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
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// BackupRegistryController reconciles an object
type BackupRegistryController struct {
	client.Client
	Scheme *runtime.Scheme

	NamespaceFilter interface {
		Check(namespace string) (ok bool)
	}
	BackupRegistry interface {
		CopyImage(_ context.Context, image string) (newImage string, ok bool, _ error)
	}
	RequeueAfter    time.Duration
	IgnoreReadiness bool

	ReconciledObjectBuilder ObjectBuilder
}

type ObjectBuilder interface {
	NewObject() client.Object
	Get(context.Context, client.Client, types.NamespacedName) (ReconciledObject, error)
}

type ReconciledObject interface {
	IsReady() bool
	PodTemplate() v1.PodTemplateSpec
	Update(context.Context, client.Client, v1.PodTemplateSpec) error
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Deployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *BackupRegistryController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	if !r.NamespaceFilter.Check(req.Namespace) {
		log.V(1).Info("Reconcile is ignored: namespace is from exclude list")
		return ctrl.Result{}, nil
	}

	obj, err := r.ReconciledObjectBuilder.Get(ctx, r.Client, req.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			log.V(1).Info("Reconcile is ignored: object not found")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("get object %+v: %w", req.NamespacedName, err)
	}

	if !r.IgnoreReadiness && !obj.IsReady() {
		log.V(1).Info("Reconcile is ignored: not ready")
		if r.RequeueAfter == 0 {
			r.RequeueAfter = 10 * time.Second
		}
		return ctrl.Result{RequeueAfter: r.RequeueAfter}, nil
	}

	var toUpdateObject bool

	template := obj.PodTemplate()
	for ci, c := range template.Spec.Containers {
		newImage, ok, err := r.BackupRegistry.CopyImage(ctx, c.Image)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("copy %q image: %w", c.Image, err)
		}

		if ok {
			log.V(1).Info(fmt.Sprintf("Image %q copied to new place %q", c.Image, newImage))

			template.Spec.Containers[ci].Image = newImage
			toUpdateObject = true
		}
	}

	if toUpdateObject {
		err = obj.Update(ctx, r.Client, template)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("update object: %w", err)
		}

		log.V(1).Info("Images are updated")
	} else {
		log.V(1).Info("Reconcile is ignored: nothing to update")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BackupRegistryController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(r.ReconciledObjectBuilder.NewObject()).
		Complete(r)
}

package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("DaemonSet controller", func() {
	const (
		daemonSetName      = "test-daemonset"
		daemonSetNamespace = "default"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating DaemonSet", func() {

		It("Should change DaemonSet Container.Image", func() {
			ctx := context.Background()

			obj := &appsv1.DaemonSet{
				TypeMeta: metav1.TypeMeta{
					// APIVersion: "batch.tutorial.kubebuilder.io/v1",
					Kind: "DaemonSet",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      daemonSetName,
					Namespace: daemonSetNamespace,
				},
				Spec: appsv1.DaemonSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "my-service",
						},
					},
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "my-service",
							},
						},
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  "my-service",
									Image: "golang:alpine",
								},
								{
									Name:  "my-service-2",
									Image: "backup.registry/username/postgres",
								},
								{
									Name:  "my-service-3",
									Image: "external.registry/library/redis",
								},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, obj)).Should(Succeed())

			Eventually(func() ([]string, error) {
				obj := &appsv1.DaemonSet{}
				err := k8sClient.Get(ctx, types.NamespacedName{
					Name:      daemonSetName,
					Namespace: daemonSetNamespace,
				}, obj)
				if err != nil {
					return nil, err
				}

				l := make([]string, 0, len(obj.Spec.Template.Spec.Containers))
				for _, c := range obj.Spec.Template.Spec.Containers {
					l = append(l, c.Image)
				}

				return l, nil
			}, timeout, interval).Should(Equal([]string{
				"backup.registry/username/golang:alpine",
				"backup.registry/username/postgres",
				"backup.registry/username/external-registry-library-redis:latest",
			}))
		})
	})
})

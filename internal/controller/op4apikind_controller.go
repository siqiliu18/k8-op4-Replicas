/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	op4apigroupv1 "op4repo/api/v1"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Op4ApiKindReconciler reconciles a Op4ApiKind object
type Op4ApiKindReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=op4apigroup.op4domain,resources=op4apikinds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=op4apigroup.op4domain,resources=op4apikinds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=op4apigroup.op4domain,resources=op4apikinds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Op4ApiKind object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *Op4ApiKindReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the PodSet instance
	podSet := &op4apigroupv1.Op4ApiKind{}
	err := r.Get(context.TODO(), req.NamespacedName, podSet)
	if err != nil {
		log.Error(err, "unable to fetch podSet")
		return ctrl.Result{}, err
	}

	// List all pods owned by this PodSet instance
	lbls := labels.Set{
		"app":     podSet.Name,
		"version": "v0.1",
	}
	existingPods := &corev1.PodList{}
	err = r.List(context.TODO(),
		existingPods,
		&client.ListOptions{
			Namespace:     req.Namespace,
			LabelSelector: labels.SelectorFromSet(lbls),
		})
	if err != nil {
		log.Error(err, "failed to list existing pods in the podSet")
		return ctrl.Result{}, err
	}
	existingPodNames := []string{}
	// Count the pods that are pending or running as available
	for _, pod := range existingPods.Items {
		if pod.GetObjectMeta().GetDeletionTimestamp() != nil {
			continue
		}
		if pod.Status.Phase == corev1.PodPending || pod.Status.Phase == corev1.PodRunning {
			existingPodNames = append(existingPodNames, pod.GetObjectMeta().GetName())
		}
	}

	fmt.Printf("current pods %v\n", existingPods.Items)

	log.Info("Checking podset", "expected replicas", podSet.Spec.Replicas, "Pod.Names", existingPodNames)
	// Update the status if necessary
	status := op4apigroupv1.Op4ApiKindStatus{
		Replicas: int32(len(existingPodNames)),
		PodNames: existingPodNames,
	}
	if !reflect.DeepEqual(podSet.Status, status) {
		log.Info("[Current State and Desired State are different]", "current status = ", podSet.Status, "desired status = ", status)
		podSet.Status = status
		err := r.Status().Update(context.TODO(), podSet)
		if err != nil {
			log.Error(err, "failed to update the podSet")
			return ctrl.Result{}, err
		}
	}

	fmt.Println()
	fmt.Printf("111. current running pods = %v\n", existingPodNames)
	fmt.Println()
	fmt.Printf("222. configured replicas = %v\n", podSet.Spec.Replicas)
	fmt.Println()

	// Scale Down Pods
	if int32(len(existingPodNames)) > podSet.Spec.Replicas {
		// delete a pod. Just one at a time (this reconciler will be called again afterwards)
		log.Info("Deleting a pod in the podset", "expected replicas", podSet.Spec.Replicas, "Pod.Names", existingPodNames)
		pod := existingPods.Items[0]
		err = r.Delete(context.TODO(), &pod)
		if err != nil {
			log.Error(err, "failed to delete a pod")
			return ctrl.Result{}, err
		}
	}

	// Scale Up Pods
	if int32(len(existingPodNames)) < podSet.Spec.Replicas {
		// create a new pod. Just one at a time (this reconciler will be called again afterwards)
		log.Info("Adding a pod in the podset", "expected replicas", podSet.Spec.Replicas, "Pod.Names", existingPodNames)
		pod := newPodForCR(podSet)
		if err := controllerutil.SetControllerReference(podSet, pod, r.Scheme); err != nil {
			log.Error(err, "unable to set owner reference on new pod")
			return ctrl.Result{}, err
		}
		err = r.Create(context.TODO(), pod)
		if err != nil {
			log.Error(err, "failed to create a pod")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *op4apigroupv1.Op4ApiKind) *corev1.Pod {
	labels := map[string]string{
		"app":     cr.Name,
		"version": "v0.1",
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: cr.Name + "-pod",
			Namespace:    cr.Namespace,
			Labels:       labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *Op4ApiKindReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&op4apigroupv1.Op4ApiKind{}).
		Complete(r)
}

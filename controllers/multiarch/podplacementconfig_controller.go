/*
Copyright 2023.

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

package multiarch

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	multiarchv1alpha1 "multiarch-operator/apis/multiarch/v1alpha1"
)

// PodPlacementConfigReconciler reconciles a PodPlacementConfig object
type PodPlacementConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=multiarch.openshift.io,resources=podplacementconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=multiarch.openshift.io,resources=podplacementconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=multiarch.openshift.io,resources=podplacementconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodPlacementConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *PodPlacementConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Lookup the PodPlacementConfig instance for this reconcile request
	podplacementconfig := &multiarchv1alpha1.PodPlacementConfig{}
	if err := r.Get(ctx, types.NamespacedName{Name: "podplacementconfig-sample", Namespace: ""}, podplacementconfig); err != nil {
		klog.Errorf("unable to fetch PodPlacementConfig %s: %v", req.Name, err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	defaultNamespaces := sets.NewString("openshift-*", "kube-*", "hypershift-*")
	tmpDefaultSet := defaultNamespaces.Difference(sets.NewString(podplacementconfig.Spec.ExcludedNamespaces...))
	podplacementconfig.Spec.ExcludedNamespaces = append(podplacementconfig.Spec.ExcludedNamespaces, tmpDefaultSet.List()...)
	err := r.Client.Update(ctx, podplacementconfig)
	if err != nil {
		klog.Errorf("unable to update the podplacementconfig %s: %v", podplacementconfig.Name, err)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodPlacementConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&multiarchv1alpha1.PodPlacementConfig{}).
		Complete(r)
}

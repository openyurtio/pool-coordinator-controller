/*
Copyright 2022 The OpenYurt Authors.

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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1beta1 "github.com/openyurtio/coordinator-controller/api/v1beta1"
)

// PoolCoordinatorReconciler reconciles a PoolCoordinator object
type PoolCoordinatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apiextensions.k8s.io.apps.openyurt.io,resources=poolcoordinators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apiextensions.k8s.io.apps.openyurt.io,resources=poolcoordinators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apiextensions.k8s.io.apps.openyurt.io,resources=poolcoordinators/finalizers,verbs=update

func (r *PoolCoordinatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var p v1beta1.PoolCoordinator
	if err := r.Get(ctx, req.NamespacedName, &p); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 1. Handle the pool coordinator creation event
	//		1.1 check whether it can be created: the number of nodes in the NodePool is less than 3, or if a pool-coordinator has been deployed in the NodePool.
	//		1.2 prepares the tls server certificate for pool-coordinator, saves the certificate in secret and mounts it to pool-coordinator
	//		1.3 create yurtAppSet for pool coordinator.
	//		1.4 generates kubeconfig for users to access pool-coordinator.

	// 2. Handle the pool coordinator update event
	//		2.1 Update the pool coordinator status if necessary

	// 3. Handle the pool coordinator deletion event
	//		3.1 check whether it can be deleted.
	//		3.2 delete yurtAppSet for pool coordinator.
	//		3.3 cleans up the certificates of pool-coordinator(tls server certificate and kubeconfig).

	return ctrl.Result{}, nil
}


// SetupWithManager sets up the controller with the Manager.
func (r *PoolCoordinatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.PoolCoordinator{}).
		Complete(r)
}

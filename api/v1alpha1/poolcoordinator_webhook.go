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

package v1alpha1

import (
	"context"
	"fmt"
	yurtapi "github.com/openyurtio/yurt-app-manager-api/pkg/yurtappmanager/apis/apps/v1alpha1"
	appsv1alpha1 "github.com/openyurtio/yurt-app-manager-api/pkg/yurtappmanager/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var klog = log.Log.WithName("poolcoordinator-webhook")

func (r *PoolCoordinator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

const (
	NODE_AUTONOMY string = "node.beta.openyurt.io/autonomy"
)

const (
	MIN_NODES_REQUIRED int = 3
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//kubebuilder:webhook:path=/mutate-pool-coordinator-openyurt-io-v1alpha1-poolcoordinator,mutating=true,failurePolicy=fail,sideEffects=None,groups=pool-coordinator.openyurt.io,resources=poolcoordinators,verbs=create;update,versions=v1alpha1,name=mpoolcoordinator.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &PoolCoordinator{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *PoolCoordinator) Default() {
	klog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-pool-coordinator-openyurt-io-v1alpha1-poolcoordinator,mutating=false,failurePolicy=fail,sideEffects=None,groups=pool-coordinator.openyurt.io,resources=poolcoordinators,verbs=create;update,versions=v1alpha1,name=vpoolcoordinator.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &PoolCoordinator{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *PoolCoordinator) ValidateCreate() error {
	klog.Info("validate create", "name", r.Name)
	return r.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *PoolCoordinator) ValidateUpdate(old runtime.Object) error {
	klog.Info("validate update", "name", r.Name)
	return r.validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *PoolCoordinator) ValidateDelete() error {
	klog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *PoolCoordinator) validate() error {
	var errmsg string
	// get all nodepools
	config := ctrl.GetConfigOrDie()
	clientset, err := appsv1alpha1.NewForConfig(config)
	if err != nil {
		return err
	}
	nps, err := clientset.AppsV1alpha1().NodePools().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return err
	}

	// reject the request if no nodepool in the cluster
	if len(nps.Items) <= 0 {
		errmsg = "No nodepool is created in the cluster!"
		err = fmt.Errorf(errmsg)
		return err
	}

	npSpec := r.Spec.NodePool
	var specNpObj yurtapi.NodePool
	var found = false
	// verify if the nodepool name specified in the Spec
	for _, np := range nps.Items {
		if np.ObjectMeta.Name == npSpec {
			specNpObj = np
			found = true
			break
		}
	}
	if !found {
		errmsg = "Nodepool " + npSpec + " does not exist in the cluster!"
		err = fmt.Errorf(errmsg)
		return err
	}
	// reject the request if a nodepool has less than 3 nodes
	if len(specNpObj.Status.Nodes) <= MIN_NODES_REQUIRED {
		errmsg = "Nodepool Autonomy will be not enabled for the nodepool with 3 or less nodes!"
		err = fmt.Errorf(errmsg)
		return err
	}

	// reject the request if any node in the nodepool has node autonomy feature enabled
	for _, node := range specNpObj.Status.Nodes {
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return err
		}
		nodeObj, err := clientset.CoreV1().Nodes().Get(context.TODO(), node, v1.GetOptions{})
		if err != nil {
			return err
		}
		annotations := nodeObj.ObjectMeta.Annotations
		klog.Info("node annotation", NODE_AUTONOMY, annotations[NODE_AUTONOMY])
		if annotations[NODE_AUTONOMY] == "true" {
			errmsg = "Node Autonomy and Nodepool Autonomy should not be enabled simultaneously! Nodepool autonomy feature will not be enabled for the Node " + node + " which already has node autonomy feature enabled!"
			err = fmt.Errorf(errmsg)
			return err
		}

	}

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

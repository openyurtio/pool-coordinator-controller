/*
Copyright 2018 The Kubernetes Authors.

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
	"net/http"

	"github.com/openyurtio/pool-coordinator/api/v1alpha1"
	appsv1alpha1 "github.com/openyurtio/yurt-app-manager-api/pkg/yurtappmanager/client/clientset/versioned"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var klog = log.Log.WithName("node-autonomy-webhook")

//+kubebuilder:webhook:path=/validate-node-autonomy-openyurt-io-v1alpha1-poolcoordinator,mutating=false,failurePolicy=fail,sideEffects=None,groups="",resources=nodes,verbs=create;update,versions=v1,name=vnode.kb.io,admissionReviewVersions={v1,v1beta1}

// NodeValidator validates Nodes
type NodeValidator struct {
	Client client.Client
	// Decoder decodes objects
	decoder *admission.Decoder
}

const (
	NODE_POOL_JOIN string = "apps.openyurt.io/desired-nodepool"
)

// Node autonomy and nodepool autonomy should not be enable at the same time,
// so when creating node autonomy annotation we need to check if nodepool autonomy
// has been enabled first, if yes, we reject the annotation creation
func (v *NodeValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	klog.Info("enter NodeValidator")

	node := &corev1.Node{}
	err := v.decoder.Decode(req, node)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	coordinators := v1alpha1.PoolCoordinatorList{}
	if err = v.Client.List(context.TODO(), &coordinators, &client.ListOptions{}); err != nil {
		return admission.Denied(err.Error())
	}

	// if we find NODE_AUTONOMY annotation in the request,
	// we check if it already set for the node before,
    // if yes, we check if we are joining the node to a nodepool with nodepool autonomy enbled
	// if we don't set NODE_AUTONOMY before, which means we are annotating it for this node now,
	// then we check if this node belongs to some nodepool
	// which has nodepool autonomy feature enabled, if yes, reject the request
	anno_req, found_req := node.Annotations[v1alpha1.NODE_AUTONOMY]
	if found_req && anno_req == "true" {
		// 1. check the node setting before
		if req.AdmissionRequest.Operation == admissionv1.Update {
			oldnode := &corev1.Node{}
			if err := v.decoder.DecodeRaw(req.OldObject, oldnode); err != nil {
				return admission.Errored(http.StatusBadRequest, err)
			}
			anno_old, found_old := oldnode.Annotations[v1alpha1.NODE_AUTONOMY]
			if found_old && anno_old == "true" {
				np_label, np_found := node.Labels[NODE_POOL_JOIN]
				if !np_found {
					return admission.Allowed("")
				}
				for _, c := range coordinators.Items {
					if c.Spec.NodePool == np_label {
						errmsg := "Node Autonomy and Nodepool Autonomy should not be enabled simultaneously! Node " + node.Name + " already has node autonomy enable is joining the nodepool " + np_label + " which has nodepool autonomy feature enabled!"
						err = fmt.Errorf(errmsg)
						return admission.Denied(errmsg)
					}
				}
			}
		}

		// 2. check nodepool and if its nodepool automomy feature enabled
		config := ctrl.GetConfigOrDie()
		clientset, err := appsv1alpha1.NewForConfig(config)
		if err != nil {
			return admission.Denied(err.Error())
		}
		nps, err := clientset.AppsV1alpha1().NodePools().List(context.TODO(), v1.ListOptions{})
		if err != nil {
			return admission.Denied(err.Error())
		}

		if len(coordinators.Items) <= 0 || len(nps.Items) <= 0 {
			return admission.Allowed("")
		}

		for i, c := range coordinators.Items {
			klog.Info(fmt.Sprint(i) + ": go through poolcoordinator: " + c.ObjectMeta.Name + ", its nodepool: " + c.Spec.NodePool)
			for j, np := range nps.Items {
				klog.Info(fmt.Sprint(j) + ": go through np: " + np.ObjectMeta.Name)
				if c.Spec.NodePool == np.ObjectMeta.Name {
					for k, n := range np.Status.Nodes {
						klog.Info(fmt.Sprint(k) + ": go through node: " + n)
						if n == node.Name {
							errmsg := "Node Autonomy and Nodepool Autonomy should not be enabled simultaneously! Node " + node.Name + " joins nodepool " + np.ObjectMeta.Name + " which has nodepool autonomy feature enabled already by " + c.ObjectMeta.Name
							err = fmt.Errorf(errmsg)
							return admission.Denied(errmsg)

						}
					}
					break
				}
			}
		}

		//return admission.Denied(fmt.Sprintf("missing annotation %s", key))
	}

	return admission.Allowed("")
}

// NodeValidator implements admission.DecoderInjector.
// A decoder will be automatically injected.

//var _ admission.DecoderInjector = &NodeValidator{}

// InjectDecoder injects the decoder.
func (v *NodeValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}

/*var _ inject.Client = &NodeValidator{}

// InjectClient injects the client into the NodeValidator
func (v *NodeValidator) InjectClient(c client.Client) error {
	v.Client = c
	return nil
}*/

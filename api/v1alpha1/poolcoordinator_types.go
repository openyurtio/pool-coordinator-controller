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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PoolCoordinator CRD
/*type PoolCoordinator struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec PoolCoordinatorSpec
	Status PoolCoordinatorStatus
}*/

type PoolCoordinatorSpec struct {
	// Version of pool-coordinator, which corresponding to the Kubernetes version
	Version string `json:"version,omitempty"`
	// The NodePool managed by pool-coordinator.
	NodePool string `json:"nodepool,omitempty"`
}

type PoolCoordinatorStatus struct {
	// The node where pool-coordinator is located.
	NodeName string `json:"nodeName,omitempty"`
	// Conditions represent the status of pool-coordinator, which is filled by the coordinator-controller.
	Conditions []PoolCoordinatorCondition `json:"conditions,omitempty"`
	// DelegatedNodes are the nodes in the node pool that are disconnected from the cloud.
	DelegatedNodes []string `json:"delegatedNodes,omitempty"`
	// OutsidePoolNodes are nodes in the node pool that cannot connect to pool-coordinator.
	OutsidePoolNodes []string `json:"outsidePoolNodes,omitempty"`
}

type PoolCoordinatorCondition struct {
	Type               PoolCoordinatorConditionType `json:"type,omitempty"`
	Status             ConditionStatus              `json:"status,omitempty"`
	LastProbeTime      metav1.Time                  `json:"lastProbeTime,omitempty"`
	LastTransitionTime metav1.Time                  `json:"lastTransitionTime,omitempty"`
	Reason             string                       `json:"reason,omitempty"`
	Message            string                       `json:"message,omitempty"`
}

type PoolCoordinatorConditionType string

const (
	// PoolCoordinatorPending indicates that the deployment of pool-coordinator is blocked.
	//This happens, for example, if the number of nodes in the node pool is less than 3.
	PoolCoordinatorPending PoolCoordinatorConditionType = "Pending"
	// PoolCoordinatorCertsReady indicates that the certificate used by pool-coordinator is ready.
	PoolCoordinatorCertsReady PoolCoordinatorConditionType = "CertsReady"
	// PoolCoordinatorReady indicates that pool-coordinator is ready.
	PoolCoordinatorReady PoolCoordinatorConditionType = "Ready"
)

type ConditionStatus string

const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PoolCoordinator is the Schema for the poolcoordinators API
type PoolCoordinator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PoolCoordinatorSpec   `json:"spec,omitempty"`
	Status PoolCoordinatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PoolCoordinatorList contains a list of PoolCoordinator
type PoolCoordinatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PoolCoordinator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PoolCoordinator{}, &PoolCoordinatorList{})
}

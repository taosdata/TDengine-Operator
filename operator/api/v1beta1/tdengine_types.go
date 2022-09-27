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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TDengineSpec defines the desired state of TDengine
type TDengineSpec struct {
	Replicas        *int32            `json:"replicas,omitempty"`
	Image           string            `json:"image"`
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	Env             []corev1.EnvVar   `json:"env,omitempty"`
	// +kubebuilder:validation:Optional
	PodResources         corev1.ResourceRequirements    `json:"podResources,omitempty"`
	VolumeClaimTemplates []corev1.PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`
}

// TDengineStatus defines the observed state of TDengine
type TDengineStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// TDengine is the Schema for the tdengines API
type TDengine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TDengineSpec   `json:"spec,omitempty"`
	Status TDengineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TDengineList contains a list of TDengine
type TDengineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TDengine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TDengine{}, &TDengineList{})
}

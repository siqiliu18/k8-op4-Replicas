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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Op4ApiKindSpec defines the desired state of Op4ApiKind
type Op4ApiKindSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Op4ApiKind. Edit op4apikind_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// Op4ApiKindStatus defines the observed state of Op4ApiKind
type Op4ApiKindStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Op4ApiKind is the Schema for the op4apikinds API
type Op4ApiKind struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Op4ApiKindSpec   `json:"spec,omitempty"`
	Status Op4ApiKindStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// Op4ApiKindList contains a list of Op4ApiKind
type Op4ApiKindList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Op4ApiKind `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Op4ApiKind{}, &Op4ApiKindList{})
}

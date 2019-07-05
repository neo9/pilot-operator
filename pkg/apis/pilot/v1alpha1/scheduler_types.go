package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SchedulerRequest struct {
	Name string `json:"name"`
	Method string `json:"method"`
	Schedule string `json:"schedule"`
	Service string `json:"service"`
	Port int32 `json:"port"`
	Path string `json:"path"`
}

// SchedulerSpec defines the desired state of Scheduler
// +k8s:openapi-gen=true
type SchedulerSpec struct {
	Requests []SchedulerRequest `json:"requests"`
}

// SchedulerStatus defines the observed state of Scheduler
// +k8s:openapi-gen=true
type SchedulerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Scheduler is the Schema for the schedulers API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Scheduler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SchedulerSpec   `json:"spec,omitempty"`
	Status SchedulerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SchedulerList contains a list of Scheduler
type SchedulerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Scheduler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Scheduler{}, &SchedulerList{})
}

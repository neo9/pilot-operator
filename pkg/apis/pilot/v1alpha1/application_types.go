package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ApplicationType string

const (
	JAVA ApplicationType = "java"
	WEB = "web"
	NODEJS = "nodejs"
)

type ApplicationSecret struct {
	Name string `json:"name"`
	Key string `json:"key"`
}

type ApplicationPod struct {
	Secrets []ApplicationSecret `json:"secrets"`
}

type ApplicationHealthCheck struct {
	Path string `json:"path"`
}

type ApplicationResources struct {
	Memory string `json:"memory"`
	CPU string `json:"cpu"`
}

type ApplicationService struct {
	Port int32 `json:"port,omitempty"`
	TargetPort int32 `json:"targetPort,omitempty"`
}

// ApplicationSpec defines the desired state of Application
// +k8s:openapi-gen=true
type ApplicationSpec struct {
	Replicas int32 `json:"replicas,omitempty"`
	Type ApplicationType `json:"type"`
	Repository string `json:"repository"`
	Tag string `json:"tag,omitempty"`
	InitTag string `json:"initTag,omitempty"`
	Service ApplicationService `json:"service,omitempty"`
	Resources ApplicationResources `json:"resources,omitempty"`
	HealthCheck ApplicationHealthCheck `json:"healthCheck,omitempty"`
	Pod ApplicationPod `json:"pod,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}


// ApplicationStatus defines the observed state of Application
// +k8s:openapi-gen=true
type ApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Application is the Schema for the applications API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApplicationList contains a list of Application
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}

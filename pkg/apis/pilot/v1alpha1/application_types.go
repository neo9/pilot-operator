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

type ApplicationCronJobRequest struct {
	Name string `json:"name"`
	Method string `json:"method"`
	Path string `json:"path"`
	Schedule string `json:"schedule"`
}

type ApplicationCronJob struct {
	Requests []ApplicationCronJobRequest `json:"requests"`
}

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
	Port int32 `json:"port"`
	TargetPort int32 `json:"targetPort"`
}

// ApplicationSpec defines the desired state of Application
// +k8s:openapi-gen=true
// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
type ApplicationSpec struct {
	Replicas int32 `json:"replicas"`
	Type ApplicationType `json:"type"`
	Repository string `json:"repository"`
	Tag string `json:"tag"`
	InitTag string `json:"initTag"`
	Service ApplicationService `json:"service"`
	Resources ApplicationResources `json:"resources"`
	HealthCheck ApplicationHealthCheck `json:"healthCheck"`
	Pod ApplicationPod `json:"pod"`
	CronJob ApplicationCronJob `json:"cronjob"`
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

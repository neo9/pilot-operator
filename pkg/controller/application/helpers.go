package application

import (
	"github.com/go-logr/logr"
)

func getLogger(namespace string, name string, kind string) logr.Logger {
	return log.WithValues("ApplicationKind", kind, "Request.Namespace", namespace, "Request.Name", name)
}

func getMergedLabels(originalLabels map[string]string, labels map[string]string) map[string]string {
	for key, value := range labels {
		// Do not override original labels
		if _, ok := originalLabels[key]; !ok {
			originalLabels[key] = value
		}
	}

	return originalLabels
}

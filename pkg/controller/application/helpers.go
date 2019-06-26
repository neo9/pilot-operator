package application

import (
	"github.com/go-logr/logr"
)

func getLogger(namespace string, name string, kind string) logr.Logger {
	return log.WithValues("ApplicationKind", kind, "Request.Namespace", namespace, "Request.Name", name)
}

package controller

import (
	"github.com/neo9/pilot-operator/pkg/controller/scheduler"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, scheduler.Add)
}

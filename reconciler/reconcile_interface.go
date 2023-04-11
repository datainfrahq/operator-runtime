package reconciler

import (
	"github.com/datainfrahq/operator-builder/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileInterface holds all the methods to create operators
type ReconcileInterface interface {
	ReconcileConfigMap() (controllerutil.OperationResult, error)
	ReconcileDeployOrSts() (controllerutil.OperationResult, error)
	ReconcileStorage() (controllerutil.OperationResult, error)
	ReconcileService() (controllerutil.OperationResult, error)
	ReconcileStore() error
}

var Reconciler ReconcileInterface = builder.NewBuilder()

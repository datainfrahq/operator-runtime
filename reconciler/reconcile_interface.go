package reconciler

import (
	"github.com/datainfrahq/operator-builder/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileInterface holds all the methods to create operators
type ReconcileInterface interface {
	ReconcileConfigMap() (controllerutil.OperationResult, error)
	ReconcileConfigMapHash() ([]builder.HashHolder, error)
	ReconcileDeployOrSts(cmHashes []builder.HashHolder) (controllerutil.OperationResult, error)
	ReconcileStorage() (controllerutil.OperationResult, error)
	ReconcileService() (controllerutil.OperationResult, error)
}

var Reconciler ReconcileInterface = builder.NewBuilder()

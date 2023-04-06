package builder

import "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

// BuilderInterface holds all the methods to create operators
type BuilderInterface interface {
	BuildConfigMap() (controllerutil.OperationResult, error)
	BuildConfigMapHash() ([]HashHolder, error)
	BuildDeployOrSts(cmHashes []HashHolder) (controllerutil.OperationResult, error)
	BuildDeployment(cmhashes []HashHolder) (controllerutil.OperationResult, error)
	BuildPvc() (controllerutil.OperationResult, error)
	BuildStatefulset(cmhashes []HashHolder) (controllerutil.OperationResult, error)
}

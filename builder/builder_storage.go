package builder

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type BuilderStorageConfig struct {
	PvcSpec *v1.PersistentVolumeClaimSpec
	CommonBuilder
}

func ToNewBuilderStorageConfig(builder []BuilderStorageConfig) func(*Builder) {
	return func(s *Builder) {
		s.StorageConfig = builder
	}
}

func (s *Builder) ReconcileStorage() (controllerutil.OperationResult, error) {

	for _, storage := range s.StorageConfig {

		cm, err := storage.MakePvc()
		if err != nil {
			return controllerutil.OperationResultNone, err
		}

		storage.DesiredState = cm
		storage.CurrentState = &v1.PersistentVolumeClaim{}

		_, err = storage.CreateOrUpdate(s.Context.Context, s.Recorder)
		if err != nil {
			return controllerutil.OperationResultNone, nil
		}

	}

	return controllerutil.OperationResultNone, nil
}

func (b BuilderStorageConfig) MakePvc() (*v1.PersistentVolumeClaim, error) {
	return &v1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "PersistentVolumeClaim",
		},
		ObjectMeta: b.ObjectMeta,
		Spec:       *b.PvcSpec,
	}, nil
}

package builder

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type BuilderConfigMap struct {
	Data map[string]string
	CommonBuilder
}

func ToNewBuilderConfigMap(builder []BuilderConfigMap) func(*Builder) {
	return func(s *Builder) {
		s.ConfigMaps = builder
	}
}

func (s *Builder) ReconcileConfigMap() (controllerutil.OperationResult, error) {

	var result controllerutil.OperationResult

	for _, configMap := range s.ConfigMaps {

		cm, err := configMap.makeConfigMap()
		if err != nil {
			return controllerutil.OperationResultNone, err
		}

		s.Put(cm.GetName(), cm.Kind)

		configMap.DesiredState = cm
		configMap.CurrentState = &v1.ConfigMap{}

		result, err = configMap.CreateOrUpdate(s.Context.Context, s.Recorder)
		if err != nil {
			return result, nil
		}
	}

	return result, nil
}

func (b *BuilderConfigMap) makeConfigMap() (*v1.ConfigMap, error) {
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: b.ObjectMeta,
		Data:       b.Data,
	}, nil
}

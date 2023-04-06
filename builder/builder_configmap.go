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

func (s *Builder) BuildConfigMap() (controllerutil.OperationResult, error) {

	for _, configMap := range s.ConfigMaps {

		cm, err := configMap.MakeConfigMap()
		if err != nil {
			return controllerutil.OperationResultNone, err
		}

		configMap.DesiredState = cm
		configMap.CurrentState = &v1.ConfigMap{}

		_, err = configMap.CreateOrUpdate(s.Context.Context, s.Recorder)
		if err != nil {
			return controllerutil.OperationResultNone, nil
		}

	}

	return controllerutil.OperationResultNone, nil
}

func (b BuilderConfigMap) MakeConfigMap() (*v1.ConfigMap, error) {
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: b.ObjectMeta,
		Data:       b.Data,
	}, nil
}

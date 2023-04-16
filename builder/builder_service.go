package builder

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type BuilderService struct {
	ServiceSpec    *v1.ServiceSpec
	SelectorLabels map[string]string
	CommonBuilder
}

func ToNewBuilderService(builder []BuilderService) func(*Builder) {
	return func(s *Builder) {
		s.Service = builder
	}
}
func (s *Builder) ReconcileService() (controllerutil.OperationResult, error) {

	var err error
	var result controllerutil.OperationResult

	for _, svc := range s.Service {

		if svc.ServiceSpec != nil {

			makeSvc := svc.makeService()

			svc.DesiredState = makeSvc
			svc.CurrentState = &v1.Service{}

			result, err = svc.CreateOrUpdate(s.Context.Context, s.Recorder)
			if err != nil {
				return controllerutil.OperationResultNone, nil
			}
		}
	}
	return result, nil
}

func (b *BuilderService) makeService() *v1.Service {
	svc := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: b.ObjectMeta,
		Spec:       *b.ServiceSpec,
	}

	svc.Spec.Selector = b.SelectorLabels

	return svc
}

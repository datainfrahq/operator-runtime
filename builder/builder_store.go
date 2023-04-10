package builder

import (
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8sObjectName string

const (
	// As per k8s naming
	configMap  K8sObjectName = "ConfigMap"
	deployment K8sObjectName = "Deployment"
	pvc        K8sObjectName = "PersistentVolumeClaim"
	svc        K8sObjectName = "Service"
)

type InternalStore struct {
	ObjectNameKind map[string]string
	CommonBuilder
}

func NewStore(
	client client.Client,
	labels map[string]string,
	namespace string,
	crObject client.Object,
) *InternalStore {
	return &InternalStore{
		ObjectNameKind: make(map[string]string),
		CommonBuilder: CommonBuilder{
			Client: client,
			Labels: labels,
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
			},
			CrObject: crObject,
		},
	}
}

func ToNewBuilderStore(builder InternalStore) func(*Builder) {
	return func(s *Builder) {
		s.Store = builder
	}
}

func (s *Builder) Put(key, value string) {
	if _, isKeyExists := s.Store.ObjectNameKind[key]; isKeyExists {
		return
	} else {
		s.Store.ObjectNameKind[key] = value
	}
}

func (s *Builder) Exists(key string) bool {
	if _, isKeyExists := s.Store.ObjectNameKind[key]; isKeyExists {
		return true
	}
	return false
}

func (s *Builder) ReconcileStore() error {

	for _, kind := range s.Store.ObjectNameKind {
		switch kind {
		case string(deployment):
			s.Store.CommonBuilder.ObjectList = &v1.DeploymentList{}
			list, err := s.Store.List(s.Context.Context, s.Recorder)
			if err != nil {
				return err
			}
			for _, deployment := range list.(*v1.DeploymentList).Items {
				if !s.Exists(deployment.GetName()) {
					s.Store.CommonBuilder.DesiredState = &deployment

					_, err := s.Store.Delete(s.Context.Context, s.Recorder)
					if err != nil {
						return err
					}
				}
			}
		case string(configMap):
			s.Store.CommonBuilder.ObjectList = &corev1.ConfigMapList{}
			list, err := s.Store.List(s.Context.Context, s.Recorder)
			if err != nil {
				return err
			}
			for _, cm := range list.(*corev1.ConfigMapList).Items {
				if !s.Exists(cm.GetName()) {
					s.Store.CommonBuilder.DesiredState = &cm
					_, err := s.Store.Delete(s.Context.Context, s.Recorder)
					if err != nil {
						return err
					}
				}
			}
		case string(svc):
			s.Store.CommonBuilder.ObjectList = &corev1.ServiceList{}
			list, err := s.Store.List(s.Context.Context, s.Recorder)
			if err != nil {
				return err
			}
			for _, svc := range list.(*corev1.ServiceList).Items {
				if !s.Exists(svc.GetName()) {
					s.Store.CommonBuilder.DesiredState = &svc
					_, err := s.Store.Delete(s.Context.Context, s.Recorder)
					if err != nil {
						return err
					}
				}
			}
		}

	}

	return nil
}

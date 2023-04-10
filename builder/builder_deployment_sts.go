package builder

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type BuilderDeploymentStatefulSet struct {
	Replicas int32
	Labels   map[string]string
	PodSpec  *v1.PodSpec
	Kind     string
	CommonBuilder
}

func ToNewBuilderDeploymentStatefulSet(builder []BuilderDeploymentStatefulSet) func(*Builder) {
	return func(s *Builder) {
		s.DeploymentOrStatefulset = builder
	}
}

func (s *Builder) ReconcileDeployOrSts(cmHashes []HashHolder) (controllerutil.OperationResult, error) {

	for _, deployorsts := range s.DeploymentOrStatefulset {
		if deployorsts.Kind == "Deployment" {
			_, err := s.buildDeployment(cmHashes)
			if err != nil {
				return controllerutil.OperationResultNone, err
			}
		} else if deployorsts.Kind == "Statefulset" {
			_, err := s.buildStatefulset(cmHashes)
			if err != nil {
				return controllerutil.OperationResultNone, err
			}
		}
	}
	return controllerutil.OperationResultNone, nil
}

func (b BuilderDeploymentStatefulSet) makeDeployment(cmhashes []HashHolder) (*appsv1.Deployment, error) {

	var podSpec v1.PodSpec

	podSpec = *b.PodSpec

	if cmhashes != nil {
		for i := 0; i == len(podSpec.Containers); i++ {
			for _, cmhash := range cmhashes {
				podSpec.Containers[i].Env = append(podSpec.Containers[i].Env, v1.EnvVar{Name: cmhash.Name, Value: cmhash.HashVaule})
			}
		}
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: b.ObjectMeta,
		Spec: appsv1.DeploymentSpec{
			Replicas: &b.Replicas,

			Selector: &metav1.LabelSelector{
				MatchLabels: b.Labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: b.Labels,
				},
				Spec: podSpec,
			},
		},
	}, nil
}

func (b BuilderDeploymentStatefulSet) MakeStatefulSet(cmhashes []HashHolder) (*appsv1.StatefulSet, error) {
	var podSpec v1.PodSpec

	podSpec = *b.PodSpec

	if cmhashes != nil {
		for i := 0; i == len(podSpec.Containers); i++ {
			for _, cmhash := range cmhashes {
				podSpec.Containers[i].Env = append(podSpec.Containers[i].Env, v1.EnvVar{Name: cmhash.Name, Value: cmhash.HashVaule})
			}
		}
	}

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Statefulset",
		},
		ObjectMeta: b.ObjectMeta,
		Spec: appsv1.StatefulSetSpec{
			Replicas: &b.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: b.Labels,
			},
			Template: v1.PodTemplateSpec{
				Spec: *b.PodSpec,
				ObjectMeta: metav1.ObjectMeta{
					Labels: b.Labels,
				},
			},
		},
	}, nil
}

func (s *Builder) buildDeployment(cmhashes []HashHolder) (controllerutil.OperationResult, error) {

	for _, deploy := range s.DeploymentOrStatefulset {
		deployment, err := deploy.makeDeployment(cmhashes)
		if err != nil {
			return controllerutil.OperationResultNone, err
		}

		s.Put(deployment.GetName(), deployment.Kind)

		deploy.DesiredState = deployment
		deploy.CurrentState = &appsv1.Deployment{}

		_, err = deploy.CreateOrUpdate(s.Context.Context, s.Recorder)
		if err != nil {
			return controllerutil.OperationResultNone, err
		}
	}
	return controllerutil.OperationResultNone, nil
}

func (s *Builder) buildStatefulset(cmhashes []HashHolder) (controllerutil.OperationResult, error) {

	for _, statefulset := range s.DeploymentOrStatefulset {
		sts, err := statefulset.MakeStatefulSet(cmhashes)
		if err != nil {
			return controllerutil.OperationResultNone, err
		}

		statefulset.DesiredState = sts
		statefulset.CurrentState = &appsv1.StatefulSet{}

		_, err = statefulset.CreateOrUpdate(s.Context.Context, s.Recorder)
		if err != nil {
			return controllerutil.OperationResultNone, err
		}
	}
	return controllerutil.OperationResultNone, nil
}

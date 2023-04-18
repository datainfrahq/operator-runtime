package builder

import (
	"context"
	"errors"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type BuilderDeploymentStatefulSet struct {
	Replicas            int32
	Labels              map[string]string
	VolumeClaimTemplate []BuilderStorageConfig
	ServiceName         string
	PodSpec             *v1.PodSpec
	Kind                string
	CommonBuilder
}

func ToNewBuilderDeploymentStatefulSet(builder []BuilderDeploymentStatefulSet) func(*Builder) {
	return func(s *Builder) {
		s.DeploymentOrStatefulset = builder
	}
}

func (s *Builder) ReconcileDeployOrSts() (controllerutil.OperationResult, error) {

	for _, deployorsts := range s.DeploymentOrStatefulset {

		if deployorsts.Kind == "Deployment" {
			result, err := s.buildDeployment(deployorsts)
			if err != nil {
				return controllerutil.OperationResultNone, err
			}
			if result == controllerutil.OperationResultUpdated {
				return controllerutil.OperationResultNone, nil
			}

			if deployorsts.CrObject.GetGeneration() > 1 {
				deployorsts.CurrentState = &appsv1.Deployment{}
				done, _ := deployorsts.isObjFullyDeployed(s.Context.Context, s.Recorder)
				if !done {
					break
				}
			}
		} else if deployorsts.Kind == "Statefulset" {
			deployorsts.CurrentState = &appsv1.StatefulSet{}

			result, err := s.buildStatefulset(deployorsts)
			if err != nil {
				return controllerutil.OperationResultNone, err
			}
			if result == controllerutil.OperationResultUpdated {
				return controllerutil.OperationResultNone, nil
			}

			if deployorsts.CrObject.GetGeneration() > 1 {
				done, _ := deployorsts.isObjFullyDeployed(s.Context.Context, s.Recorder)
				if !done {
					break
				}
			}
		}
	}
	return controllerutil.OperationResultNone, nil
}

func (b *CommonBuilder) isObjFullyDeployed(ctx context.Context, recorder BuilderRecorder) (bool, error) {

	// Get Object
	obj, err := b.Get(ctx, recorder)
	if err != nil {
		return false, err
	}

	if detectType(obj) == "*v1.StatefulSet" {
		if obj.(*appsv1.StatefulSet).Status.CurrentRevision != obj.(*appsv1.StatefulSet).Status.UpdateRevision {
			return false, nil
		} else if obj.(*appsv1.StatefulSet).Status.CurrentReplicas != obj.(*appsv1.StatefulSet).Status.ReadyReplicas {
			return false, nil
		} else {
			return obj.(*appsv1.StatefulSet).Status.CurrentRevision == obj.(*appsv1.StatefulSet).Status.UpdateRevision, nil
		}
	} else if detectType(obj) == "*v1.Deployment" {
		for _, condition := range obj.(*appsv1.Deployment).Status.Conditions {
			// This detects a failure condition, operator should send a rolling deployment failed event
			if condition.Type == appsv1.DeploymentReplicaFailure {
				return false, errors.New(condition.Reason)
			} else if condition.Type == appsv1.DeploymentProgressing && condition.Status != v1.ConditionTrue || obj.(*appsv1.Deployment).Status.ReadyReplicas != obj.(*appsv1.Deployment).Status.Replicas {
				return false, nil
			} else {
				return obj.(*appsv1.Deployment).Status.ReadyReplicas == obj.(*appsv1.Deployment).Status.Replicas, nil
			}
		}
	}
	return false, nil
}

func (b *BuilderDeploymentStatefulSet) makeDeployment() (*appsv1.Deployment, error) {

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: b.ObjectMeta,
		Spec: appsv1.DeploymentSpec{
			Replicas: &b.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"custom_resource": b.CrObject.GetName(),
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: b.Labels,
				},
				Spec: *b.PodSpec,
			},
		},
	}, nil
}

func (b *BuilderDeploymentStatefulSet) MakeStatefulSet() (*appsv1.StatefulSet, error) {

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
			ServiceName: b.ServiceName,
			Template: v1.PodTemplateSpec{
				Spec: *b.PodSpec,
				ObjectMeta: metav1.ObjectMeta{
					Labels: b.Labels,
				},
			},
		},
	}, nil
}

func (s *Builder) buildDeployment(deploy BuilderDeploymentStatefulSet) (controllerutil.OperationResult, error) {

	deployment, err := deploy.makeDeployment()
	if err != nil {
		return controllerutil.OperationResultNone, err
	}

	s.Put(deployment.GetName(), deployment.Kind)

	deploy.DesiredState = deployment
	deploy.CurrentState = &appsv1.Deployment{}

	result, err := deploy.CreateOrUpdate(s.Context.Context, s.Recorder)
	if err != nil {
		return controllerutil.OperationResultNone, err
	}

	return result, nil
}

func (s *Builder) buildStatefulset(statefulset BuilderDeploymentStatefulSet) (controllerutil.OperationResult, error) {

	sts, err := statefulset.MakeStatefulSet()
	if err != nil {
		return controllerutil.OperationResultNone, err
	}

	sts.Spec.VolumeClaimTemplates = statefulset.MakeVolumeClaimTemplates()

	statefulset.DesiredState = sts
	statefulset.CurrentState = &appsv1.StatefulSet{}

	_, err = statefulset.CreateOrUpdate(s.Context.Context, s.Recorder)
	if err != nil {
		return controllerutil.OperationResultNone, err
	}

	return controllerutil.OperationResultNone, nil
}

func (b *BuilderDeploymentStatefulSet) MakeVolumeClaimTemplates() []v1.PersistentVolumeClaim {

	var pvcs []v1.PersistentVolumeClaim
	for _, storage := range b.VolumeClaimTemplate {

		pvc, err := storage.MakePvc()
		if err != nil {
			return []v1.PersistentVolumeClaim{}
		}

		pvcs = append(pvcs, *pvc)

	}

	return pvcs
}

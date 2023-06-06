package builder

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type BuilderNetworkPolicy struct {
	NetworkPolicySpec *networkingv1.NetworkPolicySpec
	CommonBuilder
}

func ToNewBuilderNetworkPolicy(builder []BuilderNetworkPolicy) func(*Builder) {
	return func(s *Builder) {
		s.NetworkPolicy = builder
	}
}

func (s *Builder) ReconcileNetworkPolicy() (controllerutil.OperationResult, error) {

	var err error
	var result controllerutil.OperationResult

	for _, np := range s.NetworkPolicy {

		if np.NetworkPolicySpec != nil {

			makeNp := np.makeNetworkPolicy()

			np.DesiredState = makeNp
			np.CurrentState = &networkingv1.NetworkPolicy{}

			result, err = np.CreateOrUpdate(s.Context.Context, s.Recorder)
			if err != nil {
				return controllerutil.OperationResultNone, nil
			}
		}
	}
	return result, nil
}

func (b *BuilderNetworkPolicy) makeNetworkPolicy() *networkingv1.NetworkPolicy {
	networkPolicy := &networkingv1.NetworkPolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1",
			Kind:       "NetworkPolicy",
		},
		ObjectMeta: b.ObjectMeta,
		Spec:       *b.NetworkPolicySpec,
	}

	return networkPolicy
}

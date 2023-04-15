package builder

import (
	"context"

	"github.com/datainfrahq/operator-runtime/utils"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (b *CommonBuilder) CreateOrUpdate(ctx context.Context, buildRecorder BuilderRecorder) (controllerutil.OperationResult, error) {
	addOwnerRefToObject(b.DesiredState, b.OwnerRef)
	utils.AddHashToObject(b.DesiredState, b.OwnerRef.Kind+"OperatorHash")
	if err := b.Client.Get(ctx, types.NamespacedName{Name: b.DesiredState.GetName(), Namespace: b.DesiredState.GetNamespace()}, b.CurrentState); err != nil {
		if apierrors.IsNotFound(err) {
			result, err := b.Create(ctx, buildRecorder)
			if err != nil {
				return controllerutil.OperationResultNone, err
			}
			return result, nil
		} else {
			return "", err
		}
	} else {
		if b.DesiredState.GetAnnotations()[b.OwnerRef.Kind+"OperatorHash"] != b.CurrentState.GetAnnotations()[b.OwnerRef.Kind+"OperatorHash"] {
			b.DesiredState.SetResourceVersion(b.CurrentState.GetResourceVersion())
			result, err := b.Update(ctx, buildRecorder)
			if err != nil {
				return controllerutil.OperationResultNone, err
			} else {
				return result, nil
			}
		} else {
			return controllerutil.OperationResultNone, nil
		}
	}
}

func addOwnerRefToObject(obj metav1.Object, ownerRef metav1.OwnerReference) {
	trueVar := true
	ownerRef = metav1.OwnerReference{
		APIVersion: ownerRef.APIVersion,
		Kind:       ownerRef.Kind,
		Name:       ownerRef.Name,
		UID:        ownerRef.UID,
		Controller: &trueVar,
	}
	obj.SetOwnerReferences(append(obj.GetOwnerReferences(), ownerRef))
}

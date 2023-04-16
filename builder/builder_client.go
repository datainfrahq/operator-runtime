package builder

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (b *CommonBuilder) Create(ctx context.Context, buildRecorder BuilderRecorder) (controllerutil.OperationResult, error) {

	if err := b.Client.Create(ctx, b.DesiredState); err != nil {
		buildRecorder.createEvent(b.CrObject, b.DesiredState, err)
		return controllerutil.OperationResultNone, err
	} else {
		buildRecorder.createEvent(b.CrObject, b.DesiredState, nil)
		return controllerutil.OperationResultCreated, nil
	}
}

func (b *CommonBuilder) Update(ctx context.Context, buildRecorder BuilderRecorder) (controllerutil.OperationResult, error) {
	if err := b.Client.Update(ctx, b.DesiredState); err != nil {
		buildRecorder.updateEvent(b.CrObject, b.DesiredState, err)
		return controllerutil.OperationResultNone, err
	} else {
		buildRecorder.updateEvent(b.CrObject, b.DesiredState, nil)
		return controllerutil.OperationResultUpdated, nil
	}
}

func (b *CommonBuilder) Get(ctx context.Context, buildRecorder BuilderRecorder) (client.Object, error) {
	if err := b.Client.Get(ctx, *namespacedName(b.DesiredState.GetName(), b.ObjectMeta.Namespace), b.CurrentState); err != nil {
		buildRecorder.getEvent(b.CrObject, b.DesiredState, err)
		return nil, err
	} else {
		return b.CurrentState, nil
	}
}

func (b *CommonBuilder) List(ctx context.Context, buildRecorder BuilderRecorder) (client.ObjectList, error) {
	listOpts := []client.ListOption{
		client.InNamespace(b.ObjectMeta.Namespace),
		client.MatchingLabels(b.Labels),
	}

	deployment := b.ObjectList
	if err := b.Client.List(ctx, deployment, listOpts...); err != nil {
		return nil, err
	} else {
		return deployment, nil
	}
}

func (b *CommonBuilder) Delete(ctx context.Context, buildRecorder BuilderRecorder) (controllerutil.OperationResult, error) {
	if err := b.Client.Delete(ctx, b.DesiredState); err != nil {
		buildRecorder.deleteEvent(b.CrObject, b.DesiredState, err)
		return controllerutil.OperationResultNone, err
	} else {
		buildRecorder.deleteEvent(b.CrObject, b.DesiredState, nil)
		return controllerutil.OperationResultUpdated, nil
	}
}

func namespacedName(name, namespace string) *types.NamespacedName {
	return &types.NamespacedName{Name: name, Namespace: namespace}
}

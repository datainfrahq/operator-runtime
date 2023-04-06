package builder

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (b CommonBuilder) Create(ctx context.Context, buildRecorder BuilderRecorder) (controllerutil.OperationResult, error) {

	if err := b.Client.Create(ctx, b.DesiredState); err != nil {
		buildRecorder.createEvent(b.CrObject, b.DesiredState, err)
		return "", err
	} else {
		buildRecorder.createEvent(b.CrObject, b.DesiredState, nil)
		return controllerutil.OperationResultCreated, nil
	}
}

func (b CommonBuilder) Update(ctx context.Context, buildRecorder BuilderRecorder) (controllerutil.OperationResult, error) {
	if err := b.Client.Update(ctx, b.DesiredState); err != nil {
		buildRecorder.updateEvent(b.CrObject, b.DesiredState, err)
		return "", err
	} else {
		buildRecorder.updateEvent(b.CrObject, b.DesiredState, nil)
		return controllerutil.OperationResultUpdated, nil
	}
}

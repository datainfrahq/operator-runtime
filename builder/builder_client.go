package builder

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (b CommonBuilder) Create(ctx context.Context) (controllerutil.OperationResult, error) {
	if err := b.Client.Create(ctx, b.DesiredState); err != nil {
		return "", err
	} else {
		return controllerutil.OperationResultCreated, nil
	}
}

func (b CommonBuilder) Update(ctx context.Context) (controllerutil.OperationResult, error) {
	if err := b.Client.Update(ctx, b.DesiredState); err != nil {
		return "", err
	} else {
		return controllerutil.OperationResultUpdated, nil
	}
}

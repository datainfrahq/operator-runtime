package client

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type Builder struct {
	client       client.Client
	DesiredState client.Object
	CurrentState client.Object
	ObjectName   string
	CrObject     interface{}
}

func (b Builder) CreateOrUpdate() (controllerutil.OperationResult, error) {
	if err := b.client.Get(context.TODO(), types.NamespacedName{Name: b.CurrentState.GetName(), Namespace: b.CurrentState.GetNamespace()}, b.DesiredState); err != nil {
		if apierrors.IsNotFound(err) {
			// resource does not exist, create it.
			fmt.Println("Create")
		} else {
			fmt.Println("Delete")
			return "", err
		}
	} else {

		fmt.Println("Update")
	}
	return controllerutil.OperationResultNone, nil
}

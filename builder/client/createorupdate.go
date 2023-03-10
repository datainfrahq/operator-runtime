package client

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type BuilderState struct {
	Client       client.Client
	DesiredState client.Object
	CurrentState client.Object
	CrObject     client.Object
	OwnerRef     metav1.OwnerReference
}

func (b BuilderState) CreateOrUpdate() (controllerutil.OperationResult, error) {
	addOwnerRefToObject(b.DesiredState, b.OwnerRef)
	addHashToObject(b.DesiredState, b.OwnerRef.Kind+"OperatorHash")
	if err := b.Client.Get(context.TODO(), types.NamespacedName{Name: b.DesiredState.GetName(), Namespace: b.DesiredState.GetNamespace()}, b.CurrentState); err != nil {
		if apierrors.IsNotFound(err) {
			result, err := b.Create(context.TODO())
			if err != nil {
				return controllerutil.OperationResultNone, err
			}
			return result, nil
		} else {
			fmt.Println("Delete")
			return "", err
		}
	} else {
		if b.DesiredState.GetAnnotations()[b.OwnerRef.Kind+"OperatorHash"] != b.CurrentState.GetAnnotations()[b.OwnerRef.Kind+"OperatorHash"] {
			b.DesiredState.SetResourceVersion(b.CurrentState.GetResourceVersion())
			result, err := b.Update(context.TODO())
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

func (b BuilderState) Create(ctx context.Context) (controllerutil.OperationResult, error) {
	if err := b.Client.Create(ctx, b.DesiredState); err != nil {
		return "", err
	} else {
		return controllerutil.OperationResultCreated, nil
	}
}

func (b BuilderState) Update(ctx context.Context) (controllerutil.OperationResult, error) {
	if err := b.Client.Update(ctx, b.DesiredState); err != nil {
		return "", err
	} else {
		return controllerutil.OperationResultUpdated, nil
	}
}

type BuilderObject struct {
	CrObject   client.Object
	ObjectMeta metav1.ObjectMeta
	Data       map[string]string
}

func NewBuilderObject(
	crObject client.Object,
	objectMeta metav1.ObjectMeta,
	data map[string]string,
) *BuilderObject {
	return &BuilderObject{
		CrObject:   crObject,
		ObjectMeta: objectMeta,
		Data:       data,
	}
}

func (b BuilderObject) MakeConfigMap() (*v1.ConfigMap, error) {
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: b.ObjectMeta,
		Data:       b.Data,
	}, nil
}

// addOwnerRefToObject appends the desired OwnerReference to the object
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

func addHashToObject(obj client.Object, name string) error {
	if sha, err := getObjectHash(obj); err != nil {
		return err
	} else {
		annotations := obj.GetAnnotations()
		if annotations == nil {
			annotations = make(map[string]string)
			obj.SetAnnotations(annotations)
		}
		annotations[name] = sha
		return nil
	}
}

func getObjectHash(obj client.Object) (string, error) {
	if bytes, err := json.Marshal(obj); err != nil {
		return "", err
	} else {
		sha1Bytes := sha1.Sum(bytes)
		return base64.StdEncoding.EncodeToString(sha1Bytes[:]), nil
	}
}

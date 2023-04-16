package builder

import (
	"fmt"
	"reflect"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type BuilderRecorder struct {
	Recorder       record.EventRecorder
	ControllerName string
}

func ToNewBuilderRecorder(builder BuilderRecorder) func(*Builder) {
	return func(s *Builder) {
		s.Recorder = builder
	}
}

func (b *BuilderRecorder) createEvent(crObj client.Object, obj client.Object, err error) {
	if err != nil {
		b.Recorder.Event(
			crObj,
			v1.EventTypeWarning,
			fmt.Sprintf("Name [%s], Namespace [%s], Kind [%s]", obj.GetName(), obj.GetNamespace(), detectType(obj)),
			b.ControllerName+"CreateObjectFail")
	} else {
		b.Recorder.Event(
			crObj,
			v1.EventTypeNormal,
			fmt.Sprintf("Name [%s], Namespace [%s], Kind [%s]", obj.GetName(), obj.GetNamespace(), detectType(obj)),
			b.ControllerName+"CreateObjectSuccess")
	}
}

func (b *BuilderRecorder) updateEvent(crObj client.Object, obj client.Object, err error) {
	if err != nil {
		b.Recorder.Event(
			crObj,
			v1.EventTypeWarning,
			fmt.Sprintf("Name [%s], Namespace [%s], Kind [%s]", obj.GetName(), obj.GetNamespace(), detectType(obj)),
			b.ControllerName+"UpdateObjectFail")
	} else {
		b.Recorder.Event(
			crObj,
			v1.EventTypeNormal,
			fmt.Sprintf("Name [%s], Namespace [%s], Kind [%s]", obj.GetName(), obj.GetNamespace(), detectType(obj)),
			b.ControllerName+"UpdateObjectSuccess")
	}
}

func (b *BuilderRecorder) getEvent(crObj client.Object, obj client.Object, err error) {
	if err != nil {
		b.Recorder.Event(
			crObj,
			v1.EventTypeWarning,
			fmt.Sprintf("Name [%s], Namespace [%s], Kind [%s]", obj.GetName(), obj.GetNamespace(), detectType(obj)),
			b.ControllerName+"GetObjectFail")
	}
}

func (b *BuilderRecorder) listEvent(crObj client.Object, obj client.Object, err error) {
	if err != nil {
		b.Recorder.Event(
			crObj,
			v1.EventTypeWarning,
			fmt.Sprintf("Name [%s], Namespace [%s], Kind [%s]", obj.GetName(), obj.GetNamespace(), detectType(obj)),
			b.ControllerName+"ListObjectFail")
	}
}

func (b *BuilderRecorder) deleteEvent(crObj client.Object, obj client.Object, err error) {
	if err != nil {
		b.Recorder.Event(
			crObj,
			v1.EventTypeWarning,
			fmt.Sprintf("Name [%s], Namespace [%s], Kind [%s]", obj.GetName(), obj.GetNamespace(), detectType(obj)),
			b.ControllerName+"DeleteObjectFail")
	} else {
		b.Recorder.Event(
			crObj,
			v1.EventTypeNormal,
			fmt.Sprintf("Name [%s], Namespace [%s], Kind [%s]", obj.GetName(), obj.GetNamespace(), detectType(obj)),
			b.ControllerName+"DeleteObjectSuccess")
	}
}

func detectType(obj client.Object) string { return reflect.TypeOf(obj).String() }

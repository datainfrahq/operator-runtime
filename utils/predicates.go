package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type Predicates interface {
	IgnoreNamespacePredicate(obj client.Object) bool
	IgnoreObjectPredicate(obj client.Object) bool
	IgnoreUpdate(e event.UpdateEvent) bool
}

// CommonPredicates struct holds the fields to initalise common predicate methods.
type CommonPredicates struct {
	ControllerName   string
	IgnoreAnnotation string
	Log              logr.Logger
}

// NewCommonPredicates is contstructor for CommonPredicates
func NewCommonPredicates(
	controllerName, ignoreAnnotation string,
	log logr.Logger) Predicates {
	return &CommonPredicates{
		ControllerName:   controllerName,
		IgnoreAnnotation: ignoreAnnotation,
		Log:              log,
	}
}

// IgnoreNamespacePredicate is a function that, when initialized within the create and update predicates,
// filters out the namespaces that should NOT be reconciled in the DENY_LIST. This function is particularly
// useful in scenarios where the controller watches all the namespaces but wants to exclude certain ones, such
// as kube-system and default. By using the DENY_LIST feature, the controller can ensure that the excluded
// namespaces are not processed during the reconciliation process, thus reducing unnecessary overhead.
//
// Controllers can watch a single namespace, multiple namespaces, or all namespaces. When the IgnoreNamespacePredicate
// function is used, it allows the controller to further refine which namespaces are processed. This function should
// be called within the create and update predicates to ensure that the namespaces are properly filtered before
// reconciliation takes place.
func (c *CommonPredicates) IgnoreNamespacePredicate(obj client.Object) bool {
	var log = c.Log.WithName("predicates")

	namespaces := getEnvAsSlice("DENY_LIST", nil, ",")

	for _, namespace := range namespaces {
		if obj.GetNamespace() == namespace {
			msg := fmt.Sprintf("%s will not reconcile namespace [%s], alter DENY_LIST to reconcile", c.ControllerName, obj.GetNamespace())
			log.Info(msg)
			return false
		}
	}
	return true
}

// IgnoreIgnoredObjectPredicate is a function that, when initialized within the create or update predicate,
// filters out the namespaces that have an ignore annotation present. This function is particularly useful
// when a controller is reconciling multiple, all, or a single custom resource. In such scenarios, the controller
// may need to pause reconciliation on a specific custom resource. By using this function, the controller can
// filter out the namespaces with the ignore annotation.
//
// This function should be called within the create or update predicate to ensure that the namespaces with the
// ignore annotation are properly filtered before reconciliation takes place.
func (c *CommonPredicates) IgnoreObjectPredicate(obj client.Object) bool {
	var log = c.Log.WithName("predicates")

	if ignoredStatus := obj.GetAnnotations()[c.IgnoreAnnotation]; ignoredStatus == "true" {
		msg := fmt.Sprintf("%s will not re-concile ignored Druid [%s], removed annotation to re-concile", c.ControllerName, obj.GetName())
		log.Info(msg)
		return false
	}
	return true
}

// IgnoreUpdate to ignore the Update events when nothing has changed in the resource version
func (c *CommonPredicates) IgnoreUpdate(e event.UpdateEvent) bool {

	if e.ObjectOld == nil {
		return false
	}

	if e.ObjectNew == nil {
		return false
	}

	if e.ObjectNew.GetGeneration() == e.ObjectOld.GetGeneration() && e.ObjectNew.GetGeneration() != 0 {
		return false
	}
	return true
}

// lookup if denylist env exists
func getDenyListEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// getEnvAsSlice gets the DENYLIST string and returns a slice of strings seperated by ","
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getDenyListEnv(name, "")
	if valStr == "" {
		return defaultVal
	}
	// split on ","
	val := strings.Split(valStr, sep)
	return val
}

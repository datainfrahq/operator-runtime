package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"os"
	"time"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ReconcileWait string = "RECONCILE_WAIT"
)

type ConfigMapHash struct {
	Object client.Object
}

type HashHolder struct {
	Name      string
	HashVaule string
}

func MakeConfigMapHash(configHash []ConfigMapHash) ([]HashHolder, error) {
	hashHolder := []HashHolder{}

	for _, obj := range configHash {
		valueSHA, err := getObjectHash(obj.Object)
		if err != nil {
			return nil, err
		}
		hashHolder = append(hashHolder, HashHolder{
			Name:      obj.Object.GetName(),
			HashVaule: valueSHA,
		},
		)
	}

	return unique(hashHolder), nil
}

func AddHashToObject(obj client.Object, name string) error {
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

func unique[T comparable](s []T) []T {
	inResult := make(map[T]bool)
	var result []T
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

func LookupReconcileTime(log logr.Logger) time.Duration {
	val, exists := os.LookupEnv(ReconcileWait)
	if !exists {
		return time.Second * 10
	} else {
		v, err := time.ParseDuration(val)
		if err != nil {
			log.Error(err, err.Error())
			// Exit Program if not valid
			os.Exit(1)
		}
		return v
	}
}

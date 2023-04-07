package builder

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type BuilderConfigMapHash struct {
	Object client.Object
}

type HashHolder struct {
	Name      string
	HashVaule string
}

func ToNewBuilderConfigMapHash(builder []BuilderConfigMapHash) func(*Builder) {
	return func(s *Builder) {
		s.ConfigHash = builder
	}
}

func (s *Builder) ReconcileConfigMapHash() ([]HashHolder, error) {
	hashHolder := []HashHolder{}

	for _, obj := range s.ConfigHash {
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
	return hashHolder, nil
}

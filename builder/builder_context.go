package builder

import "context"

type BuilderContext struct {
	Context context.Context
}

func ToNewBuilderContext(builder BuilderContext) func(*Builder) {
	return func(s *Builder) {
		s.Context = builder
	}
}

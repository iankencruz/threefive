package project

import "github.com/iankencruz/threefive/internal/generated"

type ProjectWithMedia struct {
	*generated.Project
	Media []*generated.Media `json:"media"`
}

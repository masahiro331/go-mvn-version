package version

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type Comparer interface {
	Check(v Version) bool
}

func NewComparer(v string) (Comparer, error) {
	var errs error

	c, err := NewConstraints(v)
	if err == nil {
		return c, nil
	}
	errs = multierror.Append(errs, err)

	r, err := NewRequirements(v)
	if err == nil {
		return r, nil
	}
	errs = multierror.Append(errs, err)

	return nil, fmt.Errorf("failed to new comparer: %w", errs)
}

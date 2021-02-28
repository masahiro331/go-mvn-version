package version

import "golang.org/x/xerrors"

type Comparer interface {
	Check(v Version) bool
}

func NewComparer(v string) (Comparer, error) {
	c, err := NewConstraints(v)
	if err == nil {
		return c, nil
	}
	r, err := NewRequirements(v)
	if err == nil {
		return r, nil
	}
	return nil, xerrors.Errorf("failed to new comparer: %w", err)
}

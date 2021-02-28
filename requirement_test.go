package version

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRequirements(t *testing.T) {
	tests := []struct {
		requirement string
		wantErr     bool
	}{
		// if soft requirement for 1.0. Use constraints compare
		{"1.0", false},

		{"(, 1.0)", false},
		{"(,1.0)", false},
		{"(, 1.0]", false},
		{"(,1.0]", false},
		{"[1.0,)", false},
		{"[1.0, )", false},
		{"[1.0,]", false},
		{"[1.0, ]", false},
		{"(0.9, 1.0)", false},
		{"(0.9, 1.0]", false},
		{"[1.0, 1.1)", false},
		{"[1.0, 1.1]", false},
		{"(0.9,1.0)", false},
		{"(0.9,1.0]", false},
		{"[1.0,1.1)", false},
		{"[1.0,1.1]", false},

		{"[,]", false},
		{"[0,]", false},
		{"[,0]", false},
		{"[0,)", false},
		{"(,0)", false},
		{"(0,)", false},
		{"[0,)", false},
		{"(,0]", false},

		{"[2.4.0,2.4.2],[2.4.4]", false},
		{"[2.4.0,2.4.2],[2.4.4],[2.5.5]", false},

		{"1.0)", true},
		{"1.0]", true},
		{"(1.0", true},
		{"[1.0", true},
		{", 1.0)", true},
		{", 1.0]", true},
		{"(1.0, ", true},
		{"[1.0, ", true},
		{"<1.0", true},

		{"(0.9,1.0,1.2)", true},
		{"(0.9,1.0,1.2]", true},
		{"[1.0,1.1,1.2)", true},
		{"[1.0,1.1,1.2]", true},
	}

	for _, tt := range tests {
		t.Run(tt.requirement, func(t *testing.T) {
			_, err := NewRequirements(tt.requirement)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRequirementsCheck(t *testing.T) {
	tests := []struct {
		requirement string
		version     string
		want        bool
	}{
		{"[,1.0.0]", "0.9", true},
		{"[,1.0.0]", "1.0", true},
		{"[,1.0.0]", "1.1", false},

		{"(,1.0.0]", "0.9", true},
		{"(,1.0.0]", "1.0", true},
		{"(,1.0.0]", "1.1", false},

		{"(,1.0.0)", "0.9", true},
		{"(,1.0.0)", "1.0.0", false},
		{"(,1.0.0)", "1.0.1", false},

		{"[,1.0.0)", "0.9", true},
		{"[,1.0.0)", "1.0.0", false},
		{"[,1.0.0)", "1.0.1", false},

		{"[0,)", "0.9", true},
		{"[0,)", "1.0.0", true},
		{"[0,)", "1.0.1", true},

		{"(,0)", "0.9", false},
		{"(,0)", "1.0.0", false},
		{"(,0)", "1.0.1", false},

		{"[,]", "0.9", true},
		{"[,]", "1.0.0", true},
		{"[,]", "1.0.1", true},

		{"[1.0.0,)", "1.0.0", true},
		{"[1.0.0,)", "1.0.1", true},
		{"[1.0.0,)", "0.9", false},

		{"[1.0.0]", "1.0.0", true},
		{"[1.0.0]", "1.0.1", false},
		{"[1.0.0]", "0.9", false},

		// (1.0.0) is not defined requirements. use [1.0.0] compare.
		{"(1.0.0)", "1.0.0", true},
		{"(1.0.0)", "1.0.1", false},
		{"(1.0.0)", "0.9", false},

		{"[1.0,2.0)", "1.0.0", true},
		{"[1.0,2.0)", "1.5", true},
		{"[1.0,2.0)", "0.9", false},
		{"[1.0,2.0)", "2.0", false},

		{"[1.0,2.0]", "1.0.0", true},
		{"[1.0,2.0]", "1.5", true},
		{"[1.0,2.0]", "2.0", true},
		{"[1.0,2.0]", "0.9", false},
		{"[1.0,2.0]", "2.1", false},

		{"(1.0,2.0]", "1.0.0", false},
		{"(1.0,2.0]", "1.5", true},
		{"(1.0,2.0]", "2.0", true},
		{"(1.0,2.0]", "0.9", false},
		{"(1.0,2.0]", "2.1", false},
		{"(1.0,2.0]", "2.1", false},

		{"(,1.0.5.RELEASE],[2.0.0.RELEASE,2.0.16.RELEASE),[2.1.0.RELEASE,2.1.3.RELEASE)", "1.0.0", true},
		{"(,1.0.5.RELEASE],[2.0.0.RELEASE,2.0.16.RELEASE),[2.1.0.RELEASE,2.1.3.RELEASE)", "2.0.0", true},
		{"(,1.0.5.RELEASE],[2.0.0.RELEASE,2.0.16.RELEASE),[2.1.0.RELEASE,2.1.3.RELEASE)", "2.1.3", false},

		// soft requirement
		{"1.0", "2.0", true},
		{"1.0", "1.0", true},
		{"1.0", "0.1", true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s %s", tt.version, tt.requirement), func(t *testing.T) {
			r, err := NewRequirements(tt.requirement)
			require.NoError(t, err)

			v, err := NewVersion(tt.version)
			require.NoError(t, err)

			assert.Equal(t, tt.want, r.Check(v))
		})
	}
}

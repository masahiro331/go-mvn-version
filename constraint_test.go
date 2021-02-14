package version

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConstraints(t *testing.T) {
	tests := []struct {
		constraint string
		wantErr    bool
	}{
		{"> 1.0", false},
		{"= abc", false},
		{"> 1.0 || < foo", false},
		{">= 1.2.3, < 2.0 || => 3.0, < 4", false},
		{">= 1.0.1.v100000", false},
		{">= 1.1", false},
		{">40.50.60, < 50.70", false},
		{"2.0", false},
		{"2.3.5-20161202202307-sha.e8fc5e5", false},
		{">= bar", false},
		{">= 1.1.1.v1", false},
		{">= 1.1.1.1v", false},
		{"==1.1.1.1v", false},
		{"BAR >= 1.2.3", false},
		{"!= !=", true},
		{"bar <", true},
	}
	for _, tt := range tests {
		t.Run(tt.constraint, func(t *testing.T) {
			_, err := NewConstraints(tt.constraint)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVersion_Check(t *testing.T) {
	tests := []struct {
		constraint string
		version    string
		want       bool
	}{
		// Equal
		{"=2.0.0", "1.2.3", false},
		{"=2.0.0", "2.0.0", true},
		{"= 1.0", "1.1.5", false},
		{"= 1.0", "1.0.0", true},
		{"=4.1", "4.1.0-alpha.1", false},
		{"=4.1-alpha", "4.1.0-alpha", true},
		{"=2.0", "1.2.3", false},
		{"=2.0", "2.0.0", true},
		{"=2.0", "2.0.1", false},
		{"=0", "1.0.0", false},

		{"== 2.0.0", "1.2.3", false},
		{"==2.0.0", "2.0.0", true},
		{"== 4.1", "4.1.0-alpha.1", false},
		{"==4.1-alpha", "4.1.0-alpha", true},

		{"2", "1.0.0", false},
		{"2", "3.4.5", false},
		{"2", "2.1.1", false},
		{"2.1", "2.1.1", false},
		{"2.1", "2.2.1", false},
		{"4.1", "4.1.0", true},
		{"1.0", "1.0.0", true},

		// Not equal
		{"!=4.1.0", "4.1.0", false},
		{"!=4.1.0", "4.1.1", true},
		{"!=4.1", "5.1.0-alpha.1", true},
		{"!=4.1-alpha", "4.1.0", true},

		// Less than
		{"<0.0.5", "0.1.0", false},
		{"<1.0.0", "0.1.0", true},
		{"<0", "0.0.0-alpha", true},
		{"<0-z", "0.0.0-alpha", true},
		{"<0", "1.0.0-alpha", false},
		{"<1", "1.0.0-alpha", true},
		{"<11", "0.1.0", true},
		{"<11", "11.1.0", false},
		{"<1.1", "0.1.0", true},
		{"<1.1", "1.1.0", false},
		{"<1.1", "1.1.1", false},

		// Less than or equal
		{"<=0.2.3", "1.2.3", false},
		{"<=1.2.3", "1.2.3", true},
		{"<= 2.1.0-a", "2.0.0", true},
		{"<=11", "1.2.3", true},
		{"<=11", "12.2.3", false},
		{"<=11", "11.2.3", false}, // different
		{"<=1.1", "1.2.3", false},
		{"<=1.1", "0.1.0", true},
		{"<=1.1", "1.1.0", true},
		{"<=1.1", "1.1.1", false}, // different
		{"<=0-0", "0.0.0-alpha", true},
		{"<=0.0.0-0", "0.0.0-alpha", true},

		// Greater than
		{">5.0.0", "4.1.0", false},
		{">4.0.0", "4.1.0", true},
		{"> 2.0", "2.1.0-beta", true},
		{">0", "0.0.1-alpha", true},
		{">0.0", "0.0.1-alpha", true},
		{">0-0", "0.0.1-alpha", true},
		{">0.0-0", "0.0.1-alpha", true},
		{">0", "0.0.0-alpha", false},
		{">0-0", "0.0.0-alpha", false},
		{">0.0.0-0", "0.0.0-alpha", false},
		{">1.2.3-alpha.1", "1.2.3-alpha.2", true},
		{">1.2.3-alpha.1", "1.3.3-alpha.2", true},
		{">1.1", "4.1.0", true},
		{">1.1", "1.1.0", false},
		{">0", "0.0.0", false},
		{">0", "1.0.0", true},
		{">11", "11.1.0", true}, // different
		{">11.1", "11.1.0", false},
		{">11.1", "11.1.1", true}, // different
		{">11.1", "11.2.1", true},

		// Greater than or equal
		{">=11.1.3", "11.1.2", false},
		{">=11.1.2", "11.1.2", true},
		{">= 1.0, < 1.2", "1.1.5", true},
		{">= 2.1.0-alpha-1", "2.1.0-beta-1", true},
		{">= 2.1.0-a", "2.1.1-beta", true},
		{">= 2.0.0", "2.1.0-beta", true},
		{">= 2.1.0-a", "2.1.1", true},
		{">= 2.1.0-alpha", "2.1.0", true},
		{">=0", "0.0.1-alpha", true},
		{">=0.0", "0.0.1-alpha", true},
		{">=0-0", "0.0.1-alpha", true},
		{">=0.0-0", "0.0.1-alpha", true},
		{">=0", "0.0.0-alpha", false},
		{">=0.0.0-0", "1.2.3", true},
		{">=0.0.0-0", "3.4.5-beta.1", true},
		{">=11", "11.1.2", true},
		{">=11.1", "11.1.2", true},
		{">=11.1", "11.0.2", false},
		{">=1.1", "4.1.0", true},
		{">=1.1", "1.1.0", true},
		{">=1.1", "0.0.9", false},
		{">=0", "0.0.0", true},

		// More than 3 numbers
		{"< 1.0.0.1 || = 2.0.1.2.3", "2.0", false},
		{"< 1.0.0.1 || = 2.0.5.4.8", "2.0.5.4.8", true},
		{"> 1.0.0.0.1 < 1.0.0.1 || = 2.0.5.4.8", "1.0.0.0.9", true},

		// Leading zeroes
		{">1.2.3", "1.02.4", true},
		{"<1.3.09", "1.05.4", false},

		// Multiple constraints
		{"< 1.0 || = 2.0", "2.0", true},
		{"< 1.0 || = 2.0", "0.1", true},
		{"< 1.0 || = 2.0", "1.1", false},
		{"> 1.0, < 1.2", "1.1.5", true},
		{"> 1.0, < 1.2 || >3.0", "1.5", false},
		{"> 1.0 < 1.2 || >3.0", "1.5", false},
		{"> 1.0	< 1.2 || >3.0", "1.1", true},
		{"> 1.0, < 1.2 || >3.0", "4.2", true},
		{"> 1.0 < 1.2 || >3.0, <4.0", "4.2", false},

		// add more tests from ghsa data
		{"< 0.3.0M2", "0.3.0m1", true},
		{"= 0.3.0M2", "0.3.0m2", true},
		{"> 0.3.0M2", "0.3.0m3", true},
		{"0.3.0M2", "0.3.0-milestone-2", true},
		{"> 0.3.0M2", "0.3.0-milestone-3", true},
		{"< 9.2.25.v20180606", "9.2.25.v20180605", true},
		{"< 1.1.1.v2", "1.1.1.v1", true},
		{"< 1.1.v2", "1.1.v1", true},
		{"< 1.v2", "1.v1", true},
		{"< v2", "v1", true},
		{"< 1.1.1.2", "1.1.1.1", true},
		{"< 1.1.2", "1.1.1", true},
		{"< 1.2", "1.1", true},
		{"< 2", "1", true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s %s", tt.version, tt.constraint), func(t *testing.T) {
			c, err := NewConstraints(tt.constraint)
			require.NoError(t, err)

			v, err := NewVersion(tt.version)
			require.NoError(t, err)

			assert.Equal(t, tt.want, c.Check(v))
		})
	}
}

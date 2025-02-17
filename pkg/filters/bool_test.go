package filters

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoolFilterParse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		expressions  []string
		expected     bool
		filterResult []bool // filter on []bool{true, false}
	}{
		{
			name:         "eval true 1",
			expressions:  []string{"container"},
			expected:     true,
			filterResult: []bool{true, false},
		},
		{
			name:         "eval true 2",
			expressions:  []string{"=true"},
			expected:     true,
			filterResult: []bool{true, false},
		},
		{
			name:         "eval true 3",
			expressions:  []string{"!=false"},
			expected:     true,
			filterResult: []bool{true, false},
		},
		{
			name:         "eval false 1",
			expressions:  []string{"not-container"},
			expected:     false,
			filterResult: []bool{false, true},
		},
		{
			name:         "eval false 2",
			expressions:  []string{"=false"},
			expected:     false,
			filterResult: []bool{false, true},
		},
		{
			name:         "eval false 3",
			expressions:  []string{"!=true"},
			expected:     false,
			filterResult: []bool{false, true},
		},
		{
			name:         "eval false then true",
			expressions:  []string{"not-container", "=true"},
			expected:     true,
			filterResult: []bool{true, true},
		},
		{
			name:         "no values",
			expressions:  []string{},
			expected:     false,
			filterResult: []bool{false, false},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			filter := NewBoolFilter()
			for _, expr := range tc.expressions {
				err := filter.Parse(expr)
				require.NoError(t, err)
			}

			filter.Enable()

			assert.Equal(t, tc.expected, filter.Value())
			filterRes := []bool{}
			for _, val := range []bool{true, false} {
				filterRes = append(filterRes, filter.Filter(val))
			}
			assert.Equal(t, tc.filterResult, filterRes)
		})
	}
}

func TestBoolFilterMatchIfKeyMissing(t *testing.T) {
	t.Parallel()

	bf1 := NewBoolFilter()
	bf1.Parse("=true")
	assert.False(t, bf1.MatchIfKeyMissing())

	bf3 := NewBoolFilter()
	bf3.Parse("=true")
	bf3.Parse("=false")
	assert.False(t, bf3.MatchIfKeyMissing())

	bf2 := NewBoolFilter()
	bf2.Parse("=false")
	assert.True(t, bf2.MatchIfKeyMissing())
}

func TestBoolFilterClone(t *testing.T) {
	t.Parallel()

	filter := NewBoolFilter()
	err := filter.Parse("=false")
	require.NoError(t, err)

	copy := filter.Clone()

	opt1 := cmp.AllowUnexported(
		BoolFilter{},
	)
	if !cmp.Equal(filter, copy, opt1) {
		diff := cmp.Diff(filter, copy, opt1)
		t.Errorf("Clone did not produce an identical copy\ndiff: %s", diff)
	}

	// ensure that changes to the copy do not affect the original
	err = copy.Parse("=true")
	require.NoError(t, err)
	if cmp.Equal(filter, copy, opt1) {
		t.Errorf("Changes to copied filter affected the original")
	}
}

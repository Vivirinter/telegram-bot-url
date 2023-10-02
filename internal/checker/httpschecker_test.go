package checker

import "testing"

func TestChecker_IsRedirect(t *testing.T) {
	c := NewChecker()

	testCases := []struct {
		name string
		code int
		want bool
	}{
		{
			name: "Test 1: 302 should be a redirect",
			code: 302,
			want: true,
		},
		{
			name: "Test 2: 200 should not be a redirect",
			code: 200,
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := c.IsRedirect(tc.code); got != tc.want {
				t.Errorf("IsRedirect(%v) = %v; want %v", tc.code, got, tc.want)
			}
		})
	}
}

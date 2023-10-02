package telegram

import "testing"

func Test_isURL(t *testing.T) {
	testCases := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "test_1: correct URL",
			str:  "https://example.com",
			want: true,
		},
		{
			name: "test_2: incorrect URL",
			str:  "this_is_not_url",
			want: false,
		},
		{
			name: "test_3: URL without a scheme",
			str:  "example.com",
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isURL(tc.str); got != tc.want {
				t.Errorf("isURL(%v) = %v; want %v", tc.str, got, tc.want)
			}
		})
	}
}

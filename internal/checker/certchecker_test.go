package checker

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"
)

func Test_isSelfSigned(t *testing.T) {
	testCases := []struct {
		name       string
		cert       *x509.Certificate
		wantResult bool
	}{
		{
			name: "test_1: Self-signed certificate",
			cert: &x509.Certificate{
				Subject: pkix.Name{
					CommonName: "example.com",
				},
				Issuer: pkix.Name{
					CommonName: "example.com",
				},
			},
			wantResult: true,
		},
		{
			name: "test_2: Not self-signed certificate",
			cert: &x509.Certificate{
				Subject: pkix.Name{
					CommonName: "example.com",
				},
				Issuer: pkix.Name{
					CommonName: "Another authority",
				},
			},
			wantResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if result := isSelfSigned(tc.cert); result != tc.wantResult {
				t.Errorf("isSelfSigned(%v) = %v; want %v", tc.cert, result, tc.wantResult)
			}
		})
	}
}

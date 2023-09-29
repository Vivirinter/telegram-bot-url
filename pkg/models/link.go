package models

import (
	"crypto/x509"
	"fmt"
	"strings"
	"time"
)

type Link struct {
	URL            string
	IsHTTPS        bool
	LastChecked    time.Time
	CheckFrequency time.Duration
	Metadata       string
	RedirectURL    string

	CertificateInfo CertificateInfo
	ResponseInfo    ResponseInfo
}

func (l Link) String() string {
	return fmt.Sprintf(
		"URL: %v\n"+
			"HTTPS: %v\n"+
			"\nCertificate Info: \n%v"+
			"\nResponse Info: \n%v",
		l.URL,
		l.IsHTTPS,
		l.CertificateInfo,
		l.ResponseInfo,
	)
}

type CertificateInfo struct {
	IsCertValid        bool
	IsSelfSigned       bool
	Issuer             string
	ValidFrom          time.Time
	ValidTo            time.Time
	Subject            string
	SerialNumber       string
	PublicKeyAlgorithm string
	SignatureAlgorithm string
}

func (ci CertificateInfo) String() string {
	return fmt.Sprintf(
		"\tCertificate is Valid: %v\n"+
			"\tCertificate is Self-Signed: %v\n"+
			"\tIssuer: %v\n"+
			"\tSubject: %v\n"+
			"\tSerial Number: %v\n"+
			"\tPublic Key Algorithm: %v\n"+
			"\tValid from: %v\n"+
			"\tValid to: %v\n"+
			"\tSignature Algorithm: %v\n",
		ci.IsCertValid,
		ci.IsSelfSigned,
		ci.Issuer,
		ci.Subject,
		ci.SerialNumber,
		ci.PublicKeyAlgorithm,
		ci.ValidFrom,
		ci.ValidTo,
		ci.SignatureAlgorithm,
	)
}

func NewCertificateInfo(cert *x509.Certificate) CertificateInfo {
	sigAlg := cert.SignatureAlgorithm.String()
	if cert.SignatureAlgorithm == x509.UnknownSignatureAlgorithm {
		sigAlg = "Unknown"
	}
	return CertificateInfo{
		Issuer:             cert.Issuer.String(),
		ValidFrom:          cert.NotBefore,
		ValidTo:            cert.NotAfter,
		Subject:            cert.Subject.String(),
		SerialNumber:       cert.SerialNumber.String(),
		PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
		SignatureAlgorithm: sigAlg,
	}
}

type ResponseInfo struct {
	Headers         map[string][]string
	StatusCode      int
	SelectedHeaders []string
	HSTS            bool
	WasRedirected   bool
	RedirectedHTTPS bool
}

func (ri ResponseInfo) String() string {
	selectedHeaderLines := "\nSelected Headers: None"
	if len(ri.SelectedHeaders) > 0 {
		selectedHeaderLines = "\nSelected Headers:\n" + strings.Join(ri.SelectedHeaders, "\n")
	}
	return fmt.Sprintf(
		"\tResponse status code: %v\n"+
			"\tWas redirected: %v\n"+
			"\tRedirected to HTTPS: %v\n"+
			"%v",
		ri.StatusCode,
		ri.WasRedirected,
		ri.RedirectedHTTPS,
		selectedHeaderLines,
	)
}

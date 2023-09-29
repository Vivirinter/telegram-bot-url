package checker

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"golang.org/x/crypto/ocsp"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/Vivirinter/telegram-bot-url/pkg/models"
)

const (
	DialerTimeout = 5
)

func NewHttpClient(insecure bool) *http.Client {
	if insecure {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		}
		return &http.Client{Transport: transport}
	}
	return http.DefaultClient
}

func CheckCert(link *models.Link) error {
	parsedURL, err := url.Parse(link.URL)
	if err != nil {
		return fmt.Errorf("Failed to parse URL: %w", err)
	}

	host := parsedURL.Hostname()
	dialer := &net.Dialer{Timeout: time.Second * DialerTimeout}
	client := NewHttpClient(true)

	conn, err := DialWithClient(dialer, "tcp", host+":443", client)
	if err != nil {
		return fmt.Errorf("Failed to dial: %w", err)
	}
	defer func(conn *tls.Conn) {
		_ = conn.Close()
	}(conn)

	if err := checkCertificate(link, host, conn.ConnectionState()); err != nil {
		return fmt.Errorf("Failed to check certificate: %w", err)
	}

	link.IsHTTPS = parsedURL.Scheme == "https"

	return nil
}

func DialWithClient(dialer *net.Dialer, network, addr string, client *http.Client) (*tls.Conn, error) {
	return tls.DialWithDialer(dialer, network, addr, client.Transport.(*http.Transport).TLSClientConfig)
}

func checkCertificate(link *models.Link, host string, state tls.ConnectionState) error {
	cert := state.PeerCertificates[0]
	link.CertificateInfo.IsSelfSigned = isSelfSigned(cert)
	if link.CertificateInfo.IsSelfSigned {
		return errors.New("The certificate is self-signed")
	}

	if err := verifyCert(link, cert, state, host); err != nil {
		return fmt.Errorf("Failed to verify cert: %w", err)
	}

	return nil
}

func isSelfSigned(cert *x509.Certificate) bool {
	return cert.Issuer.CommonName == cert.Subject.CommonName
}

func verifyCert(link *models.Link, cert *x509.Certificate, state tls.ConnectionState, host string) error {
	opts := prepareVerifyOptions(state)

	if _, err := cert.Verify(opts); err != nil {
		return fmt.Errorf("Failed to verify certificate chain: %w", err)
	}

	if err := checkCertificateValidity(cert, host, state); err != nil {
		return fmt.Errorf("Failed to check certificate validity: %w", err)
	}

	link.CertificateInfo.IsCertValid = true

	return nil
}

func prepareVerifyOptions(state tls.ConnectionState) x509.VerifyOptions {
	opts := x509.VerifyOptions{
		Intermediates: x509.NewCertPool(),
	}
	for _, pc := range state.PeerCertificates[1:] {
		opts.Intermediates.AddCert(pc)
	}

	return opts
}

func checkCertificateValidity(cert *x509.Certificate, host string, state tls.ConnectionState) error {
	if time.Now().After(cert.NotAfter) {
		return errors.New("The certificate is expired")
	}

	if time.Now().Before(cert.NotBefore) {
		return errors.New("The certificate is not yet valid")
	}

	if err := cert.VerifyHostname(host); err != nil {
		return fmt.Errorf("Failed to verify hostname: %w", err)
	}

	if err := checkOCSPResponse(state, cert); err != nil {
		return fmt.Errorf("Failed to check OCSP response: %w", err)
	}

	return nil
}

func checkOCSPResponse(state tls.ConnectionState, cert *x509.Certificate) error {
	if state.OCSPResponse != nil && len(state.PeerCertificates) > 1 {
		resp, err := ocsp.ParseResponseForCert(state.OCSPResponse, cert, state.PeerCertificates[1])
		if err != nil {
			return fmt.Errorf("Failed to parse OCSP response: %w", err)
		}
		if resp.Status == ocsp.Revoked {
			return errors.New("The certificate has been revoked")
		}
	}

	return nil
}

package checker

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Vivirinter/telegram-bot-url/pkg/models"
)

type Checker struct {
	client *http.Client
}

func NewChecker() *Checker {
	return &Checker{
		client: createClient(),
	}
}

func createClient() *http.Client {
	timeout := 5 * time.Second

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &http.Client{
		Transport: tr,
		Timeout:   timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func (c *Checker) CheckHTTPS(link *models.Link) error {
	log.Printf("Starting HTTPS check for URL: %s", link.URL)

	resp, err := c.client.Get(link.URL)
	if err != nil {
		log.Printf("Error getting URL: %s. Error: %s", link.URL, err)
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return c.processResponse(resp, link)
}

func (c *Checker) processResponse(resp *http.Response, link *models.Link) error {
	c.setupResponseInfo(resp, link)

	if c.isRedirect(resp.StatusCode) {
		return c.handleRedirect(resp, link)
	}

	link.IsHTTPS = resp.TLS != nil
	if !link.IsHTTPS {
		log.Printf("URL: %s is not redirected and does not use HTTPS", link.URL)
		return nil
	}

	return c.handleHTTPS(resp, link)
}

func (c *Checker) setupResponseInfo(resp *http.Response, link *models.Link) {
	link.ResponseInfo.Headers = resp.Header
	link.ResponseInfo.StatusCode = resp.StatusCode
}

func (c *Checker) isRedirect(status int) bool {
	return status >= 300 && status < 400
}

func (c *Checker) handleRedirect(resp *http.Response, link *models.Link) error {
	loc, err := resp.Location()
	if err != nil {
		log.Printf("Error getting redirection location for URL: %s. Error: %s", link.URL, err)
		return err
	}

	c.setupRedirectInfo(loc, link)

	log.Printf("URL: %s was redirected to: %s", link.URL, loc.String())
	return nil
}

func (c *Checker) setupRedirectInfo(loc *url.URL, link *models.Link) {
	link.RedirectURL = loc.String()
	link.ResponseInfo.WasRedirected = true
	link.ResponseInfo.RedirectedHTTPS = strings.HasPrefix(loc.String(), "https://")
}

func (c *Checker) handleHTTPS(resp *http.Response, link *models.Link) error {
	if !hasCertificates(resp) {
		log.Printf("No peer certificates for URL: %s", link.URL)
		return errors.New("no peer certificates")
	}

	c.setupCertificateInfo(resp, link)

	link.ResponseInfo.HSTS = resp.Header["Strict-Transport-Security"] != nil

	return nil
}

func hasCertificates(resp *http.Response) bool {
	return len(resp.TLS.PeerCertificates) > 0
}

func (c *Checker) setupCertificateInfo(resp *http.Response, link *models.Link) {
	cert := resp.TLS.PeerCertificates[0]
	link.CertificateInfo = models.NewCertificateInfo(cert)
}

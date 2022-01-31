package CertGen

import "crypto/x509"

type Gen interface {
	NewCertificate(cert *x509.Certificate) ([]byte, error)
}

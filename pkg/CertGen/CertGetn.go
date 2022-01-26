package CertGen

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type CertGen struct {
	pool   *x509.CertPool
	CaCert *x509.Certificate
	CaKey  *rsa.PrivateKey
}

func NewCertGenerator(CaCert []byte, CaKey []byte) (CertGen, error) {
	pool := x509.NewCertPool()
	gen := CertGen{}
	cert, err := ParseCert(CaCert)
	switch err != nil {
	case true:
		return CertGen{}, err
	}
	privateKey, err := ParsePKCS1PrivateKey(CaKey)
	switch err != nil {
	case true:
		return CertGen{}, nil
	}
	pool.AddCert(cert)
	gen.pool = pool
	gen.CaCert = cert
	gen.CaKey = privateKey
	return gen, nil
}

func ParseCert(ca []byte) (*x509.Certificate, error) {
	block := DecodePem(ca)
	return x509.ParseCertificate(block.Bytes)
}

func DecodePem(cert []byte) *pem.Block {
	block, _ := pem.Decode(cert)
	return block
}

func ParsePKCS1PrivateKey(key []byte) (*rsa.PrivateKey, error) {
	return x509.ParsePKCS1PrivateKey(DecodePem(key).Bytes)
}

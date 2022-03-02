package CertGen

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
)

type CertGen struct {
	pool   *x509.CertPool
	CaCert *x509.Certificate
	CaKey  *rsa.PrivateKey
}

var max = new(big.Int)

func init() {
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))
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
		return CertGen{}, err
	}
	pool.AddCert(cert)
	gen.pool = pool
	gen.CaCert = cert
	gen.CaKey = privateKey
	return gen, nil
}

func (c CertGen) NewCertificate(cert *x509.Certificate) ([]byte, error) {
	cert.SerialNumber = GenerateUniqueId()
	return x509.CreateCertificate(rand.Reader, cert, c.CaCert, c.CaKey.PublicKey, c.CaKey)
}

func GenerateUniqueId() *big.Int {
	id, _ := rand.Int(rand.Reader, max)
	return id
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

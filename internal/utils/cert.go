package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

const cacheDir = ".cache"

func GenerateSelfSignedCert() (tls.Certificate, error) {
	certFile := filepath.Join(cacheDir, "cert.pem")
	keyFile := filepath.Join(cacheDir, "key.pem")

	if _, err := os.Stat(certFile); err == nil {
		if _, err := os.Stat(keyFile); err == nil {
			return tls.LoadX509KeyPair(certFile, keyFile)
		}
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	tmpl := x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		Subject:      pkix.Name{CommonName: "localhost"},
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	certDer, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDer})
	privPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return tls.Certificate{}, err
	}

	if err := os.WriteFile(certFile, certPem, 0644); err != nil {
		return tls.Certificate{}, err
	}

	if err := os.WriteFile(keyFile, privPem, 0644); err != nil {
		return tls.Certificate{}, err
	}

	cert, err := tls.X509KeyPair(certPem, privPem)
	if err != nil {
		return tls.Certificate{}, err
	}

	return cert, nil
}
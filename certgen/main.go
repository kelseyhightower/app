// Copyright 2016 Google, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

var (
	notBefore, notAfter time.Time
	serialNumberLimit   *big.Int
	rsaBits             = 2048
	host                string
)

var serverSubjectAlternateNames = []string{
	"*.example.com",
	"localhost",
	"127.0.0.1",
}

type certificateConfig struct {
	isCA        bool
	caCert      *x509.Certificate
	caKey       *rsa.PrivateKey
	hosts       []string
	keyUsage    x509.KeyUsage
	extKeyUsage []x509.ExtKeyUsage
}

func init() {
	notBefore = time.Now().Add(-5 * time.Minute).UTC()
	notAfter = notBefore.Add(365 * 24 * time.Hour).UTC()
	serialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)

	flag.StringVar(&host, "host", "", "Comma-separated hostnames and IPs to generate a certificate for")
}

func main() {
	flag.Parse()

	// Generate CA
	caCert, caKey, err := generateCertificate(certificateConfig{
		isCA:        true,
		hosts:       []string{""},
		keyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		extKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	})
	if err != nil {
		log.Fatal(err)
	}

	err = writeCert("ca", caCert, caKey)
	if err != nil {
		log.Fatal(err)
	}

	caParsedCertificates, err := x509.ParseCertificates(caCert)
	if err != nil {
		log.Fatal(err)
	}

	// Generate Server Certificates
	hosts := make([]string, 0)
	for _, h := range strings.Split(host, ",") {
		if h == "" {
			continue
		}
		hosts = append(hosts, h)
	}
	hosts = append(hosts, serverSubjectAlternateNames...)

	serverCert, serverKey, err := generateCertificate(certificateConfig{
		caCert:      caParsedCertificates[0],
		caKey:       caKey,
		hosts:       hosts,
		keyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		extKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	})
	if err != nil {
		log.Fatal(err)
	}

	err = writeCert("server", serverCert, serverKey)
	if err != nil {
		log.Fatal(err)
	}
}

func writeCert(name string, cert []byte, key *rsa.PrivateKey) error {
	certFilename := fmt.Sprintf("%s.pem", name)
	keyFilename := fmt.Sprintf("%s-key.pem", name)

	certFile, err := os.Create(certFilename)
	if err != nil {
		return err
	}
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
	certFile.Close()
	fmt.Printf("wrote %s\n", certFilename)

	keyFile, err := os.OpenFile(keyFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	pem.Encode(keyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	keyFile.Close()
	fmt.Printf("wrote %s\n", keyFilename)
	return nil
}

func generateCertificate(c certificateConfig) ([]byte, *rsa.PrivateKey, error) {
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// Generate the subject key ID
	derEncodedPubKey, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	pubKeyHash := sha1.New()
	pubKeyHash.Write(derEncodedPubKey)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Kubernetes"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		IsCA:                  c.isCA,
		KeyUsage:              c.keyUsage,
		ExtKeyUsage:           c.extKeyUsage,
		BasicConstraintsValid: true,
		SubjectKeyId:          pubKeyHash.Sum(nil),
	}
	if c.hosts[0] != "" {
		template.Subject.CommonName = c.hosts[0]
	}

	if c.isCA {
		c.caCert = &template
		c.caKey = key
	}

	for _, h := range c.hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, c.caCert, &key.PublicKey, c.caKey)
	if err != nil {
		return nil, nil, err
	}

	return derBytes, key, nil
}

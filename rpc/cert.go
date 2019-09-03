package rpc

import (
	"bytes"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"

	"crypto/tls"

	"notekeeper-electron-backend/appdir"
)

func (rpc *Server) createCertificate() bool {
	now := time.Now()
	notBefore := now.Add(-time.Hour * 24)
	// end of ASN.1 time
	endOfTime := time.Date(2049, 12, 31, 23, 59, 59, 0, time.UTC)
	notAfter := endOfTime // notBefore.Add(365 * 24 * time.Hour)

	host, err := os.Hostname()
	if err != nil {
		rpc.Logger.Warn("Error getting hostname - ", err)
		return false
	}

	ipAddresses := []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")}
	dnsNames := []string{host}
	if host != "localhost" {
		dnsNames = append(dnsNames, "localhost")
	}

	addIP := func(ipAddr net.IP) {
		for _, ip := range ipAddresses {
			if bytes.Equal(ip, ipAddr) {
				return
			}
		}
		ipAddresses = append(ipAddresses, ipAddr)
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		rpc.Logger.Warn("error getting interface addresses - ", err)
		return false
	}
	for _, a := range addrs {
		ipAddr, _, err := net.ParseCIDR(a.String())
		if err == nil {
			addIP(ipAddr)
		}
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		rpc.Logger.Warn("Error creating serial number - ", err)
		return false
	}

	privKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		rpc.Logger.Warn("Error generating private key - ", err)
		return false
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"NoteKeeper.io"},
			CommonName:   host,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		DNSNames:              dnsNames,
		IPAddresses:           ipAddresses,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		rpc.Logger.Warn("Error creating certificate - ", err)
		return false
	}

	certBuf := &bytes.Buffer{}
	err = pem.Encode(certBuf, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		rpc.Logger.Warn("Error encoding certificate - ", err)
		return false
	}

	keyBytes, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		rpc.Logger.Warn("Error marshaling key bytes - ", err)
		return false
	}

	/*
		keyPemBlock := pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: keyBytes, // x509.MarshalPKCS1PrivateKey(privKey),
		}
	*/
	keyBuf := &bytes.Buffer{}
	err = pem.Encode(keyBuf, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
	if err != nil {
		rpc.Logger.Warn("Error encoding key - ", err)
		return false
	}

	/*
		certPemBlock := pem.Block{
			Type:  "CERTIFICATE",
			Bytes: derBytes,
		}
	*/
	/*
		certPath := filepath.Join(appdir.AppDataPath(), "certificate")
		certOut, err := os.Create(certPath)
		if err != nil {
			rpc.Logger.Warn("Error creating certificate file - ", err)
			return false
		}

		pem.Encode(certOut, &certPemBlock)
		certOut.Close()
	*/
	//rpc.Certificate, err = tls.X509KeyPair(pem.EncodeToMemory(&certPemBlock), pem.EncodeToMemory(&keyPemBlock))
	rpc.Certificate, err = tls.X509KeyPair(certBuf.Bytes(), keyBuf.Bytes())
	if err != nil {
		rpc.Logger.Warn("Error converting certificate - ", err)
		return false
	}

	certPath := filepath.Join(appdir.AppDataPath(), "certificate")
	err = ioutil.WriteFile(certPath, certBuf.Bytes(), 0600)
	if err != nil {
		rpc.Logger.Warn("Error writing certificate - ", err)
		return false
	}

	return true
}

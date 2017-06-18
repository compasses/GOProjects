package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
	"net/http"
	"errors"
	"io/ioutil"
)

var (
	validFrom  = ""
	validFor   = 365 * 24 * time.Hour
	isCA       = true
	rsaBits    = 2048
	ecdsaCurve = "P256"
)

func pemBlockForKey(p interface{}) *pem.Block {
	switch k := p.(type) {
	case *rsa.PrivateKey:
		b, err := x509.MarshalPKIXPublicKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal RSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: b}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprint(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func CheckExist(certPath, keyPath string) error {
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return err
	}

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return err
	}

	return nil
}

func publicKey(p interface{}) interface{} {
	switch k := p.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func CertTemplate() (*x509.Certificate, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.New("failed to generate serial number: " + err.Error())
	}

	tmpl := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"Jet He."}},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(validFor), // valid for an hour
		BasicConstraintsValid: true,
	}
	return &tmpl, nil

}
func Generate(certPath, keyPath, host string) error {
	var priv interface{}
	var err error

	switch ecdsaCurve {
	case "":
		priv, err = rsa.GenerateKey(rand.Reader, rsaBits)
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		fmt.Fprintf(os.Stderr, "Unrecognized elliptic curve:%q", ecdsaCurve)
		os.Exit(1)
	}

	if err != nil {
		log.Printf("failed to generate private key: %s", err)
		return err
	}

	var startValid time.Time
	if len(validFrom) == 0 {
		startValid = time.Now()
	} else {
		startValid, err = time.Parse("Jan 2 15:04:05 2006", validFrom)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse creation date: %s\n", err)
			return err
		}
	}
	validEnd := startValid.Add(validFor)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	fmt.Printf("SerialNumberLimit :%d\n", serialNumberLimit)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Printf("failed to generate serial number:%s\n", err)
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Jet He"},
		},
		NotBefore:             startValid,
		NotAfter:              validEnd,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		log.Printf("failed to create certificate %s\n", err)
		return err
	}
	certOut, err := os.Create(certPath)
	if err != nil {
		log.Printf("failed to open "+certPath+" for writing: %s", err)
		return err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.OpenFile(keyPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		log.Print("failed to open "+keyPath+" for writing:", err)
		return err
	}
	pem.Encode(keyOut, pemBlockForKey(priv))
	keyOut.Close()
	log.Print("Writen key.pem\n")
	return nil
}

func testCertPemFile(cert string)  {
	fmt.Println("test cert pem file :" + cert)

	c, err := ioutil.ReadFile(cert)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read "+cert +" file failed %v", err)
		return
	}

	cb, re := pem.Decode(c)
	if cb == nil {
		fmt.Fprintf(os.Stderr, "decode cert file failed %v", re)
		return
	}

	fmt.Println("get certificate type is " + cb.Type)

	ct, err := x509.ParseCertificate(cb.Bytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse key error %v", err)
	}
	pub := ct.PublicKey

	switch k := pub.(type) {
	case *rsa.PublicKey:
		fmt.Println("it's a RSA public key")
	case *ecdsa.PublicKey:
		fmt.Println("It's a ecdsa public key ")
	default:
		fmt.Fprintf(os.Stderr, "unknow type %v", k)
		return
	}

}

func testKeyPemFile(cert string)  {
	fmt.Println("test key pem file :" + cert)

	c, err := ioutil.ReadFile(cert)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read "+cert +" file failed %v", err)
		return
	}

	cb, re := pem.Decode(c)
	if cb == nil {
		fmt.Fprintf(os.Stderr, "decode key pem file failed %v", re)
		return
	}

	fmt.Println("get key type is " + cb.Type)

	ct, err := x509.ParseECPrivateKey(cb.Bytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse EC private key error %v", err)
		return
	}

	pub := ct.PublicKey

	fmt.Printf("Curve is %v\n", pub.Curve)

	//plainText := []byte("The bourgeois human is a virus on the hard drive of the working robot!")

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there!")
}

func main() {
	// Check if the cert files are available.
	err := CheckExist("cert.pem", "key.pem")
	// If they are not available, generate new ones.
	if err != nil {
		err = Generate("cert.pem", "key.pem", "127.0.0.1:8081")
		if err != nil {
			log.Fatal("Error: Couldn't create https certs.")
		}
	}
	testCertPemFile("cert.pem")
	testKeyPemFile("key.pem")

	http.HandleFunc("/", handler)
	err = http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "got error %v", err)
	}

}
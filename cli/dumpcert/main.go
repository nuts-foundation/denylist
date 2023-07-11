// dumpcert provides information about a given certificate which is
// needed in order to add a nuts denylist entry
//
// To dump a usable denylist entry for a given PEM certificate file:
//   go run cli/dumpcert/main.go --cert /path/to/x509.crt
package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nuts-foundation/denylist"

	"github.com/lestrrat-go/jwx/jwk"
)

// Define a variable for the certificate file to be inspected
var certPath string

// Setup command line flags when the module is imported
func init() {
	// Define the flags for command line arguments
	flag.StringVar(&certPath, "cert", "", "path of certificate to inspect")

	// Parse the command line flags
	flag.Parse()
}

func main() {
	// Ensure a certificate path was passed
	if certPath == "" {
		log.Fatal("--cert flag is required")
	}

	// Read the certificate file
	pemBytes, err := os.ReadFile(certPath)
	if err != nil {
		log.Fatalf("Failed to open %s: %v", certPath, err)
	}

	// Parse the PEM block to get the certificate raw bytes
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		log.Fatalf("Malformed PEM certificate block: %s", block.Type)
	}

	// Parse the certificate data as x509
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse x509 data: %s", err)
	}

	// Generate a denylist entry using the shared data structure
	var entry denylist.Entry

	// Set the entry fields
	entry.Issuer = cert.Issuer.String()
	entry.SerialNumber = cert.SerialNumber.String()
	entry.JWKThumbprint = certKeyThumbprint(cert)
	entry.Reason = "TODO: Give me one reason, give me just one reason why"

	// Dump the JSON for the new entry
	jsonBytes, err := json.MarshalIndent(&entry, "", "    ")
	fmt.Println(string(jsonBytes))
}

// certKeyThumbprint returns the JWK thumbprint of the public key contained within an X.509 certificate or empty string on error
func certKeyThumbprint(cert *x509.Certificate) string {
	// Open the public key of the certificate using the JWK library, returning an
	// empty string if that fails
	key, err := jwk.New(cert.PublicKey)
	if err != nil {
		return ""
	}
	
	// Calculate and internally set the key thumbprint
	jwk.AssignKeyID(key)

	// Retrieve and return the key thumbprint
	keyID, _ := key.Get(jwk.KeyIDKey)
	return keyID.(string)
}

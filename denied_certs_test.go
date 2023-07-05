package denylist

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"testing"

        "github.com/lestrrat-go/jwx/jwa"
        "github.com/lestrrat-go/jwx/jwk"
        "github.com/lestrrat-go/jwx/jws"
)

// Define the paths for the necessary input files
const publicKeyPath = `out/pubkey.pem`
const denylistPath = `out/denylist.jws`

// loadCerts loads the certificate files from the filesystem
func loadCerts() ([]*x509.Certificate, error) {
	// Scan for certificates which should be blocked
	paths, err := filepath.Glob("certs/*")
	if err != nil {
		return nil, err
	}

	// Allocate a variable for storing the loaded certificates
	var certs []*x509.Certificate

	// Load the discovered certificates
	for _, path := range paths {
		// Read the certificate file
		pemBytes, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open %s: %w", path, err)
		}

		// Parse the PEM block to get the certificate raw bytes
		block, _ := pem.Decode(pemBytes)
		if block == nil || block.Type != "CERTIFICATE" {
			return nil, fmt.Errorf("Malformed PEM certificate block: %s", block.Type)
		}

		// Parse the certificate data as x509
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}

		// Append the loaded certificate to the slice
		certs = append(certs, cert)
	}

	// Return the loaded certificates
	return certs, nil
}

// TestDeniedCerts ensures that certs listed in the certs/ directory are included in the deny list @ out/denylist.jws
func TestDeniedCerts(t *testing.T) {
	// Read the trusted public key
	keyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		t.Fatalf("failed to read public key: %v", err)
	}

	// Parse the public key
	key, err := jwk.ParseKey(keyBytes, jwk.WithPEM(true))
	if err != nil {
		t.Fatalf("failed to parse public key: %v", err)
	}

	// Read the certificates from the filesystem
	certs, err := loadCerts()
	if err != nil {
		t.Fatalf("failed to load certs: %v", err)
	}

	// Load the denylist
	var denylist []Entry
	rawDenylist, err := os.ReadFile(denylistPath)
	if err != nil {
		t.Fatalf("failed to read denylist: %v", err)
	}

	// Parse the denylist as a JWS, verifying the signature
	payload, err := jws.Verify(rawDenylist, jwa.EdDSA, key)
	if err != nil {
		t.Fatalf("failed to parse/verify denylist JWS: %v", err)
	}

	// Unmarshal the denylist JSON
	if err := json.Unmarshal(payload, &denylist); err != nil {
		t.Fatalf("failed to parse denylist JSON: %v", err)
	}

	// Loop through the loaded certs, ensuring each one is present in the denylist
	for _, cert := range certs {
		// Calculate the public key thumbprint for this certificate
		jwkThumbprint := certKeyThumbprint(cert)

		// Run a subtests for this certificate
		t.Run(fmt.Sprintf("%s %s", cert.Issuer, cert.SerialNumber), func (t *testing.T) {
			// Loop through the denylist entries, looking for a matching entry
			for _, entry := range denylist {
				// Ensure the issuer, serial number, and key thumbprint all match this entry
				if cert.Issuer.String() == entry.Issuer && cert.SerialNumber.String() == entry.SerialNumber && jwkThumbprint == entry.JWKThumbprint {
					// Return from the function when a matching entry has been found
					return
				}
			}

			// Report a fatal test error when no matching entry was found
			t.Fatalf("no matching entry in denylist for %v", cert)
		})
	}
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

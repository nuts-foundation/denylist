/*
 * Nuts node
 * Copyright (C) 2023 Nuts community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
)

// denylistEntry contains parameters for an X.509 certificate that must not be accepted for TLS connections
type denylistEntry struct {
	// Issuer is a string representation (x509.Certificate.Issuer.String()) of the certificate
	Issuer string `json:"issuer"`

	// SerialNumber is a string representation (x509.Certificate.SerialNumber.String()) of the certificate
	SerialNumber string `json:"serialnumber"`

	// JWKThumbprint is an identifier of the public key per https://www.rfc-editor.org/rfc/rfc7638
	JWKThumbprint string `json:"jwkthumbprint"`

	// Reason is the reason for which the certificate is being denied
	Reason string `json:"reason"`
}

func encodeDenylist(entries []denylistEntry) string {
	// Encode the denylist as JSON
	payload, err := json.Marshal(&entries)
	if err != nil {
		log.Fatalf("error marshalling JSON: %s", err)
	}

	// Read the private key from the environment variable
	privateKeyPEM := os.Getenv("DENYLIST_PRIVATEKEY_PEM")
	if privateKeyPEM == "" {
		log.Fatal("DENYLIST_PRIVATEKEY_PEM environment variable must not be empty")
	}

	// Parse the private key for signing the denylist
	key, err := jwk.ParseKey([]byte(privateKeyPEM), jwk.WithPEM(true))
	if err != nil {
		log.Fatalf("error parsing key: %s", err)
	}

	// Sign the denylist as a JWS Message
	compactJWS, err := jws.Sign(payload, jwa.EdDSA, key)
	if err != nil {
		log.Fatalf("error signing payload: %s", err)
	}

	// Return the compact encoded JWS message
	return string(compactJWS)
}

func main() {
	// Load the denylist entries from JSON
	data, err := os.ReadFile("denylist.json")
	if err != nil {
		log.Fatalf("error reading denylist.json: %s", err)
	}

	var entries []denylistEntry

	// Build the denylist contents
	err = json.Unmarshal(data, &entries)
	if err != nil {
		log.Fatalf("error unmarshalling JSON: %s", err)
	}

	// Encode and print the denylist
	encoded := encodeDenylist(entries)
	fmt.Println(encoded)
}

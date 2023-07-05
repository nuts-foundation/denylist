package denylist

// denylistEntry contains parameters for an X.509 certificate that must not be accepted for TLS connections
type Entry struct {
	// Issuer is a string representation (x509.Certificate.Issuer.String()) of the certificate
	Issuer string `json:"issuer"`

	// SerialNumber is a string representation (x509.Certificate.SerialNumber.String()) of the certificate
	SerialNumber string `json:"serialnumber"`

	// JWKThumbprint is an identifier of the public key per https://www.rfc-editor.org/rfc/rfc7638
	JWKThumbprint string `json:"jwkthumbprint"`

	// Reason is the reason for which the certificate is being denied
	Reason string `json:"reason"`
}


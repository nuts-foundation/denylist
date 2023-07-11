# denylist
A filter for PKI verification

## What is the denylist?

The denylist is a list of certificates which should be blocked/rejected in the nuts-node PKI scheme. This provides a mechanism for organizations to subscribe to secure, cryptographically signed published lists of certificates which should not be considered valid.

## Where is the denylist hosted?

The current denylist is located [in this repo](./denylist-out/denylist.jws), along with all of the tools required to generate it. An externally stored secure private key is also required in order to generate a valid denylist.

For maximum performance and reliability a CDN should be used when configuring production systems to use the denylist.

## What is a jws file?

JWS, short for JSON Web Signature, is a cryptographically secure document standard defined in [RFC7515](https://datatracker.ietf.org/doc/html/rfc7515). JWS allows you to verify the authenticity of the denylist using the public key. In the simplest terms a JWS is a JSON document for which the original author can be proven.

## Why not use a CRL or OCSP?

CRL's are created in a way that is only valid at the point of the certificate issuer. OCSP is not easily prefetchable without a cooperating client implementing OCSP stapling. Both CRL's and OCSP bring significant risk to the uptime of a network. This denylist has been implemented in such a way that downtime is tolerated and the entire contents are prefetched and cacheable.

## Updating the denylist

1. Add an entry to `config/certs.json` containing the details of the certificate to be blocked. An entry template can be generated with the command `go run cli/dumpcert/main.go --cert /path/to/x509.crt`
2. If possible add a copy of the certificate to be blocked in PEM format to the `certs/` directory
3. Run `make denylist` which will update `out/denylist.jws`.
4. Run `make test` to ensure the resulting denylist works and blocks the certificates as expected.
5. Merge the resulting changes to the `main` branch with a pull request.

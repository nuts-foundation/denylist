# Denylist
A filter for PKI verification

## What is the denylist?

The denylist is a list of certificates which should be blocked/rejected in the nuts-node PKI scheme. 
This provides a mechanism for organizations to subscribe to secure, 
cryptographically signed published lists of certificates which should not be considered valid.

## Where is the denylist hosted?

The current denylist is located [in this repo](denylist/denylist.jws), along with all of the tools required to generate it. 
An externally stored secure private key is also required in order to generate a valid denylist.

For maximum performance and reliability a CDN should be used when configuring production systems to use the denylist.

## What is a jws file?

JWS, short for JSON Web Signature, is a cryptographically secure document standard defined in [RFC7515](https://datatracker.ietf.org/doc/html/rfc7515).
JWS allows you to verify the authenticity of the denylist using the public key. 
In the simplest terms a JWS is a JSON document for which the original author can be proven.

## Why not use a CRL or OCSP?

CRL's are created in a way that is only valid at the point of the certificate issuer. 
OCSP is not easily prefetchable without a cooperating client implementing OCSP stapling. 
Both CRL's and OCSP bring significant risk to the uptime of a network. 
This denylist has been implemented in such a way that downtime is tolerated and the entire contents are prefetched and cacheable.

## Supported key types

The denylist currently only supports PEM encoded ed25519 keys.

# Instructions for adding a certificate to the denylist

## Get certificate to block
The first step is to find the certificate to block.
Identify the peer whose certificate needs to be blocked and get the certificate from the peer diagnostics.
Make the following API call to a nuts-node to get the peer diagnostics.

```
GET /internal/network/v1/diagnostics/peers
```

Lookup the peer by its `peerID` and copy the certificate.
Store this PEM encoded certificate in its own file in the `certs/` directory.

## Add certificate to denylist
The unsigned denylist is maintained at `config/certs.json`.

Make a new entry using
```shell
# generate a denylist entry
go run cli/dumpcert/main.go --cert <path-to-cert>
```

This will produce something like:
```json
{
    "issuer": "CN=Nuts Development Network Root CA",
    "serialnumber": "86909351591313461694061157857735660045310273765",
    "jwkthumbprint": "xaVihytdhiID5wqCOsQbbb2_GuIb0r7m1GUgPIdBDjo",
    "reason": "TODO: Give me one reason, give me just one reason why"
}
```

Copy this into `config/certs.json` and set the reason for blocking the certificate.


## Sign denylist
The key to sign the denylist is stored in the `NutsDenylist` vault on 1password in the `PrivateKey` item.

### 1password CLI (preferred)
Install v2 of the 1password CLI by following the [instructions](https://developer.1password.com/docs/cli/get-started/) if needed.
```shell
# sign in to the correct account
op signin

# point the env variable to the correct value in 1password
export DENYLIST_PRIVATEKEY_PEM="op://NutsDenylist/PrivateKey/password"

# sign by letting 1password substitute the key
op run make denylist
```

### Copy key to CLI (only if absolutely needed)
If the 1password CLI is not an option for some reason, the denylist can also be signed after manually setting the env variable.
```
# get the key from 1password and subsitute below to sign the new denylist.
DENYLIST_PRIVATEKEY_PEM="<private-key>" make denylist
```


## Publish denylist
The previous step updated the signed denylist at `denylist/denylist.jws`. (DO NOT CHANGE THIS FILE LOCATION/NAME)
Assuming the certificate to block was added to `certs/`, validate that the new denylist actually identifies the certificate.
```shell
# validate that the added certifcate is on the denylist
make test
```
This loops over all the certificates in `certs/` and checks that they are on the denylist.

If the test passes, merge the resulting changes to the `main` branch with a pull request.


## Clear CDN cache
The denylist is hosted on Github but distributed using [www.jsdelivr.com](https://www.jsdelivr.com), which is "A free CDN for open source projects".
This CDN operates on top of multiple major CDNs and will provide us with a better uptime than Github can.
The [default cache for jsdelivr](https://www.jsdelivr.com/documentation#id-caching) is 12 hours for files on Github branches.
TODO: contact jsdelivr to see if/how we can speed this up.

The denylist file can be found at https://cdn.jsdelivr.net/gh/nuts-foundation/denylist@main/denylist/denylist.jws. 
This points to the `main` branch on the Github repo.


## Troubleshooting
- `error parsing key: failed to parse PEM encoded key: failed to decode PEM data`:  
  PEM data in an env variable can be a bit iffy. The header and footer should be separated from the key using a literal backslash followed by literal n (so `\n`) which are converted to newline characters after the environment variable is read.
  So `-----BEGIN PRIVATE KEY-----\n<the-actual-key>\n-----END PRIVATE KEY-----`

- `"NutsDenylist" isn't a vault in this account. Specify the vault with its ID or name.`:
  You are probably signed in to your personal 1password account. Change the active account using `op signin`


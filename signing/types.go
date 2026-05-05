package signing

import "crypto/ecdsa"

type SignatureFunc func(signer string, fn func(key *ecdsa.PrivateKey) ([]byte, error)) ([]byte, error)

package libcoin

// DenomSig represents the exchange's denomination signature on a coin.
//
// Taler supports both RSA and Clause-Schnorr (CS) signatures. The Cipher field
// is a discriminator that determines which fields are present. All variant fields
// use omitempty so the struct handles both types. MVP mock data uses RSA only.
type DenomSig struct {
	Cipher       string `json:"cipher"`
	RSASignature string `json:"rsa_signature,omitempty"`
	CSSignatureR string `json:"cs_signature_r,omitempty"`
	CSSignatureS string `json:"cs_signature_s,omitempty"`
}

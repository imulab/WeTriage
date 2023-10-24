package route

// Properties is the configuration property for the route package.
type Properties struct {
	// Token is the value registered with WeCom to verify incoming message signatures.
	Token string
	// AesEncodingKey is the bytes of the AES encryption key used to encrypt/decrypt incoming messages.
	AesEncodingKey []byte
}

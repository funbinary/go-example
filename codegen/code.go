package codegen

//go:generate codegen -type=int
//go:generate codegen -type=int -doc -output ../docs/guide/zh-CN/api/error_code_generated.md

// common: authorization and authentication errors.
const (
	// ErrEncrypt - 401: Error occurred while encrypting the user password.
	ErrEncrypt int = iota + 100201
)

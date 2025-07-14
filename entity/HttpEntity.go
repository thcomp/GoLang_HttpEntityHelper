package entity

type HttpEntityType int

const (
	Unknown HttpEntityType = iota
	JSONRPC_Request
	JSONRPC_Response
	JSONRPC_Error
	MultipartFormData
	UrlEncoding
)

type HttpEntity interface {
	EntityType() HttpEntityType
}

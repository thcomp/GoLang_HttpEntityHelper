package httpentityhelper

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/thcomp/GoLang_HttpEntityHelper/jsonrpc"
	"github.com/thcomp/GoLang_HttpEntityHelper/multipart"
	"github.com/thcomp/GoLang_HttpEntityHelper/urlenc"
	ThcompUtility "github.com/thcomp/GoLang_Utility"
)

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

type HttpEntityParser interface {
	Parse(interface{}) (HttpEntity, error)
}

type HttpEntityHelper struct {
	request      *http.Request
	respose      *http.Response
	parsers      [][]HttpEntityParser
	mimeType     *string
	buffer       *bytes.Buffer
	parsedEntity HttpEntity
}

func NewHttpEntityHelper(data interface{}, reusable bool) (ret *HttpEntityHelper, retErr error) {
	if data == nil {
		retErr = fmt.Errorf("data cannot be nil")
	} else {
		ret = &HttpEntityHelper{
			parsers: [][]HttpEntityParser{
				[]HttpEntityParser{jsonrpc.NewJSONRPCParser(), &multipart.MultipartFormParser{}, &urlenc.URLEncParser{}},
			},
		}
		if reusable {
			ret.buffer = bytes.NewBuffer(nil)
		}

		switch v := data.(type) {
		case *http.Request:
			ret.request = v
		case *http.Response:
			ret.respose = v
		default:
			retErr = fmt.Errorf("unsupported data type: %T", v)
			ret = nil
		}
	}

	return
}

func (helper *HttpEntityHelper) GetMimeType() string {
	if helper.mimeType != nil {
		return *helper.mimeType
	} else {
		helper.mimeType = ThcompUtility.ToStringPointer(
			ThcompUtility.TernaryOpStringFunc(
				helper.request != nil,
				func() string {
					return helper.request.Header.Get("Content-Type")
				},
				func() string {
					return helper.respose.Header.Get("Content-Type")
				},
			),
		)
		return *helper.mimeType
	}
}

func (helper *HttpEntityHelper) RegistParser(parser HttpEntityParser, priority uint) {
	if parser != nil {
		if priority >= uint(len(helper.parsers)) {
			for i := len(helper.parsers); i <= int(priority); i++ {
				if i == int(priority) {
					helper.parsers = append(helper.parsers, []HttpEntityParser{parser})
				} else {
					helper.parsers = append(helper.parsers, []HttpEntityParser{})
				}
			}
		} else {
			helper.parsers[priority] = append(helper.parsers[priority], parser)
		}
	}
}

func (helper *HttpEntityHelper) HttpEntity() HttpEntity {
	if helper.parsedEntity != nil {
		// no-op
	} else if helper.request == nil && helper.request.Body == nil && helper.respose == nil && helper.respose.Body == nil {
		ThcompUtility.LogfW("HttpEntity is nil, request and response are nil or empty")
	} else {
		if helper.buffer != nil {
			body, _ := ThcompUtility.TernaryOpInterface(helper.request.Body != nil, helper.request.Body, helper.respose.Body).(io.Reader)
			helper.buffer.Reset()
			if _, err := helper.buffer.ReadFrom(body); err != nil {
				fmt.Printf("Error reading request body: %v\n", err)
			} else {
				if helper.request.Body != nil {
					helper.request.Body = io.NopCloser(helper.buffer)
				} else if helper.respose.Body != nil {
					helper.respose.Body = io.NopCloser(helper.buffer)
				}
			}
		}

		for _, parserList := range helper.parsers {
			for _, parser := range parserList {
				if entity, err := parser.Parse(helper.request); err == nil {
					helper.parsedEntity = entity
					break
				} else {
					fmt.Printf("Error parsing request body: %v\n", err)
				}
			}

			if helper.parsedEntity != nil {
				break
			}
		}
	}

	return helper.parsedEntity
}

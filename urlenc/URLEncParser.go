package urlenc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	root "github.com/thcomp/GoLang_HttpEntityHelper"
)

type URLEncParser struct {
}

func (parser *URLEncParser) parseEntity(reader io.Reader) (ret *URLEncData, retErr error) {
	urlEncData := URLEncData{}
	if valueBytes, readErr := io.ReadAll(reader); readErr == nil {
		if queryValues, parseErr := url.ParseQuery(string(valueBytes)); parseErr == nil {
			urlEncData.queryValues = &queryValues
			ret = &urlEncData
		} else {
			retErr = parseErr
		}
	} else {
		retErr = readErr
	}

	return
}

func (parser *URLEncParser) Parse(obj interface{}) (ret root.HttpEntity, retErr error) {
	reader := io.Reader(nil)
	header := http.Header(nil)
	url := (*url.URL)(nil)

	switch v := obj.(type) {
	case *http.Request:
		header = v.Header
		reader = v.Body
		url = v.URL
	case *http.Response:
		header = v.Header
		reader = v.Body
	default:
		retErr = fmt.Errorf("can not parse on URLEnc Request, unsupported type: %v", v)
	}

	if retErr == nil {
		if contentTypeValue := header.Get(`Content-type`); contentTypeValue != `` {
			contentTypeValue = strings.ToLower(contentTypeValue)
			if strings.HasPrefix(contentTypeValue, `application/x-www-form-urlencoded`) {
			} else {
				if len(url.RawQuery) > 0 {
					reader = bytes.NewReader([]byte(url.RawQuery))
				} else {
					retErr = fmt.Errorf("no url encoded value in request")
				}
			}
		} else {
			if len(url.RawQuery) > 0 {
				reader = bytes.NewReader([]byte(url.RawQuery))
			} else {
				retErr = fmt.Errorf("no url encoded value in request")
			}
		}

		if reader != nil {
			return parser.parseEntity(reader)
		}
	}

	return
}

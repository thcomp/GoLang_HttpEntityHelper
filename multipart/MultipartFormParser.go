package multipart

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/thcomp/GoLang_HttpEntityHelper/entity"
	ThcompUtility "github.com/thcomp/GoLang_Utility"
)

type MultipartFormParser struct {
	cacheEditorFactory ThcompUtility.CacheEditorFactory
}

func (parser *MultipartFormParser) CacheEditorFactory(cacheEditorFactory ThcompUtility.CacheEditorFactory) {
	parser.cacheEditorFactory = cacheEditorFactory
}

func (parser *MultipartFormParser) Parse(obj interface{}) (ret entity.HttpEntity, retErr error) {
	formData := (*MultipartFormData)(nil)

	switch v := obj.(type) {
	case *http.Request:
		if multipartHelper, err := ThcompUtility.NewMultipartHelperFromHttpRequest(v); err == nil {
			formData = &MultipartFormData{helper: multipartHelper}
		} else {
			retErr = err
		}
	case *http.Response:
		if contentTypeValue := v.Header.Get(`Content-type`); contentTypeValue != `` {
			lowerContentTypeValue := strings.ToLower(contentTypeValue)
			if strings.HasPrefix(lowerContentTypeValue, `multipart/form-data`) {
				boundaryText := (*string)(nil)
				partialTexts := strings.Split(contentTypeValue, `boundary=`)
				if len(partialTexts) >= 2 {
					for i := 1; i < len(partialTexts); i++ {
						prevPart := strings.TrimRight(partialTexts[i-1], " \t")
						if prevPart[len(prevPart)-1] == ';' {
							subPartialTexts := strings.Split(partialTexts[i], ";")
							boundaryText = &subPartialTexts[0]
							break
						}

					}
				}

				if boundaryText != nil {
					cacheEditorFactory := parser.cacheEditorFactory
					if parser.cacheEditorFactory == nil {
						cacheEditorFactory = ThcompUtility.NewMemoryCacheEditorFactory()
					}

					if multipartHelper, err := ThcompUtility.NewMultipartHelper(v.Body, *boundaryText, cacheEditorFactory); err == nil {
						ret = &MultipartFormData{
							helper: multipartHelper,
						}
					} else {
						retErr = err
					}
				} else {
					retErr = fmt.Errorf("not exist boundary text in content-type header: %s", contentTypeValue)
				}
			} else {
				retErr = fmt.Errorf("not multipart/form-data: %s", contentTypeValue)
			}
		} else {
			retErr = fmt.Errorf("not exist content-type header")
		}
	}

	return formData, retErr
}

func (parser *MultipartFormParser) IsMultipartFormData(headers http.Header) (ret bool) {
	mimetype := headers.Get(`Content-type`)
	lowerMimetype := strings.ToLower(mimetype)
	if strings.HasPrefix(lowerMimetype, `multipart/form-data`) {
		ret = true
	}

	return
}

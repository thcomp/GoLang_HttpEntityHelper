package urlenc

import (
	"net/url"

	root "github.com/thcomp/GoLang_HttpEntityHelper"
)

type URLEncData struct {
	queryValues *url.Values
}

func (encData *URLEncData) EntityType() root.HttpEntityType {
	return root.UrlEncoding
}

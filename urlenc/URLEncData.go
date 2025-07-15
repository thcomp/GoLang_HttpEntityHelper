package urlenc

import (
	"net/url"

	"github.com/thcomp/GoLang_HttpEntityHelper/entity"
)

type URLEncData struct {
	queryValues *url.Values
}

func (encData *URLEncData) EntityType() entity.HttpEntityType {
	return entity.UrlEncoding
}

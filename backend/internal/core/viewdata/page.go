package viewdata

import (
	"fmt"
	"net/http"
)

type MetaData struct {
	Title        string
	Description  string
	CanonicalURL string
	ActiveLink   string
}

type HomePageData struct {
	MetaData
	Feature string
}

type AboutPageData struct {
	MetaData
}

type ContactPageData struct {
	MetaData
	ContactEmail string
}

func NewMeta(r *http.Request, title, description string) MetaData {
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}

	return MetaData{
		Title:        title,
		Description:  description,
		CanonicalURL: fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.Path),
		ActiveLink:   r.URL.Path,
	}
}

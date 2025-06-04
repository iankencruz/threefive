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

type LoginPageData struct {
	MetaData
	Email  string
	Errors map[string]string
}

type RegisterPageData struct {
	MetaData  MetaData
	FirstName string
	LastName  string
	Email     string
	Errors    map[string]string
}

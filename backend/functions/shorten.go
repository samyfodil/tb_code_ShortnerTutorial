package lib

import (
	"io"
	"strings"

	"github.com/taubyte/go-sdk/database"
	"github.com/taubyte/go-sdk/event"
	"github.com/taubyte/utils/multihash"
)

//go:generate go get github.com/mailru/easyjson
//go:generate go install github.com/mailru/easyjson/...@latest
//go:generate easyjson -omit_empty -all ${GOFILE}

const (
	baseUrl = "daatrnjz0.g.tau.link"
	minLen  = 5
)

type Body struct {
	BaseUrl string `json:"base_url"`
	URL     string `json:"url"`
}

type Response struct {
	URL    string `json:"url"`
	Short  string `json:"short"`
	Exists bool   `json:"exists"`
}

/* POST /shorten */
//export shorten
func shorten(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	bodyData, err := io.ReadAll(h.Body())
	if err != nil {
		return 1
	}

	body := &Body{}
	err = body.UnmarshalJSON(bodyData)
	if err != nil {
		return 1
	}

	if body.BaseUrl != baseUrl {
		return 1
	}

	urlHash := multihash.Hash(body.URL)

	proposed := strings.ToLower(urlHash[len(urlHash)-minLen:])

	db, err := database.New("urls")
	if err != nil {
		return 1
	}
	defer db.Close()

	exists := false
	_, err = db.Get(proposed)
	if err == nil {
		exists = true
	} else {
		err = db.Put(proposed, []byte(body.URL))
		if err != nil {
			return 1
		}
	}

	res, err := Response{
		URL:    "https://" + body.BaseUrl + "/?s=" + proposed,
		Short:  proposed,
		Exists: exists,
	}.MarshalJSON()
	if err != nil {
		return 1
	}

	h.Write(res)
	h.Return(200)

	return 0
}

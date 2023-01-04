package lib

import (
	"github.com/taubyte/go-sdk/database"
	"github.com/taubyte/go-sdk/event"
)

/* GET /?s=<short> */
//export redirect
func redirect(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	short, err := h.Query().Get("s")
	if err != nil {
		return 1
	}

	db, err := database.New("urls")
	if err != nil {
		return 1
	}
	defer db.Close()

	url, err := db.Get(short)
	if err != nil {
		return 1
	}

	h.Redirect(string(url)).Temporary()

	return 0
}

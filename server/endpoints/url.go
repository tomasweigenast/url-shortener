package endpoints

import (
	"log"
	"net/http"

	"tomasweigenast.com/url-shortener/services"
)

func Url(w http.ResponseWriter, req *http.Request) {
	path := req.PathValue("url")
	if len(path) == 0 {
		w.Write([]byte("missing url"))
		return
	}

	url, err := services.LinksService().FetchUrl(req.Context(), path)
	if err != nil {
		log.Println("link", path, "not found. error:", err)
		http.NotFound(w, req)
		return
	}

	http.Redirect(w, req, url, http.StatusFound)
}

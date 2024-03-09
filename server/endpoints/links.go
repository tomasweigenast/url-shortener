package endpoints

import (
	"encoding/json"
	"net/http"

	"tomasweigenast.com/url-shortener/models"
	"tomasweigenast.com/url-shortener/server/middleware"
	"tomasweigenast.com/url-shortener/server/response"
	"tomasweigenast.com/url-shortener/services"
)

func GetLinks(w http.ResponseWriter, req *http.Request) {

	uid := middleware.GetUid(req)

	links, err := services.LinksService().ListLinks(req.Context(), uid)
	if err != nil {
		response.Failed(w, models.StringError{Reason: err.Error()}, http.StatusNotFound)
		return
	}

	response.Successful(w, links)
}

func CreateLink(w http.ResponseWriter, req *http.Request) {
	uid := middleware.GetUid(req)

	decoder := json.NewDecoder(req.Body)

	// parse
	var body models.CreateLink
	err := decoder.Decode(&body)
	if err != nil {
		response.Failed(w, models.StringError{Reason: "unable to parse request body"}, http.StatusBadRequest)
		return
	}

	url, err := services.LinksService().CreateLink(req.Context(), uid, body)
	if err != nil {
		response.Failed(w, models.ParseError(err), http.StatusBadRequest)
		return
	}

	response.Successful(w, url)
}

func DeleteLink(w http.ResponseWriter, req *http.Request) {

	uid := req.Context().Value(middleware.ContextUserKey).(models.User).Id

	// fetch from users service
	user, err := services.UsersService().GetUser(req.Context(), services.GetUserBy{Id: uid})
	if err != nil {
		response.Failed(w, models.StringError{Reason: err.Error()}, http.StatusNotFound)
		return
	}

	response.Successful(w, user)
}
package endpoints

import (
	"encoding/json"
	"net/http"

	"tomasweigenast.com/url-shortener/models"
	"tomasweigenast.com/url-shortener/server/response"
	"tomasweigenast.com/url-shortener/services"
)

func Register(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	// parse
	var body models.RegisterUser
	err := decoder.Decode(&body)
	if err != nil {
		response.Failed(w, models.StringError{Reason: "unable to parse request body"}, http.StatusBadRequest)
		return
	}

	user, err := services.UsersService().CreateUser(req.Context(), body)
	if err != nil {
		response.Failed(w, models.ParseError(err), http.StatusBadRequest)
		return
	}

	response.Successful(w, user)
}

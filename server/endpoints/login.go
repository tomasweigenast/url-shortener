package endpoints

import (
	"encoding/json"
	"net/http"

	"tomasweigenast.com/url-shortener/models"
	"tomasweigenast.com/url-shortener/server/response"
	"tomasweigenast.com/url-shortener/services"
)

func Login(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	// parse
	var body models.LoginUser
	err := decoder.Decode(&body)
	if err != nil {
		response.Failed(w, models.StringError{Reason: "unable to parse request body"}, http.StatusBadRequest)
		return
	}

	user, token, err := services.AuthService().SignInUser(req.Context(), body.Email, body.Password)
	if err != nil {
		response.Failed(w, models.StringError{Reason: err.Error()})
		return
	}

	response.Successful(w, &models.LoginResponse{
		User:        *user,
		Credentials: *token,
	})
}

package endpoints

import (
	"net/http"

	"tomasweigenast.com/url-shortener/models"
	"tomasweigenast.com/url-shortener/server/middleware"
	"tomasweigenast.com/url-shortener/server/response"
	"tomasweigenast.com/url-shortener/services"
)

func GetAccount(w http.ResponseWriter, req *http.Request) {

	uid := middleware.GetUid(req)

	// fetch from users service
	user, err := services.UsersService().GetUser(req.Context(), services.GetUserBy{Id: uid})
	if err != nil {
		response.Failed(w, models.StringError{Reason: err.Error()}, http.StatusNotFound)
		return
	}

	response.Successful(w, user)
}

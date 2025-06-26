package api

import (
	"net/http"

	"github.com/alejandro-cardenas-g/bullAndCowsApp/internal/utils"
)

func WriteJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}
	return utils.WriteJSON(w, status, &envelope{Error: message})
}

func (app *Application) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	WriteJSONError(w, http.StatusInternalServerError, "The server encountered a problem")
}

func (app *Application) BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *Application) NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	WriteJSONError(w, http.StatusNotFound, "Resource Not Found")
}

func (app *Application) ConflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	WriteJSONError(w, http.StatusConflict, err.Error())
}

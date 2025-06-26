package api

import (
	"net/http"

	"github.com/alejandro-cardenas-g/bullAndCowsApp/internal/utils"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.SugaredLogger
}

func (app *Controller) WriteJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}
	return utils.WriteJSON(w, status, &envelope{Error: message})
}

func (app *Controller) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	WriteJSONError(w, http.StatusInternalServerError, "The server encountered a problem")
}

func (app *Controller) BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *Controller) NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	WriteJSONError(w, http.StatusNotFound, "Resource Not Found")
}

func (app *Controller) ConflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	WriteJSONError(w, http.StatusConflict, err.Error())
}

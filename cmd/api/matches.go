package api

import (
	"errors"
	"net/http"

	"github.com/alejandro-cardenas-g/bullAndCowsApp/contracts"
	"github.com/alejandro-cardenas-g/bullAndCowsApp/internal/services"
	"github.com/alejandro-cardenas-g/bullAndCowsApp/internal/utils"
	"github.com/gorilla/mux"
)

type MatchesController struct {
	*Controller
	matchesService contracts.IMatchesService
}

func newMatchesController(controller *Controller, matchesService contracts.IMatchesService) *MatchesController {
	return &MatchesController{
		Controller:     controller,
		matchesService: matchesService,
	}
}

func (uc *MatchesController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/matches/create", uc.createMatchHandler).Methods("POST")
	router.HandleFunc("/matches/join/{roomId}", uc.joinMatchHandler).Methods("PUT")
	router.HandleFunc("/matches/setCombination/{roomId}", uc.setCombinationHandler).Methods("PUT")
	router.HandleFunc("/matches/startGame/{roomId}", uc.startGameHandler).Methods("PUT")
}

func (uc *MatchesController) createMatchHandler(w http.ResponseWriter, r *http.Request) {

	payload := &contracts.CreateRoomCommand{}
	if err := utils.ParseJSON(r, payload); err != nil {
		uc.BadRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		uc.BadRequestError(w, r, err)
		return
	}

	res, err := uc.matchesService.CreateRoom(r.Context(), *payload)

	if err != nil {
		uc.InternalServerError(w, r, err)
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, res); err != nil {
		uc.InternalServerError(w, r, err)
		return
	}
}

func (uc *MatchesController) joinMatchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomId := vars["roomId"]

	if err := validateRoomId(roomId); err != nil {
		uc.BadRequestError(w, r, err)
		return
	}

	payload := &contracts.JoinRoomCommand{}
	if err := utils.ParseJSON(r, payload); err != nil {
		uc.BadRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		uc.BadRequestError(w, r, err)
		return
	}

	payload.RoomId = roomId

	res, err := uc.matchesService.JoinRoom(r.Context(), *payload)

	if err != nil {
		switch err {
		case services.ErrMatchNotFound:
			{
				uc.NotFoundError(w, r, err)
			}
		case services.ErrCanNotAddAnotherPlayer:
			{
				uc.ConflictError(w, r, err)
			}
		default:
			{
				uc.InternalServerError(w, r, err)
			}
		}
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, res); err != nil {
		uc.InternalServerError(w, r, err)
		return
	}
}

func (uc *MatchesController) setCombinationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomId := vars["roomId"]

	if err := validateRoomId(roomId); err != nil {
		uc.BadRequestError(w, r, err)
		return
	}

	payload := &contracts.SetCombinationCommand{}
	if err := utils.ParseJSON(r, payload); err != nil {
		uc.BadRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		uc.BadRequestError(w, r, err)
		return
	}

	payload.RoomId = roomId

	res, err := uc.matchesService.SetCombination(r.Context(), *payload)

	if err != nil {
		switch {
		case errors.Is(err, services.ErrMatchNotFullRoom):
			{
				uc.ConflictError(w, r, err)
			}
		case errors.Is(err, services.ErrInvalidCombination):
			{
				uc.BadRequestError(w, r, err)
			}
		case errors.Is(err, services.ErrMatchNotFound):
			{
				uc.NotFoundError(w, r, err)
			}
		default:
			{
				uc.InternalServerError(w, r, err)
			}
		}
		return
	}

	if err := utils.WriteJSON(w, http.StatusAccepted, res); err != nil {
		uc.InternalServerError(w, r, err)
		return
	}
}

func (uc *MatchesController) startGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomId := vars["roomId"]

	if err := validateRoomId(roomId); err != nil {
		uc.BadRequestError(w, r, err)
		return
	}

	result, err := uc.matchesService.StartGame(r.Context(), roomId)

	if err != nil {
		switch {
		case errors.Is(err, services.ErrMatchNotFullRoom):
			{
				uc.ConflictError(w, r, err)
			}
		case errors.Is(err, services.ErrMatchNotFound):
			{
				uc.NotFoundError(w, r, err)
			}
		default:
			{
				uc.InternalServerError(w, r, err)
			}
		}
		return
	}

	if err := utils.WriteJSON(w, http.StatusAccepted, result); err != nil {
		uc.InternalServerError(w, r, err)
		return
	}
}

func validateRoomId(roomId string) error {
	return Validate.Struct(struct {
		RoomId string `validate:"required,len=7"`
	}{RoomId: roomId})
}

package services

import (
	"at/constants"
	"at/tools"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

type UserController struct {
	userRepo UserRepository
}

func NewUserController(repo UserRepository) *UserController {
	return &UserController{userRepo: repo}
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request) {

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrDecodingRequestBody)
		http.Error(w, constants.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	if err := uc.userRepo.CreateUser(&user); err != nil {
		log.Error().Err(err).Msg(constants.ErrUserCreate)
		http.Error(w, constants.ErrUserCreate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrEncodingResponse)
		return
	}

	log.Info().Msg(constants.SuccessfullyUserCreate)
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {

	id, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrExtractId)
		http.Error(w, constants.ErrExtractId, http.StatusNotFound)
		return
	}

	user, err := uc.userRepo.GetUser(uint64(id))
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrUserFetching)
		http.Error(w, constants.ErrUserFetching, http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrEncodingResponse)
		return
	}

	log.Info().Msg(constants.CallControllerGetUserById)
}

package services

import (
	"at/tools"
	"at/tools/errors"
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
		log.Warn().Err(err).Msg(errors.ErrDecodingRequestBody)
		http.Error(w, errors.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	if err := uc.userRepo.CreateUser(&user); err != nil {
		log.Warn().Err(err).Msg(errors.ErrUserCreate)
		http.Error(w, errors.ErrUserCreate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrEncodingResponse)
		return
	}
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {

	id, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrExtractId)
		http.Error(w, errors.ErrExtractId, http.StatusNotFound)
		return
	}

	user, err := uc.userRepo.GetUser(uint64(id))
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrUserFetching)
		http.Error(w, errors.ErrUserFetching, http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrEncodingResponse)
		return
	}
}

func (uc *UserController) BanUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		TelegramID uint64 `json:"telegram_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrDecodingRequestBody)
		http.Error(w, errors.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	if payload.TelegramID == 0 {
		log.Warn().Msg(errors.ErrUserBanned)
		http.Error(w, errors.ErrUserBanned, http.StatusBadRequest)
		return
	}

	err = uc.userRepo.BanUser(payload.TelegramID)
	if err != nil {
		if err.Error() == errors.ErrUserNotFound {
			log.Warn().Err(err).Msg(errors.ErrUserNotFound)
			http.Error(w, errors.ErrUserNotFound, http.StatusNotFound)
			return
		}

		log.Warn().Err(err).Msg(errors.ErrUserBanned)
		http.Error(w, errors.ErrUserBanned, http.StatusInternalServerError)
		return
	}
}

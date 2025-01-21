package subscribe

import (
	"at/constants"
	"at/tools"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

type SubscribeController struct {
	subRepo SubscribeRepository
}

func NewSubscribeController(repo SubscribeRepository) *SubscribeController {
	return &SubscribeController{subRepo: repo}
}

func (sc *SubscribeController) AddSubscribe(w http.ResponseWriter, r *http.Request) {
	var sub Subscribe
	err := json.NewDecoder(r.Body).Decode(&sub)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrDecodingRequestBody)
		http.Error(w, constants.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	if err := sc.subRepo.AddSubscribe(&sub); err != nil {
		log.Error().Err(err).Msgf("%v, %v", constants.ErrSubscribeCreate, sub)
		http.Error(w, constants.ErrSubscribeCreate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (sc *SubscribeController) GetSubscribesByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := tools.ExtractID(r.URL.Path)
	if err != nil || userID <= 0 {
		log.Error().Msg(constants.ErrExtractId)
		http.Error(w, constants.ErrExtractId, http.StatusBadRequest)
		return
	}

	subs, err := sc.subRepo.GetSubscribesByUser(uint(userID))
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrSubscribesByUser)
		http.Error(w, constants.ErrSubscribesByUser, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(subs)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrEncodingResponse)
		http.Error(w, constants.ErrEncodingResponse, http.StatusInternalServerError)
		return
	}
}

func (sc *SubscribeController) GetSubscribesBySponsor(w http.ResponseWriter, r *http.Request) {
	sponsorID, err := tools.ExtractID(r.URL.Path)
	if err != nil || sponsorID <= 0 {
		log.Error().Msg(constants.ErrExtractId)
		http.Error(w, constants.ErrExtractId, http.StatusBadRequest)
		return
	}

	subs, err := sc.subRepo.GetSubscribesBySponsor(uint(sponsorID))
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrSubscribesBySponsor)
		http.Error(w, constants.ErrSubscribesBySponsor, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(subs)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrEncodingResponse)
		http.Error(w, constants.ErrEncodingResponse, http.StatusInternalServerError)
		return
	}
}

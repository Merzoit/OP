package subscribe

import (
	"at/tools"
	"at/tools/errors"
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
		log.Warn().Err(err).Msg(errors.ErrDecodingRequestBody)
		http.Error(w, errors.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	if sub, err := sc.subRepo.AddSubscribe(&sub); err != nil {
		log.Warn().Err(err).Msgf("%v, %v", errors.ErrSubscribeCreate, sub)
		http.Error(w, errors.ErrSubscribeCreate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (sc *SubscribeController) GetSubscribesByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := tools.ExtractID(r.URL.Path)
	if err != nil || userID <= 0 {
		log.Warn().Msg(errors.ErrExtractId)
		http.Error(w, errors.ErrExtractId, http.StatusBadRequest)
		return
	}

	subs, err := sc.subRepo.GetSubscribesByUser(uint(userID))
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrSubscribesByUser)
		http.Error(w, errors.ErrSubscribesByUser, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(subs)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrEncodingResponse)
		http.Error(w, errors.ErrEncodingResponse, http.StatusInternalServerError)
		return
	}
}

func (sc *SubscribeController) GetSubscribesBySponsor(w http.ResponseWriter, r *http.Request) {
	sponsorID, err := tools.ExtractID(r.URL.Path)
	if err != nil || sponsorID <= 0 {
		log.Warn().Msg(errors.ErrExtractId)
		http.Error(w, errors.ErrExtractId, http.StatusBadRequest)
		return
	}

	subs, err := sc.subRepo.GetSubscribesBySponsor(uint(sponsorID))
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrSubscribesBySponsor)
		http.Error(w, errors.ErrSubscribesBySponsor, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(subs)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrEncodingResponse)
		http.Error(w, errors.ErrEncodingResponse, http.StatusInternalServerError)
		return
	}
}

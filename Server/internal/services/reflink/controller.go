package reflink

import (
	"at/tools"
	"at/tools/errors"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

type RefLinkController struct {
	refRepo ReflinkRepository
}

func NewRefLinkController(refRepo ReflinkRepository) *RefLinkController {
	return &RefLinkController{refRepo: refRepo}
}

func (rc RefLinkController) CreateLink(w http.ResponseWriter, r *http.Request) {
	var link ReferralLink

	err := json.NewDecoder(r.Body).Decode(&link)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrDecodingRequestBody)
		http.Error(w, errors.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	l, err := rc.refRepo.CreateLink(&link)
	if err != nil {
		log.Warn().Err(err).Msgf("%v, %v", errors.ErrRefCreate, l)
		http.Error(w, errors.ErrRefCreate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(link)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrEncodingResponse)
		http.Error(w, errors.ErrEncodingResponse, http.StatusInternalServerError)
		return
	}
}

func (rc RefLinkController) GetLink(w http.ResponseWriter, r *http.Request) {
	workerID, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrExtractId)
		http.Error(w, errors.ErrExtractId, http.StatusBadRequest)
		return
	}

	link, err := rc.refRepo.GetLink(int(workerID))
	if err != nil {
		if err.Error() == errors.ErrNoRows {
			log.Warn().Err(err).Msg(errors.ErrWorkerNotFound)
			http.Error(w, errors.ErrWorkerNotFound, http.StatusNotFound)
			return
		}

		log.Warn().Err(err).Msg(errors.ErrRefGet)
		http.Error(w, errors.ErrRefGet, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(link)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrEncodingResponse)
		http.Error(w, errors.ErrEncodingResponse, http.StatusInternalServerError)
		return
	}
}

func (rc RefLinkController) UpdateLink(w http.ResponseWriter, r *http.Request) {
	workerID, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrExtractId)
		http.Error(w, errors.ErrExtractId, http.StatusBadRequest)
		return
	}

	var payload struct {
		Tag string `json:"tag"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrDecodingRequestBody)
		http.Error(w, errors.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	err = rc.refRepo.UpdateLink(int(workerID), payload.Tag)
	if err != nil {
		if err.Error() == errors.ErrRefNotFound {
			log.Warn().Err(err).Msg(errors.ErrRefNotFound)
			http.Error(w, errors.ErrRefNotFound, http.StatusNotFound)
			return
		}

		log.Warn().Err(err).Msg(errors.ErrRefUpdate)
		http.Error(w, errors.ErrRefUpdate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc RefLinkController) ClickAdd(w http.ResponseWriter, r *http.Request) {
	workerID, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrExtractId)
		http.Error(w, errors.ErrExtractId, http.StatusBadRequest)
		return
	}

	var payload struct {
		Count int `json:"count"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrDecodingRequestBody)
		http.Error(w, errors.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	err = rc.refRepo.ClickAdd(int(workerID), payload.Count)
	if err != nil {
		if err.Error() == errors.ErrRefNotFound {
			log.Warn().Err(err).Msg(errors.ErrRefNotFound)
			http.Error(w, errors.ErrRefNotFound, http.StatusNotFound)
			return
		}

		log.Warn().Err(err).Msg(errors.ErrRefClickAdd)
		http.Error(w, errors.ErrRefClickAdd, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc RefLinkController) RegistrationAdd(w http.ResponseWriter, r *http.Request) {
	workerID, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrExtractId)
		http.Error(w, errors.ErrExtractId, http.StatusBadRequest)
		return
	}

	var payload struct {
		Count int `json:"count"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrDecodingRequestBody)
		http.Error(w, errors.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	err = rc.refRepo.RegistrationAdd(int(workerID), payload.Count)
	if err != nil {
		if err.Error() == errors.ErrRefNotFound {
			log.Warn().Err(err).Msg(errors.ErrRefNotFound)
			http.Error(w, errors.ErrRefNotFound, http.StatusNotFound)
			return
		}

		log.Warn().Err(err).Msg(errors.ErrRefRegAdd)
		http.Error(w, errors.ErrRefRegAdd, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

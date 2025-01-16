package reflink

import (
	"at/constants"
	"at/tools"
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
		log.Error().Err(err).Msg(constants.ErrDecodingRequestBody)
		http.Error(w, constants.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	err = rc.refRepo.CreateLink(&link)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrRefCreate)
		http.Error(w, constants.ErrRefCreate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(link)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrEncodingResponse)
		http.Error(w, constants.ErrEncodingResponse, http.StatusInternalServerError)
		return
	}
}

func (rc RefLinkController) GetLink(w http.ResponseWriter, r *http.Request) {
	workerID, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrExtractId)
		http.Error(w, constants.ErrExtractId, http.StatusBadRequest)
		return
	}

	link, err := rc.refRepo.GetLink(int(workerID))
	if err != nil {
		if err.Error() == constants.ErrNoRows {
			log.Error().Err(err).Msg(constants.ErrWorkerNotFound)
			http.Error(w, constants.ErrWorkerNotFound, http.StatusNotFound)
			return
		}

		log.Error().Err(err).Msg(constants.ErrRefGet)
		http.Error(w, constants.ErrRefGet, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(link)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrEncodingResponse)
		http.Error(w, constants.ErrEncodingResponse, http.StatusInternalServerError)
		return
	}
}

func (rc RefLinkController) UpdateLink(w http.ResponseWriter, r *http.Request) {
	workerID, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrExtractId)
		http.Error(w, constants.ErrExtractId, http.StatusBadRequest)
		return
	}

	var payload struct {
		Tag string `json:"tag"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrDecodingRequestBody)
		http.Error(w, constants.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	err = rc.refRepo.UpdateLink(int(workerID), payload.Tag)
	if err != nil {
		if err.Error() == constants.ErrRefNotFound {
			log.Error().Err(err).Msg(constants.ErrRefNotFound)
			http.Error(w, constants.ErrRefNotFound, http.StatusNotFound)
			return
		}

		log.Error().Err(err).Msg(constants.ErrRefUpdate)
		http.Error(w, constants.ErrRefUpdate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc RefLinkController) ClickAdd(w http.ResponseWriter, r *http.Request) {
	workerID, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrExtractId)
		http.Error(w, constants.ErrExtractId, http.StatusBadRequest)
		return
	}

	var payload struct {
		Count int `json:"count"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrDecodingRequestBody)
		http.Error(w, constants.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	err = rc.refRepo.ClickAdd(int(workerID), payload.Count)
	if err != nil {
		if err.Error() == constants.ErrRefNotFound {
			log.Error().Err(err).Msg(constants.ErrRefNotFound)
			http.Error(w, constants.ErrRefNotFound, http.StatusNotFound)
			return
		}

		log.Error().Err(err).Msg(constants.ErrRefClickAdd)
		http.Error(w, constants.ErrRefClickAdd, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc RefLinkController) RegistrationAdd(w http.ResponseWriter, r *http.Request) {
	workerID, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrExtractId)
		http.Error(w, constants.ErrExtractId, http.StatusBadRequest)
		return
	}

	var payload struct {
		Count int `json:"count"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrDecodingRequestBody)
		http.Error(w, constants.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	err = rc.refRepo.RegistrationAdd(int(workerID), payload.Count)
	if err != nil {
		if err.Error() == constants.ErrRefNotFound {
			log.Error().Err(err).Msg(constants.ErrRefNotFound)
			http.Error(w, constants.ErrRefNotFound, http.StatusNotFound)
			return
		}

		log.Error().Err(err).Msg(constants.ErrRefRegAdd)
		http.Error(w, constants.ErrRefRegAdd, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

package code

import (
	"at/tools"
	"at/tools/errors"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

type CodeController struct {
	codeRepo CodeRepository
}

func NewCodeController(codeRepo CodeRepository) *CodeController {
	return &CodeController{codeRepo: codeRepo}
}

func (cc CodeController) CreateCode(w http.ResponseWriter, r *http.Request) {

	var code Code
	err := json.NewDecoder(r.Body).Decode(&code)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrDecodingRequestBody)
		http.Error(w, errors.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	if c, err := cc.codeRepo.CreateCode(&code); err != nil {
		log.Warn().Err(err).Msgf("%v, %v", errors.ErrCodeCreate, c)
		http.Error(w, errors.ErrCodeCreate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(code)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrEncodingResponse)
		return
	}
}

func (cc CodeController) GetCode(w http.ResponseWriter, r *http.Request) {

	id, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrExtractId)
		http.Error(w, errors.ErrExtractId, http.StatusBadRequest)
		return
	}

	code, err := cc.codeRepo.GetCode(int(id))
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrCodeFetching)
		http.Error(w, errors.ErrCodeFetching, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(code)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrEncodingResponse)
		return
	}
}

func (cc CodeController) DeleteCode(w http.ResponseWriter, r *http.Request) {

	id, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrExtractId)
		http.Error(w, errors.ErrExtractId, http.StatusBadRequest)
		return
	}

	err = cc.codeRepo.DeleteCode(int(id))
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrCodeDelete)
		http.Error(w, errors.ErrCodeDelete, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cc CodeController) AddRequestCount(w http.ResponseWriter, r *http.Request) {

	accessCode, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrExtractId)
		http.Error(w, errors.ErrExtractId, http.StatusBadRequest)
		return
	}

	var payload struct {
		Increment int `json:"increment"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrDecodingRequestBody)
		http.Error(w, errors.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	err = cc.codeRepo.AddRequestCount(int(accessCode), payload.Increment)
	if err != nil {
		if err.Error() == errors.ErrCodeNotFound {
			log.Warn().Err(err).Msg(errors.ErrCodeNotFound)
			http.Error(w, errors.ErrCodeNotFound, http.StatusNotFound)
			return
		}

		log.Warn().Err(err).Msg(errors.ErrCodeRequestCounter)
		http.Error(w, errors.ErrCodeRequestCounter, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cc CodeController) GetCodesByWorker(w http.ResponseWriter, r *http.Request) {

	workerID, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrExtractId)
		http.Error(w, errors.ErrExtractId, http.StatusBadRequest)
		return
	}

	codes, err := cc.codeRepo.GetCodesByWorker(uint(workerID))
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrCodeFetchingByWorker)
		http.Error(w, errors.ErrCodeFetchingByWorker, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(codes)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrEncodingResponse)
		http.Error(w, errors.ErrEncodingResponse, http.StatusInternalServerError)
		return
	}
}

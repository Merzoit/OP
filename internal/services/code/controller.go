package code

import (
	"at/constants"
	"at/tools"
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
		log.Error().Err(err).Msg(constants.ErrDecodingRequestBody)
		http.Error(w, constants.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	if err := cc.codeRepo.CreateCode(&code); err != nil {
		log.Error().Err(err).Msg(constants.ErrCodeCreate)
		http.Error(w, constants.ErrCodeCreate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(code)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrEncodingResponse)
		return
	}

	log.Info().Msg(constants.SuccessfullyCodeCreate)
}

func (cc CodeController) GetCode(w http.ResponseWriter, r *http.Request) {

	id, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrExtractId)
		http.Error(w, constants.ErrExtractId, http.StatusBadRequest)
		return
	}

	code, err := cc.codeRepo.GetCode(int(id))
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrCodeFetching)
		http.Error(w, constants.ErrCodeFetching, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(code)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrEncodingResponse)
		return
	}

	log.Info().Msg(constants.SuccessfullyControllerGetCode)
}

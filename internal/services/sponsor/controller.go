package sponsor

import (
	"at/constants"
	"at/tools"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

type SponsorController struct {
	sponsorRepo SponsorRepository
}

func NewSponsorController(repo SponsorRepository) *SponsorController {
	return &SponsorController{sponsorRepo: repo}
}

func (sc SponsorController) CreateSponsor(w http.ResponseWriter, r *http.Request) {
	var sponsor Sponsor
	err := json.NewDecoder(r.Body).Decode(&sponsor)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrDecodingRequestBody)
		http.Error(w, constants.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	if err := sc.sponsorRepo.CreateSponsor(&sponsor); err != nil {
		log.Error().Err(err).Msg(constants.ErrSponsorCreate)
		http.Error(w, constants.ErrSponsorCreate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(sponsor)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrEncodingResponse)
		return
	}

	log.Info().Msg(constants.SuccessfullySponsorCreate)
}

func (sc SponsorController) DeleteSponsor(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg(constants.CallControllerSponsorDelete)

	id, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrExtractId)
		http.Error(w, constants.ErrExtractId, http.StatusBadRequest)
		return
	}

	if err := sc.sponsorRepo.DeleteSponsor(uint64(id)); err != nil {
		log.Error().Err(err).Msg(constants.ErrSponsorDelete)
		http.Error(w, constants.ErrSponsorDelete, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Info().Msg(constants.SuccessfullyControllerSponsorDelete)
}

func (sc SponsorController) GetSponsor(w http.ResponseWriter, r *http.Request) {
	id, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrExtractId)
		http.Error(w, constants.ErrExtractId, http.StatusBadRequest)
		return
	}

	sponsor, err := sc.sponsorRepo.GetSponsor(uint64(id))
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrSponsorFetching)
		http.Error(w, constants.ErrSponsorFetching, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(sponsor)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrEncodingResponse)
		return
	}

	log.Info().Msgf(constants.SuccessfullySponsorDelete, id)
}

func (sc SponsorController) GetSponsors(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg(constants.CallControllerSponsorsGet)

	sponsors, err := sc.sponsorRepo.GetSponsors()
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrSponsorFetching)
		http.Error(w, constants.ErrSponsorFetching, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(sponsors)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrEncodingResponse)
		return
	}

	log.Info().Msgf(constants.SuccessfullySponsorsGet, len(sponsors))
}

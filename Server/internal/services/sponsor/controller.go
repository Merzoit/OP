package sponsor

import (
	"at/tools"
	"at/tools/errors"
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
		log.Warn().Err(err).Msg(errors.ErrDecodingRequestBody)
		http.Error(w, errors.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	if s, err := sc.sponsorRepo.CreateSponsor(&sponsor); err != nil {
		log.Warn().Err(err).Msgf("%v, %v", errors.ErrSponsorCreate, s)
		http.Error(w, errors.ErrSponsorCreate, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(sponsor)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrEncodingResponse)
		return
	}
}

func (sc SponsorController) DeleteSponsor(w http.ResponseWriter, r *http.Request) {
	id, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrExtractId)
		http.Error(w, errors.ErrExtractId, http.StatusBadRequest)
		return
	}

	if err := sc.sponsorRepo.DeleteSponsor(uint64(id)); err != nil {
		log.Warn().Err(err).Msg(errors.ErrSponsorDelete)
		http.Error(w, errors.ErrSponsorDelete, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (sc SponsorController) GetSponsor(w http.ResponseWriter, r *http.Request) {
	id, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrExtractId)
		http.Error(w, errors.ErrExtractId, http.StatusBadRequest)
		return
	}

	sponsor, err := sc.sponsorRepo.GetSponsor(uint64(id))
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrSponsorFetching)
		http.Error(w, errors.ErrSponsorFetching, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(sponsor)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrEncodingResponse)
		return
	}
}

func (sc SponsorController) GetSponsors(w http.ResponseWriter, r *http.Request) {
	sponsors, err := sc.sponsorRepo.GetSponsors()
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrSponsorFetching)
		http.Error(w, errors.ErrSponsorFetching, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(sponsors)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrEncodingResponse)
		return
	}
}

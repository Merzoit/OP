package role

import (
	"at/constants"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

type RoleController struct {
	roleRepo RoleRepository
}

func NewRoleController(repo RoleRepository) *RoleController {
	return &RoleController{roleRepo: repo}
}

func (rc *RoleController) GetRole(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RoleID uint `json:"role_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrDecodingRequestBody)
		http.Error(w, constants.ErrDecodingRequestBody, http.StatusBadRequest)
		return
	}

	if payload.RoleID == 0 {
		log.Error().Msg(constants.ErrExtractId)
		http.Error(w, constants.ErrExtractId, http.StatusBadRequest)
		return
	}

	role, err := rc.roleRepo.GetRole(payload.RoleID)
	if err != nil {
		if err.Error() == constants.ErrRoleNotFound {
			log.Error().Err(err).Msg(constants.ErrRoleNotFound)
			http.Error(w, constants.ErrRoleNotFound, http.StatusNotFound)
			return
		}

		log.Error().Err(err).Msg(constants.ErrRoleFetching)
		http.Error(w, constants.ErrRoleFetching, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(role)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrEncodingResponse)
		http.Error(w, constants.ErrEncodingResponse, http.StatusInternalServerError)
		return
	}
}

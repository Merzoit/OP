package role

import (
	"at/constants"
	"at/tools"
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
	id, err := tools.ExtractID(r.URL.Path)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrExtractId)
		http.Error(w, constants.ErrExtractId, http.StatusBadRequest)
		return
	}

	role, err := rc.roleRepo.GetRole(uint(id))
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrRoleFetching)
		http.Error(w, constants.ErrRoleFetching, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(role)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrEncodingResponse)
		return
	}
}

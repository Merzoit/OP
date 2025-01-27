package role

import (
	"at/tools"
	"at/tools/errors"
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
		log.Warn().Err(err).Msg(errors.ErrExtractId)
		http.Error(w, errors.ErrExtractId, http.StatusBadRequest)
		return
	}

	role, err := rc.roleRepo.GetRole(uint(id))
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrRoleFetching)
		http.Error(w, errors.ErrRoleFetching, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(role)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrEncodingResponse)
		return
	}
}

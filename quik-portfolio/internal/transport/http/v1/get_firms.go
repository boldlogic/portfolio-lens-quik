package v1

import (
	"net/http"

	"github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"
)

func (h *Handler) GetFirms(r *http.Request) (any, string, error) {
	ctx := r.Context()

	firms, err := h.service.GetFirms(ctx)
	if err != nil {
		return nil, "", err
	}
	return firmsToResp(firms), "", nil
}

func firmsToResp(firms []quik.Firm) []firmRespDTO {
	if len(firms) == 0 {
		return []firmRespDTO{}
	}

	resp := make([]firmRespDTO, 0, len(firms))
	for _, f := range firms {
		resp = append(resp, firmToResp(f))
	}
	return resp
}

func firmToResp(f quik.Firm) firmRespDTO {
	return firmRespDTO{
		Id:   f.Id,
		Code: f.Code,
		Name: f.Name,
	}
}

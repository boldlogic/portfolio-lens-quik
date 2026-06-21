package v1

import "net/http"

func (h *Handler) GetFirms(r *http.Request) (any, string, error) {
	firms, err := h.readService.GetFirms(r.Context())
	if err != nil {
		return nil, "", err
	}
	return firmsToResp(firms), "", nil
}

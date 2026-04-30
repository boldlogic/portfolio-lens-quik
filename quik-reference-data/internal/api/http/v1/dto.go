package v1

import "github.com/boldlogic/portfolio-lens-quik/pkg/models/quik"

type firmCreateReqDTO struct {
	Code string `json:"firmCode" validate:"required,min=1,max=12"`
	Name string `json:"firmName" validate:"required,min=1,max=128"`
}

type firmPatchReqDTO struct {
	Name string `json:"firmName" validate:"required,min=1,max=128"`
}

type firmRespDTO struct {
	Id   uint8  `json:"id"`
	Code string `json:"firmCode"`
	Name string `json:"firmName"`
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

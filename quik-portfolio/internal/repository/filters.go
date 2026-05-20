package repository

import mssql "github.com/microsoft/go-mssqldb"

type clientRows struct {
	ClientCode string `tvp:"client_code"`
}

func (r *Repository) makeClientCodeList(clientCodes []string) (mssql.TVP, bool) {
	if len(clientCodes) == 0 {
		return mssql.TVP{}, false
	}

	clients := make([]clientRows, 0, len(clientCodes))
	for _, code := range clientCodes {
		clients = append(clients, clientRows{ClientCode: code})
	}
	return mssql.TVP{
		TypeName: "api.client_code_list",
		Value:    clients,
	}, true
}

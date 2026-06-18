package dbrepo

import mssql "github.com/microsoft/go-mssqldb"

func MakeTVP[T, R any](rows []T, mapFn func(T) R, tvpName string) (mssql.TVP, bool) {
	if len(rows) == 0 {
		return mssql.TVP{}, false
	}
	res := make([]R, 0, len(rows))
	for _, v := range rows {
		res = append(res, mapFn(v))
	}
	return mssql.TVP{
		TypeName: tvpName,
		Value:    res,
	}, true

}

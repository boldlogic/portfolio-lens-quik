package models

import "errors"

var (
	ErrNotFound           = errors.New("данные по запросу не найдены")
	ErrRetrievingData     = errors.New("ошибка при получении данных")
	ErrSavingData         = errors.New("ошибка при изменении данных")
	ErrValidation         = errors.New("некорректные входные данные")
	ErrBusinessValidation = errors.New("некорректные данные в запросе")
	ErrConflict           = errors.New("запись с таким ключом уже существует")
)

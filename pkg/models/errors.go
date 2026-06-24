package models

import "errors"

var (
	ErrUnauthorized       = errors.New("недействительный ключ")
	ErrNotFound           = errors.New("данные по запросу не найдены")
	ErrRetrievingData     = errors.New("ошибка при получении данных")
	ErrSavingData         = errors.New("ошибка при изменении данных")
	ErrPartialSuccess     = errors.New("частичная ошибка")
	ErrValidation         = errors.New("некорректные входные данные")
	ErrBusinessValidation = errors.New("некорректные данные") //некорректные данные
	ErrConflict           = errors.New("запись с таким ключом уже существует")
)

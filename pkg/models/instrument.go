package models

// type Instrument struct {
// 	Ticker             string `gorm:"column:ticker;type:char(15);not null"`      // Код инструмента
// 	RegistrationNumber string `gorm:"column:registration_number;type:char(250)"` // Рег.номер инструмента
// 	FullName           string `gorm:"column:full_name;type:char(250);not null"`  // Полное название инструмента
// 	ShortName          string `gorm:"column:short_name;type:char(100)"`          // Краткое название
// 	ClassCode          string `gorm:"column:class_code;type:char(20)"`           // Код класса
// 	ClassName          string `gorm:"column:class_name;type:char(200)"`          // Наименование класса

// 	//InstrumentType     string     `gorm:"column:instrument_type;type:char(100)"`     // Тип инструмента
// 	//InstrumentSubtype  string     `gorm:"column:instrument_subtype;type:char(100)"`  // Подтип инструмента
// 	AssetClass      string
// 	AssetSubClass   string
// 	ISIN            string     `gorm:"column:isin;type:char(15)"`            // Международный идентификатор
// 	FaceValue       float64    `gorm:"column:face_value;type:float"`         // Номинал
// 	BaseCurrency    string     `gorm:"column:base_currency;type:char(3)"`    // Валюта номинала / базовая валюта
// 	QuoteCurrency   string     `gorm:"column:quote_currency;type:char(3)"`   // Валюта котировки
// 	CounterCurrency string     `gorm:"column:counter_currency;type:char(3)"` // Сопряженная валюта
// 	MaturityDate    *time.Time `gorm:"column:maturity_date;type:date"`       // Дата погашения
// 	CouponDuration  int        `gorm:"column:coupon_duration;type:int"`      // Длительность купона
// }

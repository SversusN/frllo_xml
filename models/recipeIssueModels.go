package models

import "encoding/xml"

// RootRecipe структура для корневого элемента XML
type RootRecipe struct {
	XMLName     xml.Name         `xml:"root"`
	InfoSysCode string           `xml:"info_sys_code,omitempty"`      // Исключаем, если пустая строка
	Documents   []DocumentRecipe `xml:"documents>document,omitempty"` // Исключаем, если массив пуст
}

// DocumentRecipe структура для документа
type DocumentRecipe struct {
	DocumentID  string        `xml:"document_id,omitempty" db:"document_id"`     // Исключаем, если пустая строка
	DocDateTime string        `xml:"doc_date_time,omitempty" db:"doc_date_time"` // Исключаем, если пустая строка
	Citizen     CitizenRecipe `xml:"citizen,omitempty"`                          // Исключаем, если Citizen пустая структура
	Recipe      Recipe        `xml:"recipe,omitempty"`                           // Исключаем, если Recipe пустая структура
}

// CitizenRecipe структура для гражданина
type CitizenRecipe struct {
	RegisterID   string      `xml:"register_id,omitempty"`                        // Исключаем, если пустая строка
	ExtCitizenID string      `xml:"ext_citizen_id,omitempty" db:"ext_citizen_id"` // Исключаем, если пустая строка
	Name         string      `xml:"name,omitempty" db:"name"`                     // Исключаем, если пустая строка
	Surname      string      `xml:"surname,omitempty" db:"surname"`               // Исключаем, если пустая строка
	Patronymic   string      `xml:"patronymic,omitempty" db:"patronymic"`         // Исключаем, если пустая строка
	Birthdate    string      `xml:"birthdate,omitempty" db:"birthdate"`           // Исключаем, если пустая строка
	Sex          string      `xml:"sex,omitempty" db:"sex"`                       // Исключаем, если пустая строка
	Region       string      `xml:"region,omitempty"`                             // Исключаем, если пустая строка
	Snils        string      `xml:"snils,omitempty" db:"snils"`                   // Исключаем, если пустая строка
	IdentifyDocs []DocRecipe `xml:"identify_docs>doc,omitempty"`                  // Исключаем, если массив пуст
}

// DocRecipe структура для документов удостоверяющих личность
type DocRecipe struct {
	DocType   string `xml:"doc_type,omitempty" db:"doc_type"`     // Исключаем, если пустая строка
	Serial    string `xml:"serial,omitempty" db:"serial"`         // Исключаем, если пустая строка
	Num       string `xml:"num,omitempty" db:"num"`               // Исключаем, если пустая строка
	DateIssue string `xml:"date_issue,omitempty" db:"date_issue"` // Исключаем, если пустая строка
	Authority string `xml:"authority,omitempty" db:"authority"`   // Исключаем, если пустая строка
}

// Recipe структура для рецепта
type Recipe struct {
	RecipeSerial      string `xml:"recipe_serial,omitempty" db:"recipe_serial"`             // Исключаем, если пустая строка
	RecipeNum         string `xml:"recipe_num,omitempty" db:"recipe_num"`                   // Исключаем, если пустая строка
	ExtRecipeID       string `xml:"ext_recipe_id,omitempty" db:"ext_recipe_id"`             // Исключаем, если пустая строка
	MedOrgOID         string `xml:"med_org_oid,omitempty" db:"med_org_oid"`                 // Исключаем, если пустая строка
	DoctorName        string `xml:"doctor_name,omitempty" db:"doctor_name"`                 // Исключаем, если пустая строка
	StaffPositionCode string `xml:"staff_position_code,omitempty" db:"staff_position_code"` // Исключаем, если пустая строка
	DoctorSnils       string `xml:"doctor_snils,omitempty" db:"doctor_snils"`               // Исключаем, если пустая строка
	MedicalCard       string `xml:"medical_card,omitempty" db:"medical_card"`               // Исключаем, если пустая строка
	BenefitCode       string `xml:"benefit_code,omitempty" db:"benefit_code"`               // Исключаем, если пустая строка
	MKB10Code         string `xml:"mkb10_code,omitempty" db:"mkb10_code"`                   // Исключаем, если пустая строка
	DrugSmnnCode      string `xml:"drug_smnn_code,omitempty" db:"drug_smnn_code"`           // Исключаем, если пустая строка
	CommissionDate    string `xml:"commission_date,omitempty" db:"commission_date"`         // Исключаем, если пустая строка
	CommissionNum     string `xml:"commission_num,omitempty" db:"commission_num"`           // Исключаем, если пустая строка
	Qty               int    `xml:"qty,omitempty" db:"qty"`                                 // Исключаем, если значение равно 0
	RecipeDate        string `xml:"recipe_date,omitempty" db:"recipe_date"`                 // Исключаем, если пустая строка
	RecipeExpiryCode  int    `xml:"recipe_expiry_code,omitempty" db:"recipe_expiry_code"`   // Исключаем, если значение равно 0
	DateExpiry        string `xml:"date_expiry,omitempty" db:"date_expiry"`                 // Исключаем, если пустая строка
	TS                int    `db:"ts"`
}

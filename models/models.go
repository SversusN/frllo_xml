package models

import (
	"database/sql"
	"encoding/xml"
)

type Document struct {
	DocumentID  string       `xml:"document_id" db:"Id"`
	DocDateTime string       `xml:"doc_date_time"`
	Citizen     Citizen      `xml:"citizen"`
	Benefits    []BenefitXML `xml:"benefits>benefit"`
}

// Citizen структура для гражданина
type Citizen struct {
	RegisterID   string `xml:"register_id,omitempty" db:"FrloId"`
	ExtCitizenID string `xml:"ext_citizen_id,omitempty" db:"Id"`
	Name         string `xml:"name,omitempty" db:"IndividualFirstName"`
	Surname      string `xml:"surname,omitempty" db:"IndividualLastName"`
	Patronymic   string `xml:"patronymic,omitempty" db:"IndividualPatronymic"`
	Birthdate    string `xml:"birthdate,omitempty"`
	Sex          string `xml:"sex,omitempty"`
	Citizenship  string `xml:"citizenship,omitempty"`
	Snils        string `xml:"snils,omitempty"`
	PolicySN     string `xml:"policy_sn,omitempty" db:"IndividualPolicy"`
	Region       string `xml:"region,omitempty"`
	IdentifyDocs []Doc  `xml:"identify_docs>doc,omitempty"`
}

// Doc структура для документов удостоверяющих личность

type Doc struct {
	DocTypeCIT     int    `xml:"doc_type,omitempty"`
	DocTypeNameCIT string `xml:"doc_type_name_cit,omitempty"`
	SerialCIT      string `xml:"serial,omitempty" db:"credential_serie"`
	NumCIT         string `xml:"num,omitempty" db:"credential_number"`
	DateIssueCIT   string `xml:"date_issue,omitempty" db:"date_issue_cit"`
	DateExpiryCIT  string `xml:"date_expiry,omitempty"`
	AuthorityCIT   string `xml:"authority,omitempty" db:"CitizenshipAuthority"`
	SerialIDEN     string `xml:"serial_iden,omitempty" db:"seria_iden,omitempty"`
	NumIDEN        string `xml:"num_iden,omitempty" db:"num_iden,omitempty"`
	DateIssueIDEN  string `xml:"date_issue_iden,omitempty" db:"date_issue_iden,omitempty"`
	AuthorityIDEN  string `xml:"authority_iden,omitempty" db:"authority_iden,omitempty"`
}

// Benefit структура для льгот
type Benefit struct {
	BenefitCode    sql.NullString `xml:"benefit_code" db:"fnsicategorycode"`
	ExtBenefitCode sql.NullString `xml:"ext_benefit_code,omitempty" db:"ext_benefit_code"`
	Diagnosis      sql.NullString `xml:"diagnosis" db:"code"`
	ReceiveDate    sql.NullString `xml:"receive_date" db:"receive_date"`
	CancelDate     sql.NullString `xml:"cancel_date,omitempty" db:"cancel_date"`
}

type BenefitXML struct {
	BenefitCode    string `xml:"benefit_code,omitempty" db:"fnsicategorycode"`
	ExtBenefitCode string `xml:"ext_benefit_code,omitempty" db:"ext_benefit_code"`
	Diagnosis      string `xml:"diagnosis,omitempty" db:"code"`
	ReceiveDate    string `xml:"receive_date,omitempty" db:"receive_date"`
	CancelDate     string `xml:"cancel_date,omitempty,omitempty" db:"cancel_date"`
}

// Root структура для корневого элемента XML
type Root struct {
	XMLName     xml.Name   `xml:"root"`
	InfoSysCode string     `xml:"info_sys_code"`
	Documents   []Document `xml:"documents>document"`
}

type DocumentRow struct {
	DocumentID     sql.NullString `db:"document_id"`           // pdt."Id" as document_id
	DocDateTime    sql.NullString `db:"doc_date_time"`         // TO_CHAR(NOW(), 'YYYY-MM-DDThh:mm:ss') as doc_date_time
	RegisterID     sql.NullString `db:"frlo_id"`               // ind."FrloId" as frlo_id
	ExtCitizenID   sql.NullString `db:"ind_id"`                // ind."Id" as ind_id
	Name           sql.NullString `db:"individual_first_name"` // ind."IndividualFirstName" as individual_first_name
	Surname        sql.NullString `db:"individual_last_name"`  // ind."IndividualLastName" as individual_last_name
	Patronymic     sql.NullString `db:"individual_patronymic"` // ind."IndividualPatronymic" as individual_patronymic
	Birthdate      sql.NullString `db:"birthdate"`             // TO_CHAR(ind."IndividualBirthDate", 'YYYY-MM-DDThh:mm:ss') as birthdate
	Sex            sql.NullString `db:"sex"`                   // case when ind."IndividualSex" = 'Мужской' then 1 else 0 end as sex
	Citizenship    sql.NullString `db:"citethenship"`          // coalesce(ind."CitizenshipNumber", '643') as citethenship
	DocTypeCIT     sql.NullString `db:"doc_type_cit"`          //  '1' as 	doc_type_cit,
	DocTypeNameCIT sql.NullString `db:"doc_type_name_cit"`     //'паспорт РФ' as   doc_type_name_cit,
	Snils          sql.NullString `db:"individual_snils"`      // replace(replace(ind."IndividualSnils",'-',''), ' ','') as individual_snils
	PolicySN       sql.NullString `db:"individual_policy"`     // ind."IndividualPolicy" as individual_policy
	Region         sql.NullString `db:"region"`                // '64000' as region
	SerialCIT      sql.NullString `db:"citizenship_serie"`     // ind."CitizenshipSerie" as citizenship_serie
	NumCIT         sql.NullString `db:"citizenship_number"`    // ind."CitizenshipNumber" as citizenship_number
	DateIssueCIT   sql.NullString `db:"date_issue_cit"`        // TO_CHAR(ind."CitizenshipDateIssue", 'YYYY-MM-DDThh:mm:ss') as date_issue_cit
	DateExpiryCIT  sql.NullString `db:"date_expiry_cit"`       // TO_CHAR(ind."CitizenshipDateExpiry", 'YYYY-MM-DDThh:mm:ss') as date_expiry_cit
	AuthorityCIT   sql.NullString `db:"citezenship_authority"` // ind."CitizenshipAuthority" as citezenship_authority
	SerialIDEN     sql.NullString `db:"serial_iden"`           // ind."CredentialSerie"  as credential_serie
	NumIDEN        sql.NullString `db:"num_iden"`              // ind."CredentialNumber" as credential_number
	DateIssueIDEN  sql.NullString `db:"date_issue_iden"`       // TO_CHAR(ind."CredentialDateIssue", 'YYYY-MM-DDThh:mm:ss') as date_issue_cit (обратите внимание на возможное совпадение алиаса)
	AuthorityIDEN  sql.NullString `db:"authority_iden"`        // ind."CredentialAuthority" as authority_iden
	TS             sql.NullInt64  `db:"ts"`
}

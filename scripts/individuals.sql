SELECT
    pdt."Id" as document_id,
    TO_CHAR(NOW(), 'YYYY-MM-DDThh:mm:ss') as doc_date_time,
    ind."FrloId" as frlo_id,
     '1' as 	doc_type_cit,
      'паспорт РФ' as   doc_type_name_cit,
    ind."Id" as ind_id,
    ind."IndividualFirstName" as individual_first_name,
    ind."IndividualLastName" as individual_last_name,
    ind."IndividualPatronymic" as individual_patronymic,
    TO_CHAR(ind."IndividualBirthDate", 'YYYY-MM-DDThh:mm:ss') as birthdate,
    case when ind."IndividualSex" = 'Мужской' then 1 else 0 end as sex,
    coalesce(ind."CitizenshipNumber", '643') as citethenship,
    replace(replace(ind."IndividualSnils",'-',''), ' ','') as individual_snils,
    ind."IndividualPolicy" as individual_policy,
    '64000' as region,
    ind."CredentialSerie" as citizenship_serie,
    ind."CredentialNumber" as citizenship_number,
    TO_CHAR(ind."CredentialDateIssue", 'YYYY-MM-DDThh:mm:ss') as date_issue_cit,
    TO_CHAR(ind."CitizenshipDateExpiry", 'YYYY-MM-DDThh:mm:ss') as date_expiry_cit,
    ind."CredentialAuthority" as citezenship_authority,
     null as serial_iden,
     null as authority_iden,
     null as num_iden,
     null as date_issue_iden,
     EXTRACT(EPOCH FROM ind."ModifyDate")::int as ts
 from
     privelege_doc_tmp pdt
        inner join dct.individual ind on ind."Id" = pdt."IndividualId"
 where EXTRACT(EPOCH FROM ind."ModifyDate") > $1
 limit 2000;
select
    pc."FnsiCategoryCode" as fnsicategorycode,
    null as ext_benefit_code,
    coalesce(mkb."Code",'')as code,
    TO_CHAR(prd."PrivilegeDocumentDateStart", 'YYYYMMDD') as receive_date,
    TO_CHAR(prd."PrivilegeDocumentDateFinish", 'YYYYMMDD') as cancel_date
 from privelege_cat_tmp cat
         inner join privelege_doc_tmp doc on cat."Id" = doc."categoryId"
         inner join dct.privilege_category pc on pc."Id" = cat."Id"
         inner join dct.privilege_document_category pdc on pdc."PrivilegeDocumentId" = doc."Id"
         inner join dct.privilege_document prd on doc."Id" = prd."Id"
         inner join dct.individual ind on ind."Id" = prd."IndividualId"
         left join dct.mkb mkb on mkb."Id" = pdc."MkbId"
 where ind."Id" = $1;
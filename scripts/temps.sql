create temporary table if not exists privelege_cat_tmp
 as
 select
    pc."Id",
    pc."FnsiCategoryCode" FROM dct.privilege_category as pc inner join dct.program as pr on pc."ProgramId" = pr."Id" where pc."FnsiCategoryCode" is not NULL
  and pr."ProgramCode" NOT IN ('VZN','S_ONLP','S_7_NOZOL','S_ORFANNY')
  and pc."IsDeleted" = false;

create temporary  table if not exists privelege_doc_tmp
 as
  select pd."Id", pd."IndividualId", pt."Id" as "categoryId", pd."PrivilegeDocumentTypeId"
 from dct.privilege_document pd
         inner join dct.individual i on i."Id" = pd."IndividualId"
         inner join dct.privilege_document_category pdc on pdc."PrivilegeDocumentId" = pd."Id"
         inner join privelege_cat_tmp pt on pdc."CategoryId" = pt."Id"
 where pd."PrivilegeDocumentDateFinish" > '2020-12-01'::date
  and pd."IsDeleted" = false;
select distinct
     r."Id" as document_id,
     TO_CHAR(NOW(), 'YYYY-MM-DDThh:mm:ss') as doc_date_time,
     i."Id" as ext_citizen_id,
     i."IndividualFirstName" as name,
     i."IndividualLastName" as surname,
     COALESCE( i."IndividualPatronymic",'') as patronymic,
     TO_CHAR(i."IndividualBirthDate",'YYYY-MM-DD') as birthdate,
     case when i."IndividualSex"= 'Мужской' then '1' else '2' end as sex,
     '64000' as region,
     REPLACE(REPLACE(i."IndividualSnils",'-',''),' ','') as snils,
     '1'::text as doc_type,
     COALESCE( replace (i."CredentialSerie",' ',''),'') as serial,
     COALESCE( i."CredentialNumber", '') as num,
     COALESCE(  TO_CHAR(i."CredentialDateIssue",'YYYY-MM-DD'),'') as date_issue,
     COALESCE( i."CredentialAuthority",'') as authority,
     rs."Name" as recipe_serial,
     b."RecipeNumber" as recipe_num,
     r."Id" as ext_recipe_id,
     COALESCE(c2."Oid",'') as med_org_oid,
     b."DoctorName" as doctor_name,
     COALESCE(j."FrmrCode",'') as staff_position_code,
     COALESCE(d."Snils",'') as doctor_snils,
     COALESCE(i."IndividualCardNumber",'') as medical_card,
     pc."FnsiCategoryCode" as benefit_code,
     r."MkbCode" as mkb10_code,
     pe."SmnnCode" as drug_smnn_code,
     '' as commission_date,
     '' as commission_num,
     ri."QuantityAtom" as qty,
     TO_CHAR(r."IssueDate",'YYYY-MM-DD') as recipe_date,
     4 as recipe_expiry_code,

     case when vp."Value"='30d' then TO_CHAR(r."IssueDate"+ interval '30 day','YYYY-MM-DD')
          when vp."Value"='15d' then TO_CHAR(r."IssueDate"+ interval '15 day','YYYY-MM-DD')
          when vp."Value"='90d' then TO_CHAR(r."IssueDate"+ interval '90 day','YYYY-MM-DD')
          else TO_CHAR(r."IssueDate",'YYYY-MM-DD') end as date_expiry,
      EXTRACT(EPOCH FROM r."ModifyDate")::int as ts

 from dct.recipe r
          left join dct.recipe_item ri on ri."RecipeId" =r."Id"
          left join dct.blank b on b."Id" =r."BlankId"
          left join dct.recipe_seria rs on rs."Id" =b."RecipeSeriaId"
          left join dct.individual i on i."Id" =r."IndividualId"
          left join dct.doctor_2_contractor dc on dc."Id" =b."Doctor2ContractorId"
          left join dct.doctor d on d."Id" =dc."DoctorId"
          left join dct.contractor c on c."Id" =dc."ContractorId"
          left join dct.contractor c2 on c2."Id" =c."ParentId"
          left join dct.job j on j."Id" =dc."JobId"
          left join dct.privilege_category pc on pc."Id" =r."PrivilegeCategoryId"
          left join dct.order_lpu_item oli on oli."Id" =ri."OrderLpuItemId"
          left join dct.contract_item ci on ci."Id" =oli."ContractItemCorrectId"
          left join dct.product_esklp pe on pe."ProductId" =ci."ProductId"
          left join dct.valid_period vp on vp."Id" =r."ValidPeriodId"
 where ri."RecipeId" in (select rer."RecipeId" from dct.recipe_exchange_response rer where rer."Code"='SOLD')
   and oli."Id" is not null and r."IsVk" = false
   and pe."IsDeleted" =false
   and coalesce(pe."SmnnCode",'') != ''
   and EXTRACT(EPOCH FROM r."ModifyDate")::int > $1
 limit 2000
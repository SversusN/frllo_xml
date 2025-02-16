package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"frllo_xml/config"
	"frllo_xml/storage"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgtype"

	m "frllo_xml/models"
)

func main() {
	configPath, _ := os.Getwd()
	cfg, err := config.Initialize(filepath.Join(configPath, "config.yaml"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	db, err := storage.NewPGStorage(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	err = db.CreateTemps()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Recipes {
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			maxTSrecipe, _ := RecipeExport(db, &wg, cfg.RecipesTS, cfg.Code)
			if maxTSrecipe > 0 {
				cfg.RecipesTS = maxTSrecipe
				cfg.SaveConfigToYAML(*cfg, filepath.Join(configPath, "config.yaml"))
			}
		}()
		wg.Wait()
	}
	// Выполнение запроса для получения документов
	rows, err := db.GetDocuments(cfg.TS)
	if err != nil {
		log.Fatal(err)
	}

	var documentRows []m.DocumentRow
	documentRows, err = pgx.CollectRows(rows, pgx.RowToStructByName[m.DocumentRow])
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var documents []m.Document
	var maxTs int64 = 0
	bar := pb.StartNew(len(documentRows))

	for _, docRow := range documentRows {
		if docRow.TS.Int64 >= maxTs {
			maxTs = docRow.TS.Int64
		}
		doc := m.Document{
			DocumentID:  docRow.DocumentID.String,
			DocDateTime: docRow.DocDateTime.String,
			Citizen: m.Citizen{
				RegisterID:   docRow.RegisterID.String,
				ExtCitizenID: docRow.ExtCitizenID.String,
				Name:         docRow.Name.String,
				Surname:      docRow.Surname.String,
				Patronymic:   docRow.Patronymic.String,
				Birthdate:    docRow.Birthdate.String,
				Sex: fmt.Sprintf("%d", func() int {
					if docRow.Sex.String == "1" {
						return 1
					} else {
						return 0
					}
				}()),
				Citizenship: docRow.Citizenship.String,
				Snils:       replaceAll(docRow.Snils.String, "-", ""),
				PolicySN:    docRow.PolicySN.String,
				Region:      docRow.Region.String,
				IdentifyDocs: []m.Doc{{
					DocTypeCIT:     toInteger(docRow.DocTypeCIT.String),
					DocTypeNameCIT: docRow.DocTypeNameCIT.String,
					SerialCIT:      docRow.SerialCIT.String,
					NumCIT:         docRow.NumCIT.String,
					DateIssueCIT:   docRow.DateIssueCIT.String,
					DateExpiryCIT:  docRow.DateExpiryCIT.String,
					AuthorityCIT:   docRow.AuthorityCIT.String,
					SerialIDEN:     docRow.SerialIDEN.String,
					NumIDEN:        docRow.NumIDEN.String,
					DateIssueIDEN:  docRow.DateIssueIDEN.String,
					AuthorityIDEN:  docRow.AuthorityIDEN.String,
				}},
			},
		}

		// Получение льгот
		benefitRows, err := db.GetBenefits(docRow.ExtCitizenID.String)

		if err != nil {
			log.Fatal(err)
		}

		var benefits []m.Benefit
		benefits, err = pgx.CollectRows(benefitRows, pgx.RowToStructByName[m.Benefit])
		if err != nil {
			log.Fatal(err)
		}

		for b := range benefits {
			doc.Benefits = append(doc.Benefits, m.BenefitXML{
				BenefitCode:    benefits[b].BenefitCode.String,
				ExtBenefitCode: benefits[b].ExtBenefitCode.String,
				Diagnosis:      benefits[b].Diagnosis.String,
				ReceiveDate:    benefits[b].ReceiveDate.String,
				CancelDate:     benefits[b].CancelDate.String,
			})
		}
		defer benefitRows.Close()
		bar.Increment()
		documents = append(documents, doc)
	}
	// Создание корневого элемента XML
	root := m.Root{
		InfoSysCode: "3.058",
		Documents:   documents,
	}

	// Преобразование в XML
	xmlData, err := xml.MarshalIndent(root, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Добавление XML заголовка
	xmlHeader := xml.Header + string(xmlData)

	// Сохранение результата на диск
	filePath := fmt.Sprint(cfg.Code, "_ExportFnsiInd_", time.Now().Unix(), ".xml")
	err = os.WriteFile(filePath, []byte(xmlHeader), 0644)
	cfg.TS = maxTs
	cfg.SaveConfigToYAML(*cfg, filepath.Join(configPath, "config.yaml"))
	if err != nil {
		log.Fatalf("Unable to write file: %v\n", err)
	}
	bar.Finish()
	fmt.Printf("XML data successfully written to %s\n", filePath)
}

func RecipeExport(storage storage.Storage, wg *sync.WaitGroup, recipeTS int64, code string) (maxRecipeTs int64, err error) {

	rows, err := storage.GetRecipes(recipeTS)
	var maxTS atomic.Int64
	maxTS.Store(0)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var documents []m.DocumentRecipe

	for rows.Next() {
		var doc m.DocumentRecipe
		var citizen m.CitizenRecipe
		var identifyDoc m.DocRecipe
		var recipe m.Recipe

		err := rows.Scan(
			&doc.DocumentID,
			&doc.DocDateTime,
			&citizen.ExtCitizenID,
			&citizen.Name,
			&citizen.Surname,
			&citizen.Patronymic,
			&citizen.Birthdate,
			&citizen.Sex,
			&citizen.Region,
			&citizen.Snils,
			&identifyDoc.DocType,
			&identifyDoc.Serial,
			&identifyDoc.Num,
			&identifyDoc.DateIssue,
			&identifyDoc.Authority,
			&recipe.RecipeSerial,
			&recipe.RecipeNum,
			&recipe.ExtRecipeID,
			&recipe.MedOrgOID,
			&recipe.DoctorName,
			&recipe.StaffPositionCode,
			&recipe.DoctorSnils,
			&recipe.MedicalCard,
			&recipe.BenefitCode,
			&recipe.MKB10Code,
			&recipe.DrugSmnnCode,
			&recipe.CommissionDate,
			&recipe.CommissionNum,
			&recipe.Qty,
			&recipe.RecipeDate,
			&recipe.RecipeExpiryCode,
			&recipe.DateExpiry,
			&recipe.TS,
		)
		if err != nil {
			log.Fatal(err)
		}
		if recipe.TS > int(maxTS.Load()) {
			maxTS.Store(int64(recipe.TS))
		}

		// Заполнение структуры Citizen
		citizen.IdentifyDocs = append(citizen.IdentifyDocs, identifyDoc)
		doc.Citizen = citizen

		// Заполнение структуры Recipe
		doc.Recipe = recipe

		documents = append(documents, doc)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Создание корневого элемента XML
	root := m.RootRecipe{
		InfoSysCode: code,
		Documents:   documents,
	}

	// Преобразование в XML
	xmlData, err := xml.MarshalIndent(root, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Добавление XML заголовка
	xmlHeader := xml.Header + string(xmlData)

	// Сохранение результата на диск

	filePath := fmt.Sprint(code, "_ExportWrittenRecipes_", time.Now().Unix(), ".xml")
	err = os.WriteFile(filePath, []byte(xmlHeader), 0644)
	if err != nil {
		log.Fatalf("Unable to write file: %v\n", err)
	}

	fmt.Printf("XML data successfully written to %s\n", filePath)

	return maxTS.Load(), nil
}

// Функция для замены всех вхождений подстроки
func replaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}
func toInteger(s string) int {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	return int(i)
}

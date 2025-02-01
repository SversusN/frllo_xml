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
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgtype"

	m "frllo_xml/models"
)

// Document структура для документа

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

	err = db.CreateTemps()
	if err != nil {
		log.Fatal(err)
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
					DocTypeCIT:     1,
					DocTypeNameCIT: "",
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
		benefitRows, err := db.GetBenefits(docRow.DocumentID.String)

		if err != nil {
			log.Fatal(err)
		}
		defer benefitRows.Close()

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
	filePath := fmt.Sprint("3.058_ExportFnsiInd_", time.Now().Unix(), ".xml")
	err = os.WriteFile(filePath, []byte(xmlHeader), 0644)
	cfg.TS = maxTs
	cfg.SaveConfigToYAML(*cfg, filepath.Join(configPath, "config.yaml"))
	if err != nil {
		log.Fatalf("Unable to write file: %v\n", err)
	}
	bar.Finish()
	fmt.Printf("XML data successfully written to %s\n", filePath)
}

// Функция для замены всех вхождений подстроки
func replaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

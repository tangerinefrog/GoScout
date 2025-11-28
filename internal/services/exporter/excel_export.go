package exporter

import (
	"context"
	"fmt"
	"job-scraper/internal/data"
	"job-scraper/internal/data/repositories"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type exporter struct {
	db *data.DB
}

const filename = "jobs.xlsx"

func NewExcelExporter(db *data.DB) *exporter {
	return &exporter{
		db: db,
	}
}

func (e *exporter) ExportToExcel(ctx context.Context) (string, error) {
	jobsRepo := repositories.NewJobsRepo(e.db)
	jobs, err := jobsRepo.List(ctx)
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()
	sheetName := "Jobs"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", err
	}

	headers := []string{"ID", "URL", "Title", "Company", "Location", "Date Posted", "Status"}
	for colIdx, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		f.SetCellValue(sheetName, cell, h)
	}

	for rowIdx, j := range jobs {
		row := rowIdx + 2

		var grade string
		if j.Grade != nil {
			grade = strconv.Itoa(*j.Grade)
		}
		var reasoning string
		if j.GradeReasoning != nil {
			reasoning = *j.GradeReasoning
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), j.Id)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), j.Url)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), j.Title)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), j.Company)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), j.Location)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), j.DatePosted.Local().Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), j.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), grade)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), reasoning)
	}

	f.SetActiveSheet(index)

	os.Remove(filename)
	if err := f.SaveAs(filename); err != nil {
		fmt.Println(err)
	}

	return filename, nil
}

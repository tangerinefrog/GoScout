package exporter

import (
	"fmt"
	"job-scraper/internal/data"
	"job-scraper/internal/data/repositories"
	"os"

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

func (e *exporter) ExportToExcel() (string, error) {
	jobsRepo := repositories.NewJobsRepo(e.db)
	jobs, err := jobsRepo.List()
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()
	sheetName := "Jobs"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", err
	}

	headers := []string{"ID", "URL", "Title", "Company", "Date Posted", "Status"}
	for colIdx, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		f.SetCellValue(sheetName, cell, h)
	}

	for rowIdx, j := range jobs {
		row := rowIdx + 2

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), j.Id)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), j.Url)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), j.Title)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), j.Company)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), j.DatePosted.Local().Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), j.Status)
	}

	f.SetActiveSheet(index)

	os.Remove(filename)
	if err := f.SaveAs(filename); err != nil {
		fmt.Println(err)
	}

	return filename, nil
}

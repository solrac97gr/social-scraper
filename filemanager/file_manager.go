package filemanager

import (
	"fmt"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

// FileManager defines the interface for file operations
type FileManager interface {
	ReadLinksFromExcel(filePath string) []string
	SaveResultsToExcel(data [][]string, outputPath string)
	EstimateProcessingTime(filePath string) (int, error)
}

// FileManagerImpl implements the FileManager interface
type FileManagerImpl struct{}

// NewFileManager creates a new FileManager instance
func NewFileManager() FileManager {
	return &FileManagerImpl{}
}

// ReadLinksFromExcel reads links from an Excel file
func (fm *FileManagerImpl) ReadLinksFromExcel(filePath string) []string {
	// Open Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatalf("Failed to open Excel file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Failed to close Excel file: %v", err)
		}
	}()

	// Get all sheet names
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		log.Fatal("No sheets found in Excel file")
	}

	// Get all rows in the first sheet
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		log.Fatalf("Failed to get rows: %v", err)
	}

	links := make([]string, 0)

	// Extract links from rows
	for _, row := range rows {
		for _, cell := range row {
			if strings.Contains(cell, "t.me/") || strings.Contains(cell, "telegram.me/") || strings.Contains(cell, "rutube.ru/") || strings.Contains(cell, "vk.com/") || strings.Contains(cell, "instagram.com/") {
				links = append(links, strings.TrimSpace(cell))
			}
		}
	}

	if len(links) == 0 {
		log.Fatal("No valid links found in the Excel file")
	}

	return links
}

// SaveResultsToExcel saves the extraction results to an Excel file
func (fm *FileManagerImpl) SaveResultsToExcel(data [][]string, outputPath string) {
	f := excelize.NewFile()

	// Set column headers to be wider
	f.SetColWidth("Sheet1", "A", "A", 30)
	f.SetColWidth("Sheet1", "B", "B", 15)
	f.SetColWidth("Sheet1", "C", "C", 40)

	// Style for header
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#DDEBF7"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	// Write data to Excel
	for r, row := range data {
		for c, cellValue := range row {
			cellName, _ := excelize.CoordinatesToCellName(c+1, r+1)
			f.SetCellValue("Sheet1", cellName, cellValue)
		}

		// Apply header style to first row
		if r == 0 {
			f.SetRowHeight("Sheet1", 1, 20)
			for c := range row {
				cellName, _ := excelize.CoordinatesToCellName(c+1, r+1)
				f.SetCellStyle("Sheet1", cellName, cellName, headerStyle)
			}
		}
	}

	// Save the Excel file
	if err := f.SaveAs(outputPath); err != nil {
		log.Fatalf("Failed to save Excel file: %v", err)
	}
}

var (
	estimationPerLink = map[string]int{
		"vk": 26, // seconds
		"tg": 1,  // seconds
		"ig": 1,  // seconds
		"rt": 1,  // seconds
	}
)

// EstimateProcessingTime estimates the processing time based on the number of links
func (fm *FileManagerImpl) EstimateProcessingTime(filePath string) (int, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	// Get all sheet names
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return 0, fmt.Errorf("no sheets found in Excel file")
	}

	// Get all rows in the first sheet
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return 0, err
	}

	linkCount := 0
	for _, row := range rows {
		for _, cell := range row {
			if cell != "" && (len(cell) > 4 && cell[:4] == "http") {
				if strings.Contains(cell, "vk.com/") {
					linkCount += estimationPerLink["vk"]
				} else if strings.Contains(cell, "t.me/") || strings.Contains(cell, "telegram.me/") {
					linkCount += estimationPerLink["tg"]
				} else if strings.Contains(cell, "instagram.com/") {
					linkCount += estimationPerLink["ig"]
				} else if strings.Contains(cell, "rutube.ru/") {
					linkCount += estimationPerLink["rt"]
				} else {
					linkCount += 0 // Default estimation for unknown links
				}
			}
		}
	}

	// Estimate the processing time (1 second per link)
	estimatedTime := linkCount

	return estimatedTime, nil
}

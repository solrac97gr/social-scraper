package filemanager

import (
	"encoding/csv"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

// FileManager defines the interface for file operations
type FileManager interface {
	ReadLinksFromExcel(filePath string) []string
	ReadLinksFromFile(filePath string) []string // Generic method that detects file type
	ReadLinksFromCSV(filePath string) []string
	ReadLinksFromText(content string) []string
	SaveResultsToExcel(data [][]string, outputPath string)
	EstimateProcessingTime(filePath string) (int, error)
}

// FileManagerImpl implements the FileManager interface
type FileManagerImpl struct{}

// NewFileManager creates a new FileManager instance
func NewFileManager() FileManager {
	return &FileManagerImpl{}
}

// normalizeLink standardizes a URL to use https and a canonical domain.
func normalizeLink(link string) string {
	link = strings.TrimSpace(link)
	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		link = "https://" + link
	}

	u, err := url.Parse(link)
	if err != nil {
		log.Printf("Warning: Could not parse link '%s'. Using it as is. Error: %v", link, err)
		return link
	}

	u.Scheme = "https"

	host := strings.ToLower(u.Host)
	if strings.HasSuffix(host, "vk.com") {
		u.Host = "vk.com"
	} else if strings.HasSuffix(host, "instagram.com") {
		u.Host = "instagram.com"
	} else if strings.HasSuffix(host, "rutube.ru") {
		u.Host = "rutube.ru"
	} else if strings.HasSuffix(host, "telegram.me") {
		u.Host = "t.me"
	} else if strings.HasSuffix(host, "t.me") {
		u.Host = "t.me"
	}

	return u.String()
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
				links = append(links, normalizeLink(cell))
			}
		}
	}

	if len(links) == 0 {
		log.Fatal("No valid links found in the Excel file")
	}

	return links
}

// ReadLinksFromFile is a generic method that detects file type and reads links accordingly
func (fm *FileManagerImpl) ReadLinksFromFile(filePath string) []string {
	ext := strings.ToLower(filepath.Ext(filePath))

	// Check if it's a text file (for text input handling)
	if strings.Contains(filePath, "_text_input.txt") {
		// Read file content and treat as text input
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to read text file: %v", err)
		}
		return fm.ReadLinksFromText(string(content))
	}

	switch ext {
	case ".xlsx", ".xls":
		return fm.ReadLinksFromExcel(filePath)
	case ".csv":
		return fm.ReadLinksFromCSV(filePath)
	default:
		log.Fatalf("Unsupported file format: %s. Supported formats: .xlsx, .xls, .csv", ext)
		return nil
	}
}

// ReadLinksFromCSV reads links from a CSV file
func (fm *FileManagerImpl) ReadLinksFromCSV(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close CSV file: %v", err)
		}
	}()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %v", err)
	}

	links := make([]string, 0)

	// Extract links from CSV records
	for _, record := range records {
		for _, cell := range record {
			if strings.Contains(cell, "t.me/") || strings.Contains(cell, "telegram.me/") ||
				strings.Contains(cell, "rutube.ru/") || strings.Contains(cell, "vk.com/") ||
				strings.Contains(cell, "instagram.com/") || strings.Contains(cell, "youtube.com/") {
				links = append(links, normalizeLink(cell))
			}
		}
	}

	if len(links) == 0 {
		log.Fatal("No valid links found in the CSV file")
	}

	return links
}

// ReadLinksFromText reads links from comma-separated text input
func (fm *FileManagerImpl) ReadLinksFromText(content string) []string {
	// Split by various separators: comma, newline, space
	content = strings.ReplaceAll(content, "\n", ",")
	content = strings.ReplaceAll(content, " ", ",")
	rawLinks := strings.Split(content, ",")
	links := make([]string, 0)

	for _, link := range rawLinks {
		trimmed := strings.TrimSpace(link)
		if trimmed != "" && (strings.Contains(trimmed, "t.me/") || strings.Contains(trimmed, "telegram.me/") ||
			strings.Contains(trimmed, "rutube.ru/") || strings.Contains(trimmed, "vk.com/") ||
			strings.Contains(trimmed, "instagram.com/") || strings.Contains(trimmed, "youtube.com/")) {
			links = append(links, normalizeLink(trimmed))
		}
	}

	if len(links) == 0 {
		log.Fatal("No valid links found in the text input")
	}

	return links
}

// SaveResultsToExcel saves the extraction results to an Excel file
func (fm *FileManagerImpl) SaveResultsToExcel(data [][]string, outputPath string) {
	f := excelize.NewFile()

	// Set column headers to be wider
	_ = f.SetColWidth("Sheet1", "A", "A", 30)
	_ = f.SetColWidth("Sheet1", "B", "B", 15)
	_ = f.SetColWidth("Sheet1", "C", "C", 40)

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
			_ = f.SetCellValue("Sheet1", cellName, cellValue)
		}

		// Apply header style to first row
		if r == 0 {
			_ = f.SetRowHeight("Sheet1", 1, 20)
			for c := range row {
				cellName, _ := excelize.CoordinatesToCellName(c+1, r+1)
				_ = f.SetCellStyle("Sheet1", cellName, cellName, headerStyle)
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
	// Use the universal file reader to get links count
	links := fm.ReadLinksFromFile(filePath)

	// Estimate based on platform type
	estimatedTime := 0
	for _, link := range links {
		if strings.Contains(link, "vk.com/") {
			estimatedTime += estimationPerLink["vk"]
		} else if strings.Contains(link, "t.me/") || strings.Contains(link, "telegram.me/") {
			estimatedTime += estimationPerLink["tg"]
		} else if strings.Contains(link, "instagram.com/") {
			estimatedTime += estimationPerLink["ig"]
		} else if strings.Contains(link, "rutube.ru/") {
			estimatedTime += estimationPerLink["rt"]
		} else {
			estimatedTime += 1 // Default estimation for unknown links
		}
	}

	return estimatedTime, nil
}

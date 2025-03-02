package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/xuri/excelize/v2"
)

func main() {
	// Check if input argument is provided
	if len(os.Args) < 2 {
		log.Fatal("Please provide the path to the Excel file as an argument")
	}

	inputFile := os.Args[1]
	links := readExcelLinks(inputFile)

	results := make([][]string, 0, len(links)+1)
	// Add header row
	results = append(results, []string{"Channel Name", "Followers Count", "Original Link"})

	// Process each link
	for _, link := range links {
		fmt.Printf("Processing: %s\n", link)
		channelName, followersCount := getChannelInfo(link)

		// Add to results
		results = append(results, []string{channelName, followersCount, link})

		// Avoid hitting rate limits
		time.Sleep(1 * time.Second)
	}

	// Create output file
	outputFilePath := "telegram_channels_followers.xlsx"
	createOutputExcel(results, outputFilePath)

	fmt.Printf("\nSuccess! Results saved to %s\n", outputFilePath)
}

func readExcelLinks(filePath string) []string {
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
			if strings.Contains(cell, "t.me/") || strings.Contains(cell, "telegram.me/") {
				links = append(links, strings.TrimSpace(cell))
			}
		}
	}

	if len(links) == 0 {
		log.Fatal("No Telegram links found in the Excel file")
	}

	return links
}

func getChannelInfo(link string) (string, string) {
	// Format the link to ensure it's accessible via http
	if !strings.HasPrefix(link, "http") {
		link = "https://" + link
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request
	resp, err := client.Get(link)
	if err != nil {
		log.Printf("Error fetching %s: %v", link, err)
		return "Error", "N/A"
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Status code error: %d %s", resp.StatusCode, resp.Status)
		return "Error", "N/A"
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error parsing HTML: %v", err)
		return "Error", "N/A"
	}

	// Extract channel name
	channelName := "Unknown"
	doc.Find("div.tgme_page_title").Each(func(i int, s *goquery.Selection) {
		channelName = strings.TrimSpace(s.Text())
	})

	// Extract followers count
	followersText := "0"
	doc.Find("div.tgme_page_extra").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(text, "subscriber") || strings.Contains(text, "member") || strings.Contains(text, "follower") {
			re := regexp.MustCompile(`[\d\s]+`)
			matches := re.FindString(text)
			if matches != "" {
				// Remove spaces and convert to number
				followersText = strings.ReplaceAll(matches, " ", "")
			}
		}
	})

	return channelName, followersText
}

func createOutputExcel(data [][]string, outputPath string) {
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

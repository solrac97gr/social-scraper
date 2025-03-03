package ruregistration

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// CheckRegistrationStatus checks if the given link is registered on the specified website
func CheckRegistrationStatus(link string) bool {
	fmt.Println("Checking registration status for", link)
	defer fmt.Println("Finished checking registration status for", link)
	cmd := exec.Command("node", "scripts/ru-registration.js", link)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error executing Puppeteer script: %v", err)
		return false
	}

	var result struct {
		IsRegistered bool `json:"isRegistered"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		log.Printf("Error parsing JSON output: %v", err)
		return false
	}

	return result.IsRegistered
}

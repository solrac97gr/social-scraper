package ruregistration

import (
	"encoding/json"
	"log"
	"os/exec"
	"sync/atomic"
)

var pendingChecks int32

// CheckRegistrationStatus checks if the given link is registered on the specified website
func CheckRegistrationStatus(link string, semaphore chan struct{}) bool {
	println("Checking registration status for:", link)
	atomic.AddInt32(&pendingChecks, 1)
	log.Printf("Checking registration status for: %s (Pending checks: %d)", link, atomic.LoadInt32(&pendingChecks))
	cmd := exec.Command("node", "scripts/ru-registration.js", link)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error executing Puppeteer script: %v", err)
		atomic.AddInt32(&pendingChecks, -1)
		<-semaphore // Release semaphore
		return false
	}

	var result struct {
		IsRegistered bool `json:"isRegistered"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		log.Printf("Error parsing JSON (%s) output: %v", link, output)
		atomic.AddInt32(&pendingChecks, -1)
		<-semaphore // Release semaphore
		return false
	}
	defer func() {
		atomic.AddInt32(&pendingChecks, -1)
		<-semaphore // Release semaphore
		log.Printf("Finished checking registration status for: %s (Pending checks: %d) %v", link, atomic.LoadInt32(&pendingChecks), result)
	}()

	return result.IsRegistered
}

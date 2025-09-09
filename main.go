package main

import (
	"context"
	"fmt"
	"job-scraper/config"
	"job-scraper/models"
	"job-scraper/scraper"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Create scraper
	jobScraper := scraper.NewJobScraper()

	// Test scraping one site
	ctx := context.Background()
	result, err := jobScraper.ScrapeSite(ctx, cfg.Sites[0])
	if err != nil {
		log.Printf("Error scraping %s: %v", cfg.Sites[0].Name, err)
	} else {
		fmt.Printf("Scraped %s: Found %d jobs in %v\n", result.Site, len(result.Jobs), result.Duration)

		// Print first few jobs
		for i, job := range result.Jobs {
			if i >= 3 { // Only show first 3 jobs
				break
			}
			fmt.Printf("Job %d: %s at %s\n", i+1, job.Title, job.Company)
		}
	}

	// Save results to YAML file
	filename := fmt.Sprintf("jobs_%s_%s.yaml", result.Site, time.Now().Format("2006-01-02_15-04-05"))
	err = saveResultsToYAML(result, filename)
	if err != nil {
		log.Printf("Error saving results to file: %v", err)
	} else {
		fmt.Printf("\nResults saved to: %s\n", filename)
	}
}

// saveResultsToYAML saves scraping results to a YAML file
func saveResultsToYAML(result *models.ScrapingResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)

	err = encoder.Encode(result)
	if err != nil {
		return fmt.Errorf("failed to encode YAML: %w", err)
	}

	return nil
}

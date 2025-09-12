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
	ctx := context.Background()

	// Scrape regular sites
	fmt.Println("=== Scraping Regular Sites ===")
	regularResults, err := jobScraper.ScrapeAllSites(ctx, cfg.Sites)
	if err != nil {
		log.Printf("Error scraping regular sites: %v", err)
	}

	// Scrape API sites
	fmt.Println("\n=== Scraping API Sites ===")
	apiResults, err := jobScraper.ScrapeAllAPISites(ctx, cfg.APISites)
	if err != nil {
		log.Printf("Error scraping API sites: %v", err)
	}

	// Combine all results
	allResults := append(regularResults, apiResults...)

	// Print summary
	totalJobs := 0
	for _, result := range allResults {
		fmt.Printf("Scraped %s: Found %d jobs in %v\n", result.Site, len(result.Jobs), result.Duration)
		totalJobs += len(result.Jobs)
	}
	fmt.Printf("\nTotal jobs found: %d\n", totalJobs)

	// Print sample jobs from each source
	for _, result := range allResults {
		if len(result.Jobs) > 0 {
			fmt.Printf("\n--- Sample jobs from %s ---\n", result.Site)
			for i, job := range result.Jobs {
				if i >= 3 { // Only show first 3 jobs per source
					break
				}
				fmt.Printf("Job %d: %s at %s (%s)\n", i+1, job.Title, job.Company, job.Location)
			}
		}
	}

	// Save results to YAML file
	if len(allResults) > 0 {
		filename := fmt.Sprintf("jobs_combined_%s.yaml", time.Now().Format("2006-01-02_15-04-05"))
		err = saveAllResultsToYAML(allResults, filename)
		if err != nil {
			log.Printf("Error saving results to file: %v", err)
		} else {
			fmt.Printf("\nResults saved to: %s\n", filename)
		}
	}
}

// saveAllResultsToYAML saves all scraping results to a YAML file
func saveAllResultsToYAML(results []models.ScrapingResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)

	err = encoder.Encode(results)
	if err != nil {
		return fmt.Errorf("failed to encode YAML: %w", err)
	}

	return nil
}

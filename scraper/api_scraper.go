package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"job-scraper/config"
	"job-scraper/models"
)

// APIScraper handles API-based job scraping
type APIScraper struct {
	client *http.Client
}

// NewAPIScraper creates a new API scraper instance
func NewAPIScraper() *APIScraper {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	return &APIScraper{
		client: client,
	}
}

// ScrapeAPISite scrapes a single API site for job listings
func (as *APIScraper) ScrapeAPISite(ctx context.Context, apiSite config.APISite) (*models.ScrapingResult, error) {
	start := time.Now()
	var jobs []models.JobListing
	var err error

	switch apiSite.Type {
	case "api":
		jobs, err = as.scrapeGenericAPI(ctx, apiSite)
		if err != nil {
			return &models.ScrapingResult{
				Site:     apiSite.Name,
				Jobs:     jobs,
				Error:    err,
				Duration: time.Since(start),
			}, err
		}
	case "greenhouse_api":
		jobs, err = as.scrapeGreenhouseAPI(ctx, apiSite)
		if err != nil {
			return &models.ScrapingResult{
				Site:     apiSite.Name,
				Jobs:     jobs,
				Error:    err,
				Duration: time.Since(start),
			}, err
		}
	}

	return &models.ScrapingResult{
		Site:     apiSite.Name,
		Jobs:     jobs,
		Duration: time.Since(start),
	}, nil
}

// scrapeGenericAPI scrapes a generic API endpoint
func (as *APIScraper) scrapeGenericAPI(ctx context.Context, apiSite config.APISite) ([]models.JobListing, error) {
	var jobs []models.JobListing

	// Build URL with query parameters
	reqURL, err := url.Parse(apiSite.URL)
	if err != nil {
		return jobs, fmt.Errorf("invalid URL: %v", err)
	}

	// Add query parameters
	query := reqURL.Query()
	for key, value := range apiSite.Params {
		query.Add(key, value)
	}
	reqURL.RawQuery = query.Encode()

	// Create request
	req, err := http.NewRequestWithContext(ctx, apiSite.Method, reqURL.String(), nil)
	if err != nil {
		return jobs, err
	}

	// Set headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	// Make request
	resp, err := as.client.Do(req)
	if err != nil {
		return jobs, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return jobs, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return jobs, err
	}

	// Parse JSON response
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return jobs, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Extract jobs from response (this will need to be customized per API)
	jobs = as.extractJobsFromGenericAPI(response, apiSite.Name)

	return jobs, nil
}

// scrapeGreenhouseAPI scrapes multiple Greenhouse company boards
func (as *APIScraper) scrapeGreenhouseAPI(ctx context.Context, apiSite config.APISite) ([]models.JobListing, error) {
	var allJobs []models.JobListing

	for _, company := range apiSite.Companies {
		jobs, err := as.scrapeSingleGreenhouseBoard(ctx, apiSite.BaseURL, company)
		if err != nil {
			fmt.Printf("Error scraping %s: %v\n", company, err)
			continue
		}
		allJobs = append(allJobs, jobs...)
	}

	return allJobs, nil
}

// scrapeSingleGreenhouseBoard scrapes a single Greenhouse company board
func (as *APIScraper) scrapeSingleGreenhouseBoard(ctx context.Context, baseURL, company string) ([]models.JobListing, error) {
	var jobs []models.JobListing

	// Build API URL
	apiURL := fmt.Sprintf("%s/%s/jobs?content=true", baseURL, company)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return jobs, err
	}

	// Set headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	// Make request
	resp, err := as.client.Do(req)
	if err != nil {
		return jobs, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return jobs, fmt.Errorf("API returned status %d for %s", resp.StatusCode, company)
	}

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return jobs, err
	}

	// Parse Greenhouse API response
	var response struct {
		Jobs []struct {
			ID       int    `json:"id"`
			Title    string `json:"title"`
			Location struct {
				Name string `json:"name"`
			} `json:"location"`
			AbsoluteURL string `json:"absolute_url"`
			CreatedAt   string `json:"created_at"`
			UpdatedAt   string `json:"updated_at"`
		} `json:"jobs"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return jobs, fmt.Errorf("failed to parse Greenhouse JSON: %v", err)
	}

	// Convert to our job format
	for _, job := range response.Jobs {
		// Check if job is remote and relevant
		if as.isRemoteJob(job.Title, company, job.Location.Name, "") && as.isRelevantRole(job.Title, "") {
			jobListing := models.JobListing{
				ID:         fmt.Sprintf("greenhouse-%s-%d", company, job.ID),
				Title:      job.Title,
				Company:    strings.Title(company), // Capitalize company name
				Location:   job.Location.Name,
				URL:        job.AbsoluteURL,
				Source:     fmt.Sprintf("Greenhouse - %s", strings.Title(company)),
				ScrapedAt:  time.Now(),
				PostedDate: time.Now(), // Greenhouse doesn't always provide creation date
			}
			jobs = append(jobs, jobListing)
		}
	}

	return jobs, nil
}

// extractJobsFromGenericAPI extracts jobs from a generic API response
func (as *APIScraper) extractJobsFromGenericAPI(response map[string]interface{}, source string) []models.JobListing {
	var jobs []models.JobListing

	// This is a generic parser - you'll need to customize based on the specific API
	// For now, we'll look for common patterns
	if jobsData, ok := response["jobs"].([]interface{}); ok {
		for _, jobData := range jobsData {
			if jobMap, ok := jobData.(map[string]interface{}); ok {
				title := as.getStringFromMap(jobMap, "title", "name", "position")
				company := as.getStringFromMap(jobMap, "company", "employer", "organization")
				location := as.getStringFromMap(jobMap, "location", "city", "place")
				jobURL := as.getStringFromMap(jobMap, "url", "link", "absolute_url")

				if title != "" && as.isRemoteJob(title, company, location, "") && as.isRelevantRole(title, "") {
					job := models.JobListing{
						ID:         fmt.Sprintf("%s-%d", source, len(jobs)+1),
						Title:      title,
						Company:    company,
						Location:   location,
						URL:        jobURL,
						Source:     source,
						ScrapedAt:  time.Now(),
						PostedDate: time.Now(),
					}
					jobs = append(jobs, job)
				}
			}
		}
	}

	return jobs
}

// getStringFromMap safely extracts a string value from a map, trying multiple keys
func (as *APIScraper) getStringFromMap(m map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if value, ok := m[key].(string); ok && value != "" {
			return value
		}
	}
	return ""
}

// isRemoteJob checks if a job is remote based on various indicators
func (as *APIScraper) isRemoteJob(title, company, location, fullText string) bool {
	// Convert to lowercase for case-insensitive matching
	titleLower := strings.ToLower(title)
	locationLower := strings.ToLower(location)
	fullTextLower := strings.ToLower(fullText)

	// Remote job indicators
	remoteKeywords := []string{
		"remote",
		"work from home",
		"wfh",
		"virtual",
		"distributed",
		"telecommute",
		"flexible location",
		"anywhere",
		"global",
		"worldwide",
		"us-remote",
		"remote in",
		"remote us",
		"remote canada",
	}

	// Check if any remote keywords appear in title, location, or full text
	for _, keyword := range remoteKeywords {
		if strings.Contains(titleLower, keyword) ||
			strings.Contains(locationLower, keyword) ||
			strings.Contains(fullTextLower, keyword) {
			return true
		}
	}

	// Check for common non-remote indicators (if found, likely not remote)
	nonRemoteKeywords := []string{
		"on-site",
		"onsite",
		"in-person",
		"office",
		"headquarters",
		"relocation required",
		"must be local",
	}

	for _, keyword := range nonRemoteKeywords {
		if strings.Contains(titleLower, keyword) ||
			strings.Contains(locationLower, keyword) ||
			strings.Contains(fullTextLower, keyword) {
			return false
		}
	}

	// If location is empty, might be remote
	if location == "" {
		return true
	}

	return false
}

// isRelevantRole checks if a job title/description matches relevant engineering roles
func (as *APIScraper) isRelevantRole(title, fullText string) bool {
	// Convert to lowercase for case-insensitive matching
	titleLower := strings.ToLower(title)
	fullTextLower := strings.ToLower(fullText)

	// Relevant role keywords
	relevantKeywords := []string{
		"system engineer",
		"systems engineer",
		"devops",
		"dev ops",
		"cloud engineer",
		"sre",
		"site reliability engineer",
		"platform engineer",
		"infrastructure engineer",
		"reliability engineer",
		"automation engineer",
		"build engineer",
		"release engineer",
		"deployment engineer",
		"kubernetes engineer",
		"container engineer",
		"aws engineer",
		"azure engineer",
		"gcp engineer",
		"google cloud engineer",
		"terraform engineer",
		"ansible engineer",
		"jenkins engineer",
		"ci/cd engineer",
		"monitoring engineer",
		"observability engineer",
		"security engineer",
		"compliance engineer",
		"network engineer",
		"linux engineer",
		"unix engineer",
		"operations engineer",
		"ops engineer",
		"production engineer",
		"backend engineer",
		"api engineer",
		"microservices engineer",
		"distributed systems engineer",
		"scalability engineer",
		"performance engineer",
		"data engineer",
		"ml engineer",
		"machine learning engineer",
		"ai engineer",
		"artificial intelligence engineer",
	}

	// Check if any relevant keywords appear in title or full text
	for _, keyword := range relevantKeywords {
		if strings.Contains(titleLower, keyword) ||
			strings.Contains(fullTextLower, keyword) {
			return true
		}
	}

	// Check for specific technology keywords that indicate relevant roles
	techKeywords := []string{
		"kubernetes",
		"docker",
		"terraform",
		"ansible",
		"puppet",
		"chef",
		"jenkins",
		"gitlab ci",
		"github actions",
		"aws",
		"azure",
		"gcp",
		"google cloud",
		"amazon web services",
		"microservices",
		"containerization",
		"orchestration",
		"monitoring",
		"observability",
		"prometheus",
		"grafana",
		"elk stack",
		"elasticsearch",
		"splunk",
		"datadog",
		"new relic",
		"pagerduty",
		"incident response",
		"disaster recovery",
		"high availability",
		"load balancing",
		"auto scaling",
		"infrastructure as code",
		"configuration management",
		"linux",
		"unix",
		"bash",
		"shell scripting",
		"python",
		"golang",
		"go",
		"ruby",
		"powershell",
		"networking",
		"security",
		"compliance",
		"automation",
	}

	// Count how many tech keywords appear
	techCount := 0
	for _, keyword := range techKeywords {
		if strings.Contains(titleLower, keyword) ||
			strings.Contains(fullTextLower, keyword) {
			techCount++
		}
	}

	// If we find 2 or more tech keywords, it's likely a relevant role
	return techCount >= 2
}

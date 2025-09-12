package scraper

import (
	"context"
	"fmt"
	"job-scraper/config"
	"job-scraper/models"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// Scraper defines the interface for job scraping
type Scraper interface {
	// ScrapeSite scrapes a single site for job listings
	ScrapeSite(ctx context.Context, site config.Site) (*models.ScrapingResult, error)

	// ScrapeAPISite scrapes a single API site for job listings
	ScrapeAPISite(ctx context.Context, apiSite config.APISite) (*models.ScrapingResult, error)

	// ScrapeAllSites scrapes all configured sites
	ScrapeAllSites(ctx context.Context, sites []config.Site) ([]models.ScrapingResult, error)

	// ScrapeAllAPISites scrapes all configured API sites
	ScrapeAllAPISites(ctx context.Context, apiSites []config.APISite) ([]models.ScrapingResult, error)
}

// JobScraper implements the Scraper interface
type JobScraper struct {
	client *http.Client
}

// redirectPolicyFunc handles HTTP redirects
func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	// Limit redirects to prevent infinite loops
	if len(via) >= 10 {
		return http.ErrUseLastResponse
	}
	return nil
}

// NewJobScraper creates a new job scraper instance
func NewJobScraper() *JobScraper {
	// 1. Initialize HTTP client with proper timeouts
	client := &http.Client{
		Timeout: 30 * time.Second, // Total request timeout
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second, // Connection timeout
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second, // TLS handshake timeout
			ResponseHeaderTimeout: 10 * time.Second, // Response header timeout
			ExpectContinueTimeout: 1 * time.Second,  // Expect: 100-continue timeout
			MaxIdleConns:          100,              // Maximum idle connections
			MaxIdleConnsPerHost:   10,               // Maximum idle connections per host
			IdleConnTimeout:       90 * time.Second, // Idle connection timeout
		},
		CheckRedirect: redirectPolicyFunc,
	}

	return &JobScraper{
		client: client,
	}
}

// ScrapeSite scrapes a single site for job listings
func (js *JobScraper) ScrapeSite(ctx context.Context, site config.Site) (*models.ScrapingResult, error) {
	start := time.Now()
	var jobs []models.JobListing

	// Create a new collector
	c := colly.NewCollector()

	// Set user agent to avoid being blocked
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	// Set up the callback for job listings
	c.OnHTML(site.Selector, func(e *colly.HTMLElement) {
		// Extract basic job information
		title := strings.TrimSpace(e.ChildText("h2, h3, .title, .job-title, a, .job-title-text"))
		company := strings.TrimSpace(e.ChildText(".company, .company-name, .employer, .job-card-container__company-name"))
		location := strings.TrimSpace(e.ChildText(".location, .job-location, .where, .job-card-container__metadata-item"))
		url := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))

		// If URL is empty, try to get it from the element itself
		if url == "" {
			url = e.Request.AbsoluteURL(e.Attr("href"))
		}

		// Check if this is a remote job
		isRemote := js.isRemoteJob(title, company, location, e.Text)

		// Check if this is a relevant role
		isRelevantRole := js.isRelevantRole(title, e.Text)

		// Create job listing if we have at least a title, it's remote, and it's a relevant role
		if title != "" && isRemote && isRelevantRole {
			job := models.JobListing{
				ID:         fmt.Sprintf("%s-%d", site.Name, len(jobs)+1),
				Title:      title,
				Company:    company,
				Location:   location,
				URL:        url,
				Source:     site.Name,
				ScrapedAt:  time.Now(),
				PostedDate: time.Now(), // We'll improve this later
			}
			jobs = append(jobs, job)
		}
	})

	// Set up error handling
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error scraping %s: %v\n", site.URL, err)
	})

	// Visit the site
	err := c.Visit(site.URL)
	if err != nil {
		return &models.ScrapingResult{
			Site:     site.Name,
			Jobs:     jobs,
			Error:    err,
			Duration: time.Since(start),
		}, err
	}

	return &models.ScrapingResult{
		Site:     site.Name,
		Jobs:     jobs,
		Duration: time.Since(start),
	}, nil
}

// isRemoteJob checks if a job is remote based on various indicators
func (js *JobScraper) isRemoteJob(title, company, location, fullText string) bool {
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

	// If location is empty or very generic, might be remote
	if location == "" ||
		strings.Contains(locationLower, "united states") ||
		strings.Contains(locationLower, "usa") ||
		strings.Contains(locationLower, "us") {
		return true
	}

	return false
}

// isRelevantRole checks if a job title/description matches relevant engineering roles
func (js *JobScraper) isRelevantRole(title, fullText string) bool {
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

// ScrapeAPISite scrapes a single API site for job listings
func (js *JobScraper) ScrapeAPISite(ctx context.Context, apiSite config.APISite) (*models.ScrapingResult, error) {
	// Delegate to the API scraper
	apiScraper := NewAPIScraper()
	return apiScraper.ScrapeAPISite(ctx, apiSite)
}

// ScrapeAllSites scrapes all configured sites
func (js *JobScraper) ScrapeAllSites(ctx context.Context, sites []config.Site) ([]models.ScrapingResult, error) {
	var results []models.ScrapingResult

	for _, site := range sites {
		result, err := js.ScrapeSite(ctx, site)
		if err != nil {
			fmt.Printf("Error scraping %s: %v\n", site.Name, err)
		}
		results = append(results, *result)
	}

	return results, nil
}

// ScrapeAllAPISites scrapes all configured API sites
func (js *JobScraper) ScrapeAllAPISites(ctx context.Context, apiSites []config.APISite) ([]models.ScrapingResult, error) {
	var results []models.ScrapingResult

	for _, apiSite := range apiSites {
		result, err := js.ScrapeAPISite(ctx, apiSite)
		if err != nil {
			fmt.Printf("Error scraping API site %s: %v\n", apiSite.Name, err)
		}
		results = append(results, *result)
	}

	return results, nil
}

// extractJobDetails extracts job details from a job page
func (js *JobScraper) extractJobDetails(html string, site config.Site) (*models.JobListing, error) {
	// TODO: Implement this function
	// 1. Parse HTML using goquery
	// 2. Extract job details based on site type
	// 3. Handle different site formats (Indeed, LinkedIn, etc.)
	// 4. Return JobListing struct

	return &models.JobListing{}, nil
}

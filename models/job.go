package models

import (
	"time"
)

// JobListing represents a job posting
type JobListing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Company     string    `json:"company"`
	Location    string    `json:"location"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	Salary      string    `json:"salary"`
	PostedDate  time.Time `json:"posted_date"`
	Source      string    `json:"source"` // which site it came from

	// Relevance scoring
	RelevanceScore float64  `json:"relevance_score"`
	MatchedSkills  []string `json:"matched_skills"`

	// Metadata
	ScrapedAt time.Time `json:"scraped_at"`
}

// SeenJob represents a job that has been seen before
type SeenJob struct {
	JobID  string    `json:"job_id"`
	URL    string    `json:"url"`
	SeenAt time.Time `json:"seen_at"`
	Source string    `json:"source"`
}

// ScrapingResult represents the result of scraping a site
type ScrapingResult struct {
	Site     string        `json:"site"`
	Jobs     []JobListing  `json:"jobs"`
	Error    error         `json:"error,omitempty"`
	Duration time.Duration `json:"duration"`
}

// JobMatch represents a job that matches the user's criteria
type JobMatch struct {
	Job            JobListing `json:"job"`
	RelevanceScore float64    `json:"relevance_score"`
	MatchedSkills  []string   `json:"matched_skills"`
}

// ScrapingStats represents statistics from a scraping run
type ScrapingStats struct {
	TotalJobsFound    int           `json:"total_jobs_found"`
	RelevantJobsFound int           `json:"relevant_jobs_found"`
	NewJobsFound      int           `json:"new_jobs_found"`
	TotalDuration     time.Duration `json:"total_duration"`
	SitesScraped      int           `json:"sites_scraped"`
	SitesFailed       int           `json:"sites_failed"`
}

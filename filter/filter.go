package filter

import (
	"job-scraper/config"
	"job-scraper/models"
)

// JobFilter handles filtering and scoring of job listings
type JobFilter struct {
	config   *config.Config
	seenJobs map[string]models.SeenJob
}

// NewJobFilter creates a new job filter
func NewJobFilter(cfg *config.Config) *JobFilter {
	// TODO: Implement this function
	// 1. Load seen jobs from file if exists
	// 2. Initialize the filter with config
	// 3. Return the filter instance

	return &JobFilter{
		config:   cfg,
		seenJobs: make(map[string]models.SeenJob),
	}
}

// FilterJobs filters and scores job listings based on relevance
func (jf *JobFilter) FilterJobs(jobs []models.JobListing) []models.JobMatch {
	var matches []models.JobMatch

	// TODO: Implement this function
	// 1. For each job:
	//    - Check if already seen (skip if yes)
	//    - Check if job location matches criteria
	//    - Calculate relevance score
	//    - Check if score meets minimum threshold
	//    - If relevant, add to matches
	// 2. Return filtered matches

	return matches
}

// CalculateRelevanceScore calculates how relevant a job is based on keywords
func (jf *JobFilter) CalculateRelevanceScore(job models.JobListing) (float64, []string) {
	// TODO: Implement this function
	// 1. Combine job title, description, and company into searchable text
	// 2. Convert to lowercase for case-insensitive matching
	// 3. Count how many keywords match
	// 4. Calculate score (e.g., matches / total keywords)
	// 5. Return score and list of matched keywords

	return 0.0, []string{}
}

// IsJobSeen checks if a job has been seen before
func (jf *JobFilter) IsJobSeen(job models.JobListing) bool {
	// TODO: Implement this function
	// 1. Generate a unique ID for the job (URL hash or similar)
	// 2. Check if ID exists in seenJobs map
	// 3. Return true if seen, false otherwise

	return false
}

// MarkJobAsSeen marks a job as seen
func (jf *JobFilter) MarkJobAsSeen(job models.JobListing) {
	// TODO: Implement this function
	// 1. Generate unique ID for the job
	// 2. Add to seenJobs map with current timestamp
	// 3. Optionally save to file periodically
}

// SaveSeenJobs saves the seen jobs to file
func (jf *JobFilter) SaveSeenJobs() error {
	// TODO: Implement this function
	// 1. Marshal seenJobs to JSON
	// 2. Write to file specified in config
	// 3. Return error if any

	return nil
}

// LoadSeenJobs loads seen jobs from file
func (jf *JobFilter) LoadSeenJobs() error {
	// TODO: Implement this function
	// 1. Read file specified in config
	// 2. Unmarshal JSON to seenJobs map
	// 3. Return error if any

	return nil
}

// IsLocationMatch checks if a job's location matches the filtering criteria
func (jf *JobFilter) IsLocationMatch(job models.JobListing) bool {
	// TODO: Implement this function
	// 1. Check if global search is enabled and no exclusions match
	// 2. Check if job location is in allowed countries
	// 3. Check if job location is in allowed cities
	// 4. Check remote work preferences
	// 5. Return true if location matches criteria

	return true
}

// IsRemoteJob checks if a job is remote, hybrid, or on-site
func (jf *JobFilter) IsRemoteJob(job models.JobListing) (bool, bool) {
	// TODO: Implement this function
	// 1. Parse job location for remote keywords ("remote", "work from home", etc.)
	// 2. Check if it's hybrid (mix of remote and on-site)
	// 3. Return (isRemote, isHybrid)

	return false, false
}

// IsCountryAllowed checks if a country is in the allowed list
func (jf *JobFilter) IsCountryAllowed(country string) bool {
	// TODO: Implement this function
	// 1. If global search is enabled, check exclusions
	// 2. If specific countries are listed, check if country is in list
	// 3. Return true if country is allowed

	return true
}

// ParseJobLocation extracts country and city information from job location string
func (jf *JobFilter) ParseJobLocation(location string) (string, string, bool) {
	// TODO: Implement this function
	// 1. Parse location string for country and city
	// 2. Handle common formats like "City, Country", "Remote", "City, State, Country"
	// 3. Return (country, city, isRemote)
	// 4. Use string manipulation and regex if needed

	return "", "", false
}

// generateJobID generates a unique ID for a job
func (jf *JobFilter) generateJobID(job models.JobListing) string {
	// TODO: Implement this function
	// 1. Create a hash of job URL or title+company
	// 2. Return the hash as string

	return ""
}

package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the main configuration for the job scraper
type Config struct {
	// Sites to scrape for job listings
	Sites []Site `json:"sites"`

	// Resume keywords or skills to match against
	ResumeKeywords []string `json:"resume_keywords"`

	// Location filtering configuration
	Location LocationConfig `json:"location"`

	// Email configuration for notifications
	Email EmailConfig `json:"email"`

	// Minimum relevance score to consider a job (0.0 to 1.0)
	MinRelevanceScore float64 `json:"min_relevance_score"`

	// Path to store seen jobs (to avoid duplicates)
	SeenJobsPath string `json:"seen_jobs_path"`

	// Resume file path (optional, for parsing resume)
	ResumePath string `json:"resume_path"`
}

// Site represents a job site to scrape
type Site struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Type     string `json:"type"`     // "indeed", "linkedin", "glassdoor", etc.
	Selector string `json:"selector"` // CSS selector for job listings
}

// LocationConfig represents location filtering settings
type LocationConfig struct {
	// Global search - if true, search worldwide
	Global bool `json:"global"`

	// Specific countries to include (ISO country codes like "US", "CA", "GB")
	Countries []string `json:"countries"`

	// Specific cities to include
	Cities []string `json:"cities"`

	// Remote work preference
	Remote RemoteConfig `json:"remote"`

	// Exclude specific countries (useful when using global search)
	ExcludeCountries []string `json:"exclude_countries"`
}

// RemoteConfig represents remote work preferences
type RemoteConfig struct {
	// Accept remote jobs
	Accept bool `json:"accept"`

	// Require remote jobs only
	Required bool `json:"required"`

	// Accept hybrid jobs (mix of remote and on-site)
	Hybrid bool `json:"hybrid"`
}

// EmailConfig represents email notification settings
type EmailConfig struct {
	Enabled  bool   `json:"enabled"`
	SMTPHost string `json:"smtp_host"`
	SMTPPort int    `json:"smtp_port"`
	Username string `json:"username"`
	Password string `json:"password"`
	To       string `json:"to"`
	From     string `json:"from"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(configPath string) (*Config, error) {
	// 1. Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// 2. Parse JSON into Config struct
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	// 3. Validate the configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// 4. Debug logging (if enabled)
	if os.Getenv("JOB_SCRAPER_DEBUG") == "true" {
		fmt.Printf("DEBUG: Config loaded successfully\n")
		fmt.Printf("DEBUG: Sites: %d, Keywords: %d, Email enabled: %v\n",
			len(cfg.Sites), len(cfg.ResumeKeywords), cfg.Email.Enabled)
	}

	// 5. Return the config
	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// TODO: Implement this function
	// 1. Check if at least one site is configured
	if len(c.Sites) == 0 {
		return fmt.Errorf("config.json requires at least 1 site configured")
	}
	// 2. Check if at least one keyword is provided
	if len(c.ResumeKeywords) == 0 {
		return fmt.Errorf("config.json requires at least 1 resumeKeyword configured")
	}
	// 3. Validate email settings if enabled
	if c.Email.Enabled {
		// add validation
	}
	// 4. Validate location settings
	// 5. Return error if validation fails

	return nil
}

// SaveConfig saves configuration to a JSON file
// Nice to have - leaving as placeholder
func SaveConfig(config *Config, configPath string) error {
	// TODO: Implement this function
	// 1. Marshal Config struct to JSON

	// 2. Write to file
	// 3. Return error if any

	return nil
}

package notifier

import (
	"fmt"

	"job-scraper/config"
	"job-scraper/models"
)

// Notifier handles sending notifications about job matches
type Notifier struct {
	config *config.Config
}

// NewNotifier creates a new notifier
func NewNotifier(cfg *config.Config) *Notifier {
	return &Notifier{
		config: cfg,
	}
}

// NotifyJobMatches sends notifications about new job matches
func (n *Notifier) NotifyJobMatches(matches []models.JobMatch, stats models.ScrapingStats) error {
	if len(matches) == 0 {
		// TODO: Optionally send "no matches" notification
		return nil
	}

	// TODO: Implement this function
	// 1. Generate email content from matches
	// 2. If email is configured, send email
	// 3. Otherwise, print to console or save to file
	// 4. Return error if any

	return nil
}

// SendEmail sends an email with job matches
func (n *Notifier) SendEmail(subject, body string) error {
	if !n.config.Email.Enabled {
		return fmt.Errorf("email notifications are disabled")
	}

	// TODO: Implement this function
	// 1. Create gomail.Message
	// 2. Set headers (From, To, Subject)
	// 3. Set body (HTML or plain text)
	// 4. Create dialer with SMTP settings
	// 5. Send the email
	// 6. Return error if any

	return nil
}

// GenerateEmailContent generates HTML content for job matches email
func (n *Notifier) GenerateEmailContent(matches []models.JobMatch, stats models.ScrapingStats) (string, error) {
	// TODO: Implement this function
	// 1. Create HTML template for job matches
	// 2. Include job details, relevance scores, links
	// 3. Add statistics about the scraping run
	// 4. Return formatted HTML string

	return "", nil
}

// SaveToFile saves job matches to a text file
func (n *Notifier) SaveToFile(matches []models.JobMatch, filename string) error {
	// TODO: Implement this function
	// 1. Open file for writing
	// 2. Format job matches as text
	// 3. Write to file
	// 4. Close file and return error if any

	return nil
}

// PrintToConsole prints job matches to console
func (n *Notifier) PrintToConsole(matches []models.JobMatch, stats models.ScrapingStats) {
	// TODO: Implement this function
	// 1. Print header with statistics
	// 2. For each match, print formatted job details
	// 3. Include relevance score and matched skills
	// 4. Print job URL for easy access
}

// FormatJobMatch formats a single job match for display
func (n *Notifier) FormatJobMatch(match models.JobMatch) string {
	// TODO: Implement this function
	// 1. Format job title, company, location
	// 2. Include relevance score
	// 3. Include matched skills
	// 4. Include job URL
	// 5. Return formatted string

	return ""
}

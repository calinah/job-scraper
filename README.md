# Job Scraper

A Go-based job scraping application that automatically finds and filters job listings based on your resume and preferences.

## Features

- **Multi-site scraping**: Scrape job listings from multiple job sites (Indeed, LinkedIn, etc.)
- **Resume-based filtering**: Extract keywords from your resume or use predefined keywords
- **Location filtering**: Filter jobs by global search, specific countries, cities, and remote work preferences
- **Relevance scoring**: Score jobs based on keyword matches and relevance
- **Duplicate detection**: Avoid re-notifying about jobs you've already seen
- **Email notifications**: Get notified about new relevant jobs via email
- **Configurable**: Easy configuration via JSON file

## Project Structure

```
job-scraper/
├── config/          # Configuration management
├── filter/          # Job filtering and relevance scoring
├── models/          # Data structures and types
├── notifier/        # Email and notification system
├── resume/          # Resume parsing and keyword extraction
├── scraper/         # Web scraping functionality
├── main.go          # Application entry point
├── config.json      # Configuration file
└── README.md        # This file
```


## Setup

1. **Install Go** (version 1.21 or later)
2. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd job-scraper
   ```

3. **Install dependencies**:
   ```bash
   go mod tidy
   ```

4. **Configure the application**:
   - Edit `config.json` with your preferences
   - Add your resume file (optional)
   - Configure email settings (optional)

## Usage

Run the application:

```bash
go run main.go
```

Or build and run:

```bash
go build -o job-scraper
./job-scraper
```

## Implementation Guide

### Phase 1: Configuration and Basic Structure

Start by implementing the configuration system:

1. **Implement `config.LoadConfig()`**:
   - Read JSON file using `os.ReadFile()`
   - Unmarshal JSON into Config struct
   - Add validation for required fields

2. **Implement `config.SaveConfig()`**:
   - Marshal Config struct to JSON
   - Write to file using `os.WriteFile()`

### Phase 2: Data Models and Persistence

1. **Implement seen jobs persistence**:
   - Create functions to save/load seen jobs
   - Use JSON for storage
   - Handle file not found gracefully

2. **Implement job ID generation**:
   - Create unique IDs for jobs (hash of URL or title+company)
   - Use crypto/sha256 for hashing

### Phase 3: Web Scraping

1. **Implement basic scraping**:
   - Use `github.com/gocolly/colly` for web scraping
   - Start with a simple site (like a job board)
   - Extract basic job information

2. **Handle different site types**:
   - Create site-specific scrapers
   - Handle different HTML structures
   - Add rate limiting and delays

### Phase 4: Filtering and Scoring

1. **Implement relevance scoring**:
   - Count keyword matches in job description
   - Calculate score based on match percentage
   - Consider job title weight more heavily

2. **Implement filtering logic**:
   - Check minimum relevance score
   - Filter out already seen jobs
   - Sort by relevance score

### Phase 5: Notifications

1. **Implement email notifications**:
   - Use `gopkg.in/gomail.v2` for sending emails
   - Create HTML templates for job listings
   - Handle email configuration

2. **Add console output**:
   - Print job matches to console
   - Show statistics and progress

### Phase 6: Resume Parsing

1. **Implement PDF parsing**:
   - Use a PDF library to extract text
   - Parse for common programming skills
   - Extract keywords automatically

2. **Add keyword extraction**:
   - Define common programming languages/frameworks
   - Search for these in resume text
   - Combine with manual keywords

## Configuration

Edit `config.json` to customize:

- **Sites to scrape**: Add job sites with URLs and selectors
- **Keywords**: Add skills and technologies to match
- **Location filtering**: Configure global search, countries, cities, and remote work preferences
- **Email settings**: Configure SMTP for notifications
- **Relevance threshold**: Set minimum score for job matches

### Location Configuration Options

The `location` section in `config.json` supports:

- **`global`**: Set to `true` for worldwide search
- **`countries`**: Array of ISO country codes (e.g., "US", "CA", "GB")
- **`cities`**: Array of specific cities to include
- **`remote`**: Remote work preferences:
  - `accept`: Accept remote jobs
  - `required`: Require remote jobs only
  - `hybrid`: Accept hybrid jobs
- **`exclude_countries`**: Countries to exclude (useful with global search)

### Example Location Configurations

**Global search with exclusions:**
```json
{
  "location": {
    "global": true,
    "exclude_countries": ["CN", "RU"]
  }
}
```

**Specific countries only:**
```json
{
  "location": {
    "global": false,
    "countries": ["US", "CA", "GB"],
    "cities": ["San Francisco", "New York", "Toronto"]
  }
}
```

**Remote work only:**
```json
{
  "location": {
    "global": true,
    "remote": {
      "accept": true,
      "required": true,
      "hybrid": false
    }
  }
}
```

## Next Steps

1. **Start with configuration**: Implement `config.LoadConfig()`
2. **Add basic scraping**: Implement one site scraper
3. **Add filtering**: Implement relevance scoring
4. **Add persistence**: Save seen jobs
5. **Add notifications**: Email or console output
6. **Improve and expand**: Add more sites, better parsing, etc.


## Dependencies

- `github.com/gocolly/colly`: Web scraping
- `github.com/PuerkitoBio/goquery`: HTML parsing
- `gopkg.in/gomail.v2`: Email sending

## License

MIT License - feel free to use and modify as needed for learning! 
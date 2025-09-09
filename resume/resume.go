package resume

import (
	"job-scraper/config"
)

// ResumeParser handles parsing resume files to extract keywords
type ResumeParser struct {
	config *config.Config
}

// NewResumeParser creates a new resume parser
func NewResumeParser(cfg *config.Config) *ResumeParser {
	return &ResumeParser{
		config: cfg,
	}
}

// ExtractKeywords extracts keywords from a resume file
func (rp *ResumeParser) ExtractKeywords(resumePath string) ([]string, error) {
	// TODO: Implement this function
	// 1. Read the resume file (PDF, DOCX, or TXT)
	// 2. Extract text content
	// 3. Parse for skills, technologies, keywords
	// 4. Return list of keywords

	return []string{}, nil
}

// ParsePDF extracts text from a PDF resume
func (rp *ResumeParser) ParsePDF(filePath string) (string, error) {
	// TODO: Implement this function
	// 1. Use a PDF parsing library (like unidoc/unipdf)
	// 2. Extract text content from PDF
	// 3. Clean and normalize text
	// 4. Return text content

	return "", nil
}

// ParseDOCX extracts text from a DOCX resume
func (rp *ResumeParser) ParseDOCX(filePath string) (string, error) {
	// TODO: Implement this function
	// 1. Use a DOCX parsing library
	// 2. Extract text content from DOCX
	// 3. Clean and normalize text
	// 4. Return text content

	return "", nil
}

// ExtractSkillsFromText extracts skills and keywords from text
func (rp *ResumeParser) ExtractSkillsFromText(text string) []string {
	// TODO: Implement this function
	// 1. Define common programming languages, frameworks, tools
	// 2. Search for these terms in the text
	// 3. Also look for common skill patterns
	// 4. Return unique list of found skills

	return []string{}
}

// LoadKeywordsFromFile loads keywords from a simple text file
func (rp *ResumeParser) LoadKeywordsFromFile(filePath string) ([]string, error) {
	// TODO: Implement this function
	// 1. Read file line by line
	// 2. Skip empty lines and comments
	// 3. Trim whitespace from each line
	// 4. Return list of keywords

	return []string{}, nil
}

// GetKeywords returns the keywords to use for job matching
func (rp *ResumeParser) GetKeywords() ([]string, error) {
	// TODO: Implement this function
	// 1. If resume path is provided, extract keywords from resume
	// 2. If keywords are provided in config, use those
	// 3. Combine both sources if available
	// 4. Return final keyword list

	return rp.config.ResumeKeywords, nil
}

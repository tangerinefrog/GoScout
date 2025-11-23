package validator

import (
	"strings"
)

type keywordValidator struct {
}

func NewKeywordValidator() *keywordValidator {
	return &keywordValidator{}
}

func (v *keywordValidator) ValidateKeywords(keywords []string, text string) bool {
	text = strings.ToLower(strings.TrimSpace(text))
	for _, keyword := range keywords {
		keyword = strings.ToLower(strings.TrimSpace(keyword))
		if keyword == "" {
			continue
		}
		matched := strings.Contains(text, keyword)
		if matched {
			return true
		}
	}

	return false
}

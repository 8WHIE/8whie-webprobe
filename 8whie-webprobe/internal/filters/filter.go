package filters

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/8whie/8whie-webprobe/pkg/model"
)

// ResponseFilter evaluates probe outcomes against configured match and filter rules.
type ResponseFilter struct {
	criteria model.FilterCriteria
}

// NewResponseFilter creates a filter evaluator from filter criteria.
func NewResponseFilter(criteria model.FilterCriteria) *ResponseFilter {
	return &ResponseFilter{
		criteria: criteria,
	}
}

// Evaluate determines whether a probe response satisfies matching thresholds.
func (f *ResponseFilter) Evaluate(res *model.ProbeResponse) (bool, string) {
	if res == nil {
		return false, "nil response"
	}

	// 1. Status Code Matching
	if len(f.criteria.MatchStatusCodes) > 0 {
		matched := false
		for _, code := range f.criteria.MatchStatusCodes {
			if res.StatusCode == code {
				matched = true
				break
			}
		}
		if !matched {
			return false, fmt.Sprintf("status code %d not in match list", res.StatusCode)
		}
	}

	// 2. Status Code Filtering
	if len(f.criteria.FilterStatusCodes) > 0 {
		for _, code := range f.criteria.FilterStatusCodes {
			if res.StatusCode == code {
				return false, fmt.Sprintf("status code %d excluded by filter", res.StatusCode)
			}
		}
	}

	// 3. Response Size Matching
	if len(f.criteria.MatchSizes) > 0 {
		matched := false
		for _, sz := range f.criteria.MatchSizes {
			if res.ContentLength == sz {
				matched = true
				break
			}
		}
		if !matched {
			return false, fmt.Sprintf("size %d not in match list", res.ContentLength)
		}
	}

	// 4. Response Size Filtering
	if len(f.criteria.FilterSizes) > 0 {
		for _, sz := range f.criteria.FilterSizes {
			if res.ContentLength == sz {
				return false, fmt.Sprintf("size %d excluded by filter", res.ContentLength)
			}
		}
	}

	// 5. Word Count Matching
	if len(f.criteria.MatchWords) > 0 {
		matched := false
		for _, w := range f.criteria.MatchWords {
			if res.WordCount == w {
				matched = true
				break
			}
		}
		if !matched {
			return false, fmt.Sprintf("word count %d not in match list", res.WordCount)
		}
	}

	// 6. Word Count Filtering
	if len(f.criteria.FilterWords) > 0 {
		for _, w := range f.criteria.FilterWords {
			if res.WordCount == w {
				return false, fmt.Sprintf("word count %d excluded by filter", res.WordCount)
			}
		}
	}

	// 7. Line Count Matching
	if len(f.criteria.MatchLines) > 0 {
		matched := false
		for _, l := range f.criteria.MatchLines {
			if res.LineCount == l {
				matched = true
				break
			}
		}
		if !matched {
			return false, fmt.Sprintf("line count %d not in match list", res.LineCount)
		}
	}

	// 8. Line Count Filtering
	if len(f.criteria.FilterLines) > 0 {
		for _, l := range f.criteria.FilterLines {
			if res.LineCount == l {
				return false, fmt.Sprintf("line count %d excluded by filter", res.LineCount)
			}
		}
	}

	return true, "matched"
}

// CountWordsAndLines calculates word and line tallies from raw response bytes.
func CountWordsAndLines(data []byte) (wordCount int, lineCount int) {
	if len(data) == 0 {
		return 0, 0
	}

	// Count lines
	lineCount = bytes.Count(data, []byte("\n"))
	if !bytes.HasSuffix(data, []byte("\n")) {
		lineCount++
	}

	// Count words
	inWord := false
	for _, b := range data {
		r := rune(b)
		if unicode.IsSpace(r) {
			if inWord {
				inWord = false
			}
		} else {
			if !inWord {
				inWord = true
				wordCount++
			}
		}
	}

	return wordCount, lineCount
}

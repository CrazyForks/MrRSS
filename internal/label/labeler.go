package label

import (
	"sort"
	"strings"
)

// Labeler provides local label generation using NLP algorithms
type Labeler struct{}

// NewLabeler creates a new Labeler instance
func NewLabeler() *Labeler {
	return &Labeler{}
}

// GenerateLabels generates labels from the given text using TF-IDF and keyword extraction
func (l *Labeler) GenerateLabels(text string, maxLabels int) LabelResult {
	// Clean the text
	cleanedText := cleanText(text)

	// Check if text is too short
	if len(cleanedText) < MinContentLength {
		return LabelResult{
			Labels:     []string{},
			IsTooShort: true,
		}
	}

	// Validate maxLabels parameter
	if maxLabels <= 0 || maxLabels > MaxLabels {
		maxLabels = MaxLabels
	}

	// Check if text is primarily Chinese
	isChinese := isChineseText(cleanedText)

	// Extract keywords/phrases
	var labels []string
	if isChinese {
		labels = l.extractChineseLabels(cleanedText, maxLabels)
	} else {
		labels = l.extractEnglishLabels(cleanedText, maxLabels)
	}

	// Normalize and validate labels
	var validLabels []string
	for _, label := range labels {
		normalized := normalizeLabel(label)
		if validateLabel(normalized) {
			validLabels = append(validLabels, normalized)
		}
		if len(validLabels) >= maxLabels {
			break
		}
	}

	return LabelResult{
		Labels:     validLabels,
		IsTooShort: false,
	}
}

// extractEnglishLabels extracts labels from English text using word frequency
func (l *Labeler) extractEnglishLabels(text string, maxLabels int) []string {
	words := splitWords(strings.ToLower(text), false)

	// Count word frequencies
	wordFreq := make(map[string]int)
	for _, word := range words {
		wordFreq[word]++
	}

	// Calculate TF scores
	type scoredWord struct {
		word  string
		score float64
	}
	var scored []scoredWord

	totalWords := float64(len(words))
	for word, freq := range wordFreq {
		// TF (Term Frequency) = frequency of word / total words
		tf := float64(freq) / totalWords

		// Simple scoring: prefer words that appear multiple times but not too frequently
		// Boost words that appear 2-10 times
		score := tf
		if freq >= 2 && freq <= 10 {
			score *= 1.5
		}

		scored = append(scored, scoredWord{word: word, score: score})
	}

	// Sort by score descending
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Extract top labels
	var labels []string
	for i := 0; i < len(scored) && len(labels) < maxLabels; i++ {
		labels = append(labels, scored[i].word)
	}

	return labels
}

// extractChineseLabels extracts labels from Chinese text using n-gram frequency
func (l *Labeler) extractChineseLabels(text string, maxLabels int) []string {
	ngrams := splitWords(text, true)

	// Count n-gram frequencies
	ngramFreq := make(map[string]int)
	for _, ngram := range ngrams {
		ngramFreq[ngram]++
	}

	// Calculate scores
	type scoredNGram struct {
		ngram string
		score float64
	}
	var scored []scoredNGram

	totalNGrams := float64(len(ngrams))
	for ngram, freq := range ngramFreq {
		// Prefer trigrams over bigrams
		lengthBoost := 1.0
		if len([]rune(ngram)) == 3 {
			lengthBoost = 1.5
		}

		// TF with length boost
		tf := float64(freq) / totalNGrams
		score := tf * lengthBoost

		// Boost n-grams that appear 2-8 times
		if freq >= 2 && freq <= 8 {
			score *= 1.5
		}

		scored = append(scored, scoredNGram{ngram: ngram, score: score})
	}

	// Sort by score descending
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Extract top labels, preferring longer n-grams
	var labels []string
	seen := make(map[string]bool)

	// First, collect trigrams
	for _, item := range scored {
		if len([]rune(item.ngram)) == 3 && !seen[item.ngram] {
			labels = append(labels, item.ngram)
			seen[item.ngram] = true
			if len(labels) >= maxLabels {
				return labels
			}
		}
	}

	// Then fill with bigrams if needed
	for _, item := range scored {
		if len([]rune(item.ngram)) == 2 && !seen[item.ngram] {
			labels = append(labels, item.ngram)
			seen[item.ngram] = true
			if len(labels) >= maxLabels {
				break
			}
		}
	}

	return labels
}

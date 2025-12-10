// Package label provides automatic label generation for articles.
// It supports both local NLP algorithms and AI-based labeling.
package label

// LabelResult contains the generated labels and related metadata
type LabelResult struct {
	Labels     []string `json:"labels"`      // Generated labels
	IsTooShort bool     `json:"is_too_short"` // True if content is too short to generate labels
}

// MinContentLength is the minimum content length required for label generation (in characters)
const MinContentLength = 50

// MaxLabels is the maximum number of labels to generate
const MaxLabels = 5

// MinLabelLength is the minimum length for a label (in characters)
const MinLabelLength = 2

// MaxLabelLength is the maximum length for a label (in characters)
const MaxLabelLength = 30

// MaxInputCharsForAI is the maximum input size for AI labeling to save tokens
const MaxInputCharsForAI = 5000

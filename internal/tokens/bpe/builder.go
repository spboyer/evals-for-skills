package tokenizer

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"strings"
)

var modelPrefixToEncoding = []struct {
	Prefix   string
	Encoding string
}{
	{Prefix: "gpt-5", Encoding: "o200k_base"},
	{Prefix: "gpt-4.", Encoding: "o200k_base"},
}

// ModelToEncoding maps model names to the tokenizer encoding used by the model.
var ModelToEncoding = map[string]string{
	"gpt-4o":                       "o200k_base",
	"gpt-4":                        "cl100k_base",
	"gpt-3.5-turbo":                "cl100k_base",
	"text-davinci-003":             "p50k_base",
	"text-davinci-002":             "p50k_base",
	"text-davinci-001":             "r50k_base",
	"text-curie-001":               "r50k_base",
	"text-babbage-001":             "r50k_base",
	"text-ada-001":                 "r50k_base",
	"davinci":                      "r50k_base",
	"curie":                        "r50k_base",
	"babbage":                      "r50k_base",
	"ada":                          "r50k_base",
	"code-davinci-002":             "p50k_base",
	"code-davinci-001":             "p50k_base",
	"code-cushman-002":             "p50k_base",
	"code-cushman-001":             "p50k_base",
	"davinci-codex":                "p50k_base",
	"cushman-codex":                "p50k_base",
	"text-davinci-edit-001":        "p50k_edit",
	"code-davinci-edit-001":        "p50k_edit",
	"text-embedding-ada-002":       "cl100k_base",
	"text-similarity-davinci-001":  "r50k_base",
	"text-similarity-curie-001":    "r50k_base",
	"text-similarity-babbage-001":  "r50k_base",
	"text-similarity-ada-001":      "r50k_base",
	"text-search-davinci-doc-001":  "r50k_base",
	"text-search-curie-doc-001":    "r50k_base",
	"text-search-babbage-doc-001":  "r50k_base",
	"text-search-ada-doc-001":      "r50k_base",
	"code-search-babbage-code-001": "r50k_base",
	"code-search-ada-code-001":     "r50k_base",
	"gpt2":                         "gpt2",
}

const (
	endOfText   = "<|endoftext|>"
	endOfPrompt = "<|endofprompt|>"
)

// Regex patterns mirrored from tokenizer_ts with stdlib-compatible whitespace handling.
const (
	regexPatternLegacy = `'s|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+`
	regexPatternModern = `(?:'s|'S|'t|'T|'re|'RE|'Re|'eR|'ve|'VE|'vE|'Ve|'m|'M|'ll|'lL|'Ll|'LL|'d|'D)|[^\r\n\p{L}\p{N}]?\p{L}+|\p{N}{1,3}| ?[^\s\p{L}\p{N}]+[\r\n]*|\s*[\r\n]+|\s+`
)

var regexPatternO200k = strings.Join([]string{
	`[^\r\n\p{L}\p{N}]?[\p{Lu}\p{Lt}\p{Lm}\p{Lo}\p{M}]*[\p{Ll}\p{Lm}\p{Lo}\p{M}]+(?:'s|'S|'t|'T|'re|'RE|'Re|'eR|'ve|'VE|'vE|'Ve|'m|'M|'ll|'lL|'Ll|'LL|'d|'D)?`,
	`[^\r\n\p{L}\p{N}]?[\p{Lu}\p{Lt}\p{Lm}\p{Lo}\p{M}]+[\p{Ll}\p{Lm}\p{Lo}\p{M}]*(?:'s|'S|'t|'T|'re|'RE|'Re|'eR|'ve|'VE|'vE|'Ve|'m|'M|'ll|'lL|'Ll|'LL|'d|'D)?`,
	`\p{N}{1,3}`,
	` ?[^\s\p{L}\p{N}]+[\r\n/]*`,
	`\s*[\r\n]+`,
	`\s+`,
}, "|")

//go:embed model
var modelFS embed.FS

func modelDir() fs.FS {
	sub, err := fs.Sub(modelFS, "model")
	if err != nil {
		panic("embedded model directory missing: " + err.Error())
	}
	return sub
}

func encodingForModelName(modelName string) string {
	if encoding, ok := ModelToEncoding[modelName]; ok {
		return encoding
	}
	for _, entry := range modelPrefixToEncoding {
		if strings.HasPrefix(modelName, entry.Prefix) {
			return entry.Encoding
		}
	}
	return ""
}

func mergeSpecialTokens(base map[string]int, extra map[string]int) map[string]int {
	out := cloneStringIntMap(base)
	for token, rank := range extra {
		out[token] = rank
	}
	return out
}

func bpeFileForEncoding(encoding string) (string, string, error) {
	switch encoding {
	case "o200k_base":
		return regexPatternO200k, "o200k_base.tiktoken", nil
	default:
		return "", "", errors.New(encoding + " encoding isn't supported")
	}
}

func SpecialTokensForEncoding(encoding string) map[string]int {
	specialTokens := map[string]int{endOfText: 50256}
	switch encoding {
	case "o200k_base":
		specialTokens = map[string]int{
			endOfText:   199999,
			endOfPrompt: 200018,
		}
	}
	return specialTokens
}

func SpecialTokensForModel(modelName string) map[string]int {
	return SpecialTokensForEncoding(encodingForModelName(modelName))
}

func RegexForEncoding(encoding string) string {
	switch encoding {
	case "o200k_base":
		return regexPatternO200k
	default:
		return regexPatternLegacy
	}
}

func RegexForModel(modelName string) string {
	return RegexForEncoding(encodingForModelName(modelName))
}

func EncodingForModel(modelName string) (string, error) {
	encoding := encodingForModelName(modelName)
	if encoding == "" {
		return "", fmt.Errorf("doesn't support this model [%s]", modelName)
	}
	return encoding, nil
}

func NewTokenizerForModel(modelName string, extraSpecialTokens map[string]int) (*Tokenizer, error) {
	encoding, err := EncodingForModel(modelName)
	if err != nil {
		return nil, err
	}
	return NewTokenizerForEncoding(encoding, extraSpecialTokens)
}

func NewTokenizerForEncoding(encoding string, extraSpecialTokens map[string]int) (*Tokenizer, error) {
	regexPattern, fileName, err := bpeFileForEncoding(encoding)
	if err != nil {
		return nil, err
	}

	specialTokens := SpecialTokensForEncoding(encoding)
	if extraSpecialTokens != nil {
		specialTokens = mergeSpecialTokens(specialTokens, extraSpecialTokens)
	}

	f, err := modelDir().Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open embedded model file %s: %w", fileName, err)
	}
	defer f.Close()

	return NewTokenizerFromReader(f, specialTokens, regexPattern, defaultCacheSize)
}

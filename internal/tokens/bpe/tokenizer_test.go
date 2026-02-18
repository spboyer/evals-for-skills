package tokenizer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var (
	imStart = "<|im_start|>"
	imEnd   = "<|im_end|>"
)

func readTestData(t *testing.T, fileName string) []byte {
	t.Helper()
	path := filepath.Join("testdata", fileName)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	return data
}

func readJSONTokenIDs(t *testing.T, fileName string) []int {
	t.Helper()
	var ids []int
	if err := json.Unmarshal(readTestData(t, fileName), &ids); err != nil {
		t.Fatalf("failed to parse %s: %v", fileName, err)
	}
	return ids
}

func requireEqualTokenIDs(t *testing.T, got, want []int) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("token mismatch\nwant: %v\ngot:  %v", want, got)
	}
}

func newTokenizerForModel(t *testing.T, model string, special map[string]int) *Tokenizer {
	t.Helper()
	tokenizer, err := NewTokenizerForModel(model, special)
	if err != nil {
		t.Fatalf("failed to create tokenizer for model %s: %v", model, err)
	}
	return tokenizer
}

func newTokenizerForEncoding(t *testing.T, encoding string, special map[string]int) *Tokenizer {
	t.Helper()
	tokenizer, err := NewTokenizerForEncoding(encoding, special)
	if err != nil {
		t.Fatalf("failed to create tokenizer for encoding %s: %v", encoding, err)
	}
	return tokenizer
}

func TestTikTokenizerCore(t *testing.T) {
	specialTokens := map[string]int{
		imStart: 100264,
		imEnd:   100265,
	}
	tokenizer := newTokenizerForModel(t, "gpt-3.5-turbo", specialTokens)

	t.Run("single punctuation", func(t *testing.T) {
		str := "!"
		encoded := tokenizer.Encode(str, nil)
		requireEqualTokenIDs(t, encoded, []int{0})
		if decoded := tokenizer.Decode(encoded); decoded != str {
			t.Fatalf("expected %q got %q", str, decoded)
		}
	})

	t.Run("empty string", func(t *testing.T) {
		str := ""
		encoded := tokenizer.Encode(str, nil)
		requireEqualTokenIDs(t, encoded, []int{})
		if decoded := tokenizer.Decode(encoded); decoded != str {
			t.Fatalf("expected %q got %q", str, decoded)
		}
	})

	t.Run("hello world", func(t *testing.T) {
		str := "Hello World!"
		encoded := tokenizer.Encode(str, nil)
		requireEqualTokenIDs(t, encoded, []int{9906, 4435, 0})
		if decoded := tokenizer.Decode(encoded); decoded != str {
			t.Fatalf("expected %q got %q", str, decoded)
		}
	})

	t.Run("special tokens 1", func(t *testing.T) {
		str := "<|im_start|>Hello World<|im_end|>"
		encoded := tokenizer.Encode(str, []string{imStart, imEnd})
		requireEqualTokenIDs(t, encoded, []int{100264, 9906, 4435, 100265})
		if decoded := tokenizer.Decode(encoded); decoded != str {
			t.Fatalf("expected %q got %q", str, decoded)
		}
	})

	t.Run("special tokens 2", func(t *testing.T) {
		str := "<|im_start|>Hello<|im_end|> World"
		encoded := tokenizer.Encode(str, []string{imStart, imEnd})
		requireEqualTokenIDs(t, encoded, []int{100264, 9906, 100265, 4435})
		if decoded := tokenizer.Decode(encoded); decoded != str {
			t.Fatalf("expected %q got %q", str, decoded)
		}
	})

	t.Run("special tokens unicode", func(t *testing.T) {
		str := "<|im_start|>Hello ⭐ World<|im_end|>"
		encoded := tokenizer.Encode(str, []string{imStart, imEnd})
		requireEqualTokenIDs(t, encoded, []int{100264, 9906, 2928, 99834, 4435, 100265})
		if decoded := tokenizer.Decode(encoded); decoded != str {
			t.Fatalf("expected %q got %q", str, decoded)
		}
	})

	t.Run("encode trim suffix", func(t *testing.T) {
		str := "<|im_start|>Hello World<|im_end|>"
		encodedStr := "<|im_start|>Hello World"

		encoded := tokenizer.EncodeTrimSuffix(str, 4, []string{imStart, imEnd})
		requireEqualTokenIDs(t, encoded.TokenIDs, []int{100264, 9906, 4435, 100265})
		if encoded.Text != str {
			t.Fatalf("expected %q got %q", str, encoded.Text)
		}

		encoded = tokenizer.EncodeTrimSuffix(str, 5, []string{imStart, imEnd})
		requireEqualTokenIDs(t, encoded.TokenIDs, []int{100264, 9906, 4435, 100265})
		if encoded.Text != str {
			t.Fatalf("expected %q got %q", str, encoded.Text)
		}

		encoded = tokenizer.EncodeTrimSuffix(str, 3, []string{imStart, imEnd})
		requireEqualTokenIDs(t, encoded.TokenIDs, []int{100264, 9906, 4435})
		if encoded.Text != encodedStr {
			t.Fatalf("expected %q got %q", encodedStr, encoded.Text)
		}
	})

	t.Run("encode trim suffix 2", func(t *testing.T) {
		str := "<|im_start|>Hello TempWorld<|im_end|>"
		encodedStr := "<|im_start|>Hello TempWorld"

		encoded := tokenizer.EncodeTrimSuffix(str, 5, []string{imStart, imEnd})
		requireEqualTokenIDs(t, encoded.TokenIDs, []int{100264, 9906, 20539, 10343, 100265})
		if encoded.Text != str {
			t.Fatalf("expected %q got %q", str, encoded.Text)
		}

		encoded = tokenizer.EncodeTrimSuffix(str, 6, []string{imStart, imEnd})
		requireEqualTokenIDs(t, encoded.TokenIDs, []int{100264, 9906, 20539, 10343, 100265})
		if encoded.Text != str {
			t.Fatalf("expected %q got %q", str, encoded.Text)
		}

		encoded = tokenizer.EncodeTrimSuffix(str, 3, []string{imStart, imEnd})
		requireEqualTokenIDs(t, encoded.TokenIDs, []int{100264, 9906, 20539})
		if encoded.Text != encodedStr {
			t.Fatalf("expected %q got %q", encodedStr, encoded.Text)
		}
	})

	t.Run("encode trim suffix 3", func(t *testing.T) {
		str := strings.Repeat("t", 4000)
		encoded := tokenizer.Encode(str, nil)
		trimmed := tokenizer.EncodeTrimSuffix(str, 5, []string{})
		if len(trimmed.TokenIDs) != 5 {
			t.Fatalf("expected 5 token ids, got %d", len(trimmed.TokenIDs))
		}
		requireEqualTokenIDs(t, trimmed.TokenIDs, encoded[:5])
	})

	t.Run("encode trim prefix", func(t *testing.T) {
		str := "<|im_start|>Hello World<|im_end|>"
		encodedStr := "Hello World<|im_end|>"

		encoded := tokenizer.EncodeTrimPrefix(str, 4, []string{imStart, imEnd})
		requireEqualTokenIDs(t, encoded.TokenIDs, []int{100264, 9906, 4435, 100265})
		if encoded.Text != str {
			t.Fatalf("expected %q got %q", str, encoded.Text)
		}

		encoded = tokenizer.EncodeTrimPrefix(str, 5, []string{imStart, imEnd})
		requireEqualTokenIDs(t, encoded.TokenIDs, []int{100264, 9906, 4435, 100265})
		if encoded.Text != str {
			t.Fatalf("expected %q got %q", str, encoded.Text)
		}

		encoded = tokenizer.EncodeTrimPrefix(str, 3, []string{imStart, imEnd})
		requireEqualTokenIDs(t, encoded.TokenIDs, []int{9906, 4435, 100265})
		if encoded.Text != encodedStr {
			t.Fatalf("expected %q got %q", encodedStr, encoded.Text)
		}
	})

	t.Run("encode trim prefix 2", func(t *testing.T) {
		str := "<|im_start|>HelloTemp World<|im_end|>"
		encodedStr := " World<|im_end|>"

		encoded := tokenizer.EncodeTrimPrefix(str, 5, []string{imStart, imEnd})
		requireEqualTokenIDs(t, encoded.TokenIDs, []int{100264, 9906, 12427, 4435, 100265})
		if encoded.Text != str {
			t.Fatalf("expected %q got %q", str, encoded.Text)
		}

		encoded = tokenizer.EncodeTrimPrefix(str, 6, []string{imStart, imEnd})
		requireEqualTokenIDs(t, encoded.TokenIDs, []int{100264, 9906, 12427, 4435, 100265})
		if encoded.Text != str {
			t.Fatalf("expected %q got %q", str, encoded.Text)
		}

		encoded = tokenizer.EncodeTrimPrefix(str, 3, []string{imStart, imEnd})
		requireEqualTokenIDs(t, encoded.TokenIDs, []int{4435, 100265})
		if encoded.Text != encodedStr {
			t.Fatalf("expected %q got %q", encodedStr, encoded.Text)
		}
	})

	t.Run("encode trim prefix 3", func(t *testing.T) {
		str := strings.Repeat("t", 4000)
		encoded := tokenizer.Encode(str, nil)
		trimmed := tokenizer.EncodeTrimPrefix(str, 5, []string{})
		if len(trimmed.TokenIDs) != 5 {
			t.Fatalf("expected 5 token ids, got %d", len(trimmed.TokenIDs))
		}
		requireEqualTokenIDs(t, trimmed.TokenIDs, encoded[len(encoded)-5:])
	})

	t.Run("tokenize source code gpt-3.5", func(t *testing.T) {
		source := string(readTestData(t, "lib.rs.txt"))
		want := readJSONTokenIDs(t, "tokens_gpt_3.5_turbo.json")
		encoded := tokenizer.Encode(source, []string{imStart, imEnd})
		if len(encoded) != 5584 {
			t.Fatalf("expected token length 5584, got %d", len(encoded))
		}
		requireEqualTokenIDs(t, encoded, want)
		if decoded := tokenizer.Decode(encoded); decoded != source {
			t.Fatalf("decoded output mismatch")
		}
	})
}

func TestTokenizerGoldenByEncodingAndModel(t *testing.T) {
	specialTokens := map[string]int{
		imStart: 100264,
		imEnd:   100265,
	}
	source := string(readTestData(t, "lib.rs.txt"))

	t.Run("p50k_base", func(t *testing.T) {
		tokenizer := newTokenizerForEncoding(t, "p50k_base", specialTokens)
		want := readJSONTokenIDs(t, "tokens_p50k_base.json")
		encoded := tokenizer.Encode(source, []string{imStart, imEnd})
		if len(encoded) != 7230 {
			t.Fatalf("expected token length 7230, got %d", len(encoded))
		}
		requireEqualTokenIDs(t, encoded, want)
		if decoded := tokenizer.Decode(encoded); decoded != source {
			t.Fatalf("decoded output mismatch")
		}
	})

	t.Run("p50k_edit", func(t *testing.T) {
		tokenizer := newTokenizerForEncoding(t, "p50k_edit", specialTokens)
		want := readJSONTokenIDs(t, "tokens_p50k_edit.json")
		encoded := tokenizer.Encode(source, []string{imStart, imEnd})
		if len(encoded) != 7230 {
			t.Fatalf("expected token length 7230, got %d", len(encoded))
		}
		requireEqualTokenIDs(t, encoded, want)
		if decoded := tokenizer.Decode(encoded); decoded != source {
			t.Fatalf("decoded output mismatch")
		}
	})

	t.Run("r50k_base", func(t *testing.T) {
		tokenizer := newTokenizerForEncoding(t, "r50k_base", specialTokens)
		want := readJSONTokenIDs(t, "tokens_r50k_base.json")
		encoded := tokenizer.Encode(source, []string{imStart, imEnd})
		if len(encoded) != 11378 {
			t.Fatalf("expected token length 11378, got %d", len(encoded))
		}
		requireEqualTokenIDs(t, encoded, want)
		if decoded := tokenizer.Decode(encoded); decoded != source {
			t.Fatalf("decoded output mismatch")
		}
	})

	t.Run("gpt2", func(t *testing.T) {
		tokenizer := newTokenizerForModel(t, "gpt2", specialTokens)
		want := readJSONTokenIDs(t, "tokens_gpt2.json")
		encoded := tokenizer.Encode(source, []string{imStart, imEnd})
		if len(encoded) != 11378 {
			t.Fatalf("expected token length 11378, got %d", len(encoded))
		}
		requireEqualTokenIDs(t, encoded, want)
		if decoded := tokenizer.Decode(encoded); decoded != source {
			t.Fatalf("decoded output mismatch")
		}
	})

	t.Run("gpt-4o", func(t *testing.T) {
		gpt4oSpecialTokens := map[string]int{
			"<|endoftext|>":   199999,
			"<|endofprompt|>": 200018,
		}
		tokenizer := newTokenizerForModel(t, "gpt-4o", gpt4oSpecialTokens)
		want := readJSONTokenIDs(t, "tokens_gpt_4o.json")
		encoded := tokenizer.Encode(source, []string{"<|endoftext|>", "<|endofprompt|>"})
		if len(encoded) != 5609 {
			t.Fatalf("expected token length 5609, got %d", len(encoded))
		}
		requireEqualTokenIDs(t, encoded, want)
		if decoded := tokenizer.Decode(encoded); decoded != source {
			t.Fatalf("decoded output mismatch")
		}
	})
}

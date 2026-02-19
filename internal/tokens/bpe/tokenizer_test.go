package tokenizer

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	imStart = "<|im_start|>"
	imEnd   = "<|im_end|>"
)

func requireEqualTokenIDs(t *testing.T, got, want []int) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("token mismatch\nwant: %v\ngot:  %v", want, got)
	}
}

func TestTikTokenizerCore(t *testing.T) {
	specialTokens := map[string]int{
		imStart: 100264,
		imEnd:   100265,
	}
	tokenizer, err := NewTokenizerForModel("gpt-5", specialTokens)
	require.NoError(t, err)

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

	t.Run("encode trim suffix 3", func(t *testing.T) {
		str := strings.Repeat("t", 4000)
		encoded := tokenizer.Encode(str, nil)
		trimmed := tokenizer.EncodeTrimSuffix(str, 5, []string{})
		if len(trimmed.TokenIDs) != 5 {
			t.Fatalf("expected 5 token ids, got %d", len(trimmed.TokenIDs))
		}
		requireEqualTokenIDs(t, trimmed.TokenIDs, encoded[:5])
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
}

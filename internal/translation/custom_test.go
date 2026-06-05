package translation

import "testing"

func TestReplacePlaceholders_QuotedJSONPlaceholders(t *testing.T) {
	translator := NewCustomTranslator(&CustomTranslatorConfig{})
	template := `{"text":"{{text}}","source_lang":"auto","target_lang":"{{target_lang}}"}`

	got := translator.replacePlaceholders(template, "hello \"world\"\nnext", "ZH")
	want := `{"text":"hello \"world\"\nnext","source_lang":"auto","target_lang":"ZH"}`

	if got != want {
		t.Fatalf("replacePlaceholders() = %q, want %q", got, want)
	}
}

func TestReplacePlaceholders_UnquotedJSONPlaceholders(t *testing.T) {
	translator := NewCustomTranslator(&CustomTranslatorConfig{})
	template := `{"text":{{text}},"target_lang":{{target_lang}},"source_lang":{{source_lang}}}`

	got := translator.replacePlaceholders(template, "hello", "ZH")
	want := `{"text":"hello","target_lang":"ZH","source_lang":"auto"}`

	if got != want {
		t.Fatalf("replacePlaceholders() = %q, want %q", got, want)
	}
}

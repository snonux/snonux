package processor

import (
	"strings"
	"testing"
)

func TestStripURLTrailing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "plain", in: "https://example.com", want: "https://example.com"},
		{name: "trailing period", in: "https://example.com.", want: "https://example.com"},
		{name: "multiple punctuation", in: "https://a.b/c).", want: "https://a.b/c"},
		{name: "empty", in: "", want: ""},
		{name: "only punctuation", in: "...", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := stripURLTrailing(tt.in)
			if got != tt.want {
				t.Fatalf("stripURLTrailing(%q) = %q; want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestAutolinkLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "no url escapes",
			in:   `hello <world>`,
			want: `hello &lt;world&gt;`,
		},
		{
			name: "single url",
			in:   "see https://foo.test ok",
			want: `see <a href="https://foo.test" target="_blank" rel="noopener noreferrer">https://foo.test</a> ok`,
		},
		{
			name: "url with trailing period in prose",
			in:   "Visit https://foo.test.",
			want: `Visit <a href="https://foo.test" target="_blank" rel="noopener noreferrer">https://foo.test</a>.`,
		},
		{
			name: "two urls",
			in:   "a http://a.com b https://b.org c",
			want: `a <a href="http://a.com" target="_blank" rel="noopener noreferrer">http://a.com</a> b <a href="https://b.org" target="_blank" rel="noopener noreferrer">https://b.org</a> c`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := autolinkLine(tt.in)
			if got != tt.want {
				t.Fatalf("autolinkLine(%q) = %q; want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestFormatParagraph(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "single line",
			in:   "hello",
			want: "hello",
		},
		{
			name: "line break",
			in:   "line one\nline two",
			want: "line one<br>\nline two",
		},
		{
			name: "skips blank lines inside para",
			in:   "a\n\nb",
			want: "a<br>\nb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := formatParagraph(tt.in)
			if got != tt.want {
				t.Fatalf("formatParagraph(%q) = %q; want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestFormatParagraph_autolinkMultiline(t *testing.T) {
	t.Parallel()
	got := formatParagraph("u https://x.y\nv")
	if !strings.Contains(got, `<a href="https://x.y"`) {
		t.Fatalf("expected autolink in multiline paragraph, got %q", got)
	}
	if !strings.Contains(got, "<br>") {
		t.Fatalf("expected br between lines, got %q", got)
	}
}

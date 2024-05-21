package golden

import (
	"testing"
)

func TestRegexReplaceSemverLike(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: `{
{"version": "v1.2.3"},
{"version": "v0.21.1-dev.0.0.20230223051417-a164292d2464"},
			}`,
			want: `{
{"version": "<<PRESENCE>>"},
{"version": "<<PRESENCE>>"},
			}`,
		},
		{
			input: `{{"version": "v1.2.3"}}`,
			want:  `{{"version": "<<PRESENCE>>"}}`,
		},
		{
			input: `{
{"non-go-version": "1.2.3"},
{"non-go-version": "v1.2"},
}`,
			want: `{
{"non-go-version": "1.2.3"},
{"non-go-version": "v1.2"},
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := regexReplaceGoSemverLike(tt.input,
				"<<PRESENCE>>"); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexReplaceGUID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: `{{"guid": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"}}`,
			want:  `{{"guid": "<<PRESENCE>>"}}`,
		},
		{
			input: `{{"non-guid": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a1h"}}`,
			want:  `{{"non-guid": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a1h"}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := regexReplaceGUID(tt.input, "<<PRESENCE>>"); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexReplaceElapsed(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: `{
{"elapsed": "552h59m59.999999999s"},
{"elapsed": "999.999999ms"},
{"elapsed": "999.999µs"},
{"elapsed": "999ns"},
{"elapsed": "552h59m59s"},
}`,
			want: `{
{"elapsed": "<<PRESENCE>>"},
{"elapsed": "<<PRESENCE>>"},
{"elapsed": "<<PRESENCE>>"},
{"elapsed": "<<PRESENCE>>"},
{"elapsed": "<<PRESENCE>>"},
}`,
		},
		{
			input: `{{"elapsed": "552h59m59.999999999s"}}`,
			want:  `{{"elapsed": "<<PRESENCE>>"}}`,
		},
		{
			input: `{"unmatched": "552h59m59.999999999s"}`,
			want:  `{"unmatched": "552h59m59.999999999s"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := regexReplaceElapsed(tt.input, "<<PRESENCE>>"); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexReplaceElapsedSeconds(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: `{
{"elapsed_seconds": 12.1},
{"elapsed_seconds": 12.123456789},
{"elapsed_seconds": 12},
{"elapsed_seconds": 12.0},
}`,
			want: `{
{"elapsed_seconds": 0.123},
{"elapsed_seconds": 0.123},
{"elapsed_seconds": 0.123},
{"elapsed_seconds": 0.123},
}`,
		},
		{
			input: `{{"elapsed_seconds": 12.1}}`,
			want:  `{{"elapsed_seconds": 0.123}}`,
		},
		{
			input: `{"unmatched": 12.1}`,
			want:  `{"unmatched": 12.1}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := regexReplaceElapsedSeconds(tt.input,
				0.123); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexReplaceStart(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: `{
{"start": "2023-01-12T21:22:57.581596Z"},
{"start": "2023-01-12T21:22:57.000001Z"},
{"start": "2023-01-14T03:49:30.018319141+01:00"},
{"start": "2023-01-12T21:22:57.000000003+00:00"},
{"start": "2023-01-12T21:22:57.000000003+01:00"},
{"start": "2023-01-12T21:22:57.000000003-01:00"},
{"start": "2023-01-12T21:22:57.000000003+01:30"},
}`,
			want: `{
{"start": "<<PRESENCE>>"},
{"start": "<<PRESENCE>>"},
{"start": "<<PRESENCE>>"},
{"start": "<<PRESENCE>>"},
{"start": "<<PRESENCE>>"},
{"start": "<<PRESENCE>>"},
{"start": "<<PRESENCE>>"},
}`,
		},
		{
			input: `{{"start": "2023-01-12T21:22:57.581596Z"}}`,
			want:  `{{"start": "<<PRESENCE>>"}}`,
		},
		{
			input: `{"unmatched": "2023-01-12T21:22:57.581596Z"}`,
			want:  `{"unmatched": "2023-01-12T21:22:57.581596Z"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := regexReplaceStart(tt.input, "<<PRESENCE>>"); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexReplaceCustom(t *testing.T) {
	tests := []struct {
		input       string
		regex       string
		placeHolder string
		want        string
	}{
		{
			input: `{
{"custom": "foo"},
{"custom": "bar"},
{"custom": "baz"},
}`,
			regex:       `foo|bar|baz`,
			placeHolder: "<<PRESENCE>>",
			want: `{
{"custom": "<<PRESENCE>>"},
{"custom": "<<PRESENCE>>"},
{"custom": "<<PRESENCE>>"},
}`,
		},
		{
			input:       `{{"run_id": "pizza-delivery-sl-Ql_TVR"}}`,
			regex:       `pizza-delivery-[0-9,a-z,A-Z,-_]{9}`,
			placeHolder: "<<PRESENCE>>",
			want:        `{{"run_id": "<<PRESENCE>>"}}`,
		},
		{
			input: ` v0.21.0            darwin_amd64         
                    darwin_arm64         
                    linux_amd64          
                    linux_arm64          
                                         
 v0.21.1-0.202302…  darwin_amd64         
                    darwin_arm64         
                    linux_amd64          
                    linux_arm64          
`,
			regex: `\s+\n\s+v\d+\.\d+\.\d+\-[0-9,a-f,\-,\.]+…` +
				`\s+darwin_amd64\s+darwin_arm64\s+linux_amd64\s+linux_arm64\s+`,
			placeHolder: "",
			want: ` v0.21.0            darwin_amd64         
                    darwin_arm64         
                    linux_amd64          
                    linux_arm64`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := regexReplaceCustom(
				tt.input,
				tt.placeHolder,
				tt.regex); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

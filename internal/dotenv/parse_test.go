package dotenv

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewParser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		ctx      context.Context
		opts     []ParserOption
		want     Parser
		checkErr func(t *testing.T, err error)
	}{
		{
			name: "success.default",
			ctx:  context.Background(),
			want: &parser{
				LineSeparator: "\n",
			},
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("NewParser: %v", err)
				}
			},
		},
		{
			name: "success.with_line_separator",
			ctx:  context.Background(),
			opts: []ParserOption{
				ParserOptionWithLineSeparator("\r\n"),
			},
			want: &parser{
				LineSeparator: "\r\n",
			},
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("NewParser: %v", err)
				}
			},
		},
		{
			name: "failure.context_canceled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			want: nil,
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("NewParser: %v, want: %v", err, context.Canceled)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewParser(tt.ctx, tt.opts...)
			if tt.checkErr != nil {
				tt.checkErr(t, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("NewParser() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		ctx      context.Context
		data     string
		want     *Dotenv
		checkErr func(t *testing.T, err error)
	}{
		{
			name: "success.quoted_value_with_comment",
			ctx:  context.Background(),
			data: `# comment line 1` + "\n" +
				`# comment line 2` + "\n" +
				`AAA='aaa'` + "\n" +
				`BBB="bbb"` + "\n" +
				`SKIP_IF_LINE_NOT_HAS_EQUAL` + "\n" +
				`UNTERMINATED="unterminated`, // TODO: discuss if this should be an error
			want: &Dotenv{
				Env: []Env{
					{Comment: "comment line 1\ncomment line 2\n", Key: "AAA", Value: "aaa"},
					{Comment: "", Key: "BBB", Value: "bbb"},
					{Comment: "", Key: "UNTERMINATED", Value: "\"unterminated"},
				},
			},
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("Parse: %v", err)
				}
			},
		},
		{
			name: "failure.context_canceled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			data: `AAA='aaa'`,
			want: nil,
			checkErr: func(t *testing.T, err error) {
				t.Helper()
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("Parse: %v, want: %v", err, context.Canceled)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parser, err := NewParser(t.Context())
			if err != nil {
				t.Fatalf("NewParser: %v", err)
			}
			got, err := parser.Parse(tt.ctx, tt.data)
			if tt.checkErr != nil {
				tt.checkErr(t, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Parse() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

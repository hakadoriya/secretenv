package dotenv

import (
	"context"
	"strings"
)

// Dotenv is a result of parsing .env file.
type Dotenv struct {
	Env []Env
}

// Env is a key-value pair (with comment) in .env file.
type Env struct {
	Comment string
	Key     string
	Value   string
}

// Parser is a parser for .env file.
type Parser interface {
	Parse(ctx context.Context, data string) (*Dotenv, error)
}

type parser struct {
	LineSeparator string
}

var _ Parser = (*parser)(nil)

// NewParser creates a new parser for .env file.
//
// The parser will use the default line separator "\n".
//
// If the line separator is not set, the parser will use the default line separator.
func NewParser(ctx context.Context, opts ...ParserOption) (Parser, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	p := &parser{
		LineSeparator: "\n",
	}
	for _, opt := range opts {
		opt.apply(p)
	}

	return p, nil
}

func (p *parser) Parse(ctx context.Context, data string) (*Dotenv, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// split by new line
	lines := strings.Split(data, p.LineSeparator)

	envs := make([]Env, 0, len(lines))
	var env Env

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)

		// skip comment line
		if comment, found := strings.CutPrefix(line, "#"); found {
			env.Comment += strings.TrimLeft(comment, " ") + "\n"
			continue
		}

		// trim `export ` prefix
		line = strings.TrimPrefix(line, "export ")
		line = strings.TrimSpace(line)

		// split by `=`
		const partsLen = 2
		parts := strings.SplitN(line, "=", partsLen)
		if len(parts) < partsLen {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if q := `'`; strings.HasPrefix(value, q) && strings.HasSuffix(value, q) {
			// if value is single quoted, trim the single quotes
			value = strings.Trim(value, q)
			value = strings.ReplaceAll(value, `\\`, `\`) // unescape backslash
			value = strings.ReplaceAll(value, `\`+q, q)  // unescape quote
		} else if q := `"`; strings.HasPrefix(value, q) && strings.HasSuffix(value, q) {
			// if value is double quoted, trim the double quotes
			value = strings.Trim(value, q)
			value = strings.ReplaceAll(value, `\\`, `\`) // unescape backslash
			value = strings.ReplaceAll(value, `\`+q, q)  // unescape quote
		}

		// set key and value
		env.Key = key
		env.Value = value

		// append to envs
		copied := env
		envs = append(envs, copied)

		// reset env
		//
		//nolint:exhaustruct
		env = Env{}
	}

	return &Dotenv{
		Env: envs,
	}, nil
}

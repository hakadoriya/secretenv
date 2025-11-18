package dotenv

// ParserOption is an option for the .env parser.
type ParserOption interface {
	apply(p *parser)
}

type parserOptionFunc func(p *parser)

func (f parserOptionFunc) apply(p *parser) { f(p) }

// ParserOptionWithLineSeparator sets the line separator for the parser.
//
// The default line separator is "\n".
//
// If the line separator is not set, the parser will use the default line separator.
func ParserOptionWithLineSeparator(lineSeparator string) ParserOption {
	return parserOptionFunc(func(p *parser) { p.LineSeparator = lineSeparator })
}

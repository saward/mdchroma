// Package bfchroma provides an easy and extensible blackfriday renderer that
// uses the chroma syntax highlighter to render code blocks.
package bfchroma

import (
	"io"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	md "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
)

// Option defines the functional option type
type Option func(r *Renderer)

// Style is a function option allowing to set the style used by chroma
// Default : "monokai"
func Style(s string) Option {
	return func(r *Renderer) {
		r.Style = styles.Get(s)
	}
}

// ChromaStyle is an option to directly set the style of the renderer using a
// chroma style instead of a string
func ChromaStyle(s *chroma.Style) Option {
	return func(r *Renderer) {
		r.Style = s
	}
}

// WithoutAutodetect disables chroma's language detection when no codeblock
// extra information is given. It will fallback to a sane default instead of
// trying to detect the language.
func WithoutAutodetect() Option {
	return func(r *Renderer) {
		r.Autodetect = false
	}
}

// EmbedCSS will embed CSS needed for html.WithClasses() in beginning of the document
func EmbedCSS() Option {
	return func(r *Renderer) {
		r.embedCSS = true
	}
}

// ChromaOptions allows to pass Chroma html.Option such as Standalone()
// WithClasses(), ClassPrefix(prefix)...
func ChromaOptions(options ...html.Option) Option {
	return func(r *Renderer) {
		r.ChromaOptions = options
	}
}

// Extend allows to specify the blackfriday renderer which is extended
func Extend(br md.Renderer) Option {
	return func(r *Renderer) {
		r.Base = br
	}
}

// NewRenderer will return a new bfchroma renderer with sane defaults
func NewRenderer(options ...Option) *Renderer {
	r := &Renderer{
		Base: mdhtml.NewRenderer(mdhtml.RendererOptions{
			Flags: mdhtml.CommonFlags,
		}),
		Style:      styles.Monokai,
		Autodetect: true,
	}
	for _, option := range options {
		option(r)
	}
	r.Formatter = html.New(r.ChromaOptions...)
	return r
}

// RenderWithChroma will render the given text to the w io.Writer
func (r *Renderer) RenderWithChroma(w io.Writer, text []byte, data ast.CodeBlock) error {
	var lexer chroma.Lexer

	// Determining the lexer to use
	if len(data.Info) > 0 {
		lexer = lexers.Get(string(data.Info))
	} else if r.Autodetect {
		lexer = lexers.Analyse(string(text))
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Tokenize the code
	iterator, err := lexer.Tokenise(nil, string(text))
	if err != nil {
		return err
	}
	return r.Formatter.Format(w, r.Style, iterator)
}

// Renderer is a custom Blackfriday renderer that uses the capabilities of
// chroma to highlight code with triple backtick notation
type Renderer struct {
	Base          md.Renderer
	Autodetect    bool
	ChromaOptions []html.Option
	Style         *chroma.Style
	Formatter     *html.Formatter
	embedCSS      bool
}

// RenderNode satisfies the Renderer interface
func (r *Renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	if _, ok := node.(*ast.Document); ok {
		if entering && r.embedCSS {
			w.Write([]byte("<style>"))
			r.Formatter.WriteCSS(w, r.Style)
			w.Write([]byte("</style>"))
		}
		return r.Base.RenderNode(w, node, entering)
	}
	if nodeData, ok := node.(*ast.CodeBlock); ok {
		if err := r.RenderWithChroma(w, nodeData.Literal, *nodeData); err != nil {
			return r.Base.RenderNode(w, node, entering)
		}
		return ast.SkipChildren
	}
	return r.Base.RenderNode(w, node, entering)
}

// RenderHeader satisfies the Renderer interface
func (r *Renderer) RenderHeader(w io.Writer, ast ast.Node) {
	r.Base.RenderHeader(w, ast)
}

// RenderFooter satisfies the Renderer interface
func (r *Renderer) RenderFooter(w io.Writer, ast ast.Node) {
	r.Base.RenderFooter(w, ast)
}

// ChromaCSS returns CSS used with chroma's html.WithClasses() option
func (r *Renderer) ChromaCSS(w io.Writer) error {
	return r.Formatter.WriteCSS(w, r.Style)
}


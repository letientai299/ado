package styles

import (
	"github.com/charmbracelet/glamour/ansi"
)

// Theme configuration for the application.
//
// Colors can be specified in several formats (consult lipgloss for examples):
//   - Hex: "#ffffff" or "#fff"
//
// - ANSI 16: "red", "green", "yellow", "blue", "magenta", "cyan", "white", "black" (and "bright"
// variants)
//   - ANSI 256: "21" (0-255)
//
// For shared themes, use the `include!` directive to load from external files:
//
//	theme:
//	  include!: "~/.config/ado/themes/tokyo-night.yaml"
type Theme struct {
	Name      string `json:"name"       yaml:"name"`
	TrueColor bool   `json:"true_color" yaml:"true_color"`
	Tokens    Tokens `json:"tokens"     yaml:"tokens"`
}

// Tokens contain the color code that the application will use to render outputs
type Tokens struct {
	Warn    string `json:"warning" yaml:"warning"`
	Error   string `json:"error"   yaml:"error"`
	Success string `json:"success" yaml:"success"`

	Markdown MarkdownTokens `json:"markdown" yaml:"markdown"`
	Chroma   ChromaTokens   `json:"chroma"   yaml:"chroma"`
}

type MarkdownTokens struct {
	Text           string `json:"text"            yaml:"text"`
	Heading        string `json:"heading"         yaml:"heading"`
	H1             string `json:"h1"              yaml:"h1"`
	H1Background   string `json:"h1_background"   yaml:"h1_background"`
	HorizontalRule string `json:"horizontal_rule" yaml:"horizontal_rule"`
	Link           string `json:"link"            yaml:"link"`
	LinkText       string `json:"link_text"       yaml:"link_text"`
	Image          string `json:"image"           yaml:"image"`
	ImageText      string `json:"image_text"      yaml:"image_text"`
	Code           string `json:"code"            yaml:"code"`
	CodeBackground string `json:"code_background" yaml:"code_background"`
	CodeBlock      string `json:"code_block"      yaml:"code_block"`
}

type ChromaTokens struct {
	Text                string `json:"text"                  yaml:"text"`
	Error               string `json:"error"                 yaml:"error"`
	ErrorBackground     string `json:"error_background"      yaml:"error_background"`
	Comment             string `json:"comment"               yaml:"comment"`
	CommentPreproc      string `json:"comment_preproc"       yaml:"comment_preproc"`
	Keyword             string `json:"keyword"               yaml:"keyword"`
	KeywordReserved     string `json:"keyword_reserved"      yaml:"keyword_reserved"`
	KeywordNamespace    string `json:"keyword_namespace"     yaml:"keyword_namespace"`
	KeywordType         string `json:"keyword_type"          yaml:"keyword_type"`
	Operator            string `json:"operator"              yaml:"operator"`
	Punctuation         string `json:"punctuation"           yaml:"punctuation"`
	Name                string `json:"name"                  yaml:"name"`
	NameBuiltin         string `json:"name_builtin"          yaml:"name_builtin"`
	NameTag             string `json:"name_tag"              yaml:"name_tag"`
	NameAttribute       string `json:"name_attribute"        yaml:"name_attribute"`
	NameClass           string `json:"name_class"            yaml:"name_class"`
	NameConstant        string `json:"name_constant"         yaml:"name_constant"`
	NameDecorator       string `json:"name_decorator"        yaml:"name_decorator"`
	Function            string `json:"function"              yaml:"function"`
	String              string `json:"string"                yaml:"string"`
	LiteralStringEscape string `json:"literal_string_escape" yaml:"literal_string_escape"`
	Number              string `json:"number"                yaml:"number"`
	GenericDeleted      string `json:"generic_deleted"       yaml:"generic_deleted"`
	GenericInserted     string `json:"generic_inserted"      yaml:"generic_inserted"`
	GenericSubheading   string `json:"generic_subheading"    yaml:"generic_subheading"`
	Background          string `json:"background"            yaml:"background"`
}

// glamourStyle creates a StyleConfig for glamour.TermRenderer.
func (t Theme) glamourStyle() ansi.StyleConfig {
	return ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "\n",
				BlockSuffix: "\n",
				Color:       ptr(t.Tokens.Markdown.Text),
			},
			Margin: ptr(uint(defaultMargin)),
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{},
			Indent:         ptr(uint(1)),
			IndentToken:    ptr("│ "),
		},
		List: ansi.StyleList{
			LevelIndent: defaultListIndent,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       ptr(t.Tokens.Markdown.Heading),
				Bold:        ptr(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color:           ptr(t.Tokens.Markdown.H1),
				BackgroundColor: ptr(t.Tokens.Markdown.H1Background),
				Bold:            ptr(true),
			},
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "## ",
			},
		},
		H3: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "### ",
			},
		},
		H4: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "#### ",
			},
		},
		H5: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "##### ",
			},
		},
		H6: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "###### ",
				Bold:   ptr(false),
			},
		},
		Strikethrough: ansi.StylePrimitive{
			CrossedOut: ptr(true),
		},
		Emph: ansi.StylePrimitive{
			Italic: ptr(true),
		},
		Strong: ansi.StylePrimitive{
			Bold: ptr(true),
		},
		HorizontalRule: ansi.StylePrimitive{
			Color:  ptr(t.Tokens.Markdown.HorizontalRule),
			Format: "\n--------\n",
		},
		Item: ansi.StylePrimitive{
			BlockPrefix: "• ",
		},
		Enumeration: ansi.StylePrimitive{
			BlockPrefix: ". ",
		},
		Task: ansi.StyleTask{
			StylePrimitive: ansi.StylePrimitive{},
			Ticked:         "[✓] ",
			Unticked:       "[ ] ",
		},
		Link: ansi.StylePrimitive{
			Color:     ptr(t.Tokens.Markdown.Link),
			Underline: ptr(true),
		},
		LinkText: ansi.StylePrimitive{
			Color: ptr(t.Tokens.Markdown.LinkText),
			Bold:  ptr(true),
		},
		Image: ansi.StylePrimitive{
			Color:     ptr(t.Tokens.Markdown.Image),
			Underline: ptr(true),
		},
		ImageText: ansi.StylePrimitive{
			Color:  ptr(t.Tokens.Markdown.ImageText),
			Format: "Image: {{.text}} →",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color:           ptr(t.Tokens.Markdown.Code),
				BackgroundColor: ptr(t.Tokens.Markdown.CodeBackground),
			},
		},

		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Markdown.CodeBlock),
				},
				Indent: ptr(uint(2)),
			},
			Chroma: &ansi.Chroma{
				Text: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.Text),
				},
				Error: ansi.StylePrimitive{
					Color:           ptr(t.Tokens.Chroma.Error),
					BackgroundColor: ptr(t.Tokens.Chroma.ErrorBackground),
				},
				Comment: ansi.StylePrimitive{
					Color:  ptr(t.Tokens.Chroma.Comment),
					Italic: ptr(true),
				},
				CommentPreproc: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.CommentPreproc),
				},
				Keyword: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.Keyword),
				},
				KeywordReserved: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.KeywordReserved),
				},
				KeywordNamespace: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.KeywordNamespace),
				},
				KeywordType: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.KeywordType),
				},
				Operator: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.Operator),
				},
				Punctuation: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.Punctuation),
				},
				Name: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.Name),
				},
				NameBuiltin: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.NameBuiltin),
				},
				NameTag: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.NameTag),
				},
				NameAttribute: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.NameAttribute),
				},
				NameClass: ansi.StylePrimitive{
					Color:     ptr(t.Tokens.Chroma.NameClass),
					Underline: ptr(true),
					Bold:      ptr(true),
				},
				NameConstant: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.NameConstant),
				},
				NameDecorator: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.NameDecorator),
				},
				NameFunction: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.Function),
				},
				LiteralNumber: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.Number),
				},
				LiteralString: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.String),
				},
				LiteralStringEscape: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.LiteralStringEscape),
				},
				GenericDeleted: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.GenericDeleted),
				},
				GenericEmph: ansi.StylePrimitive{
					Italic: ptr(true),
				},
				GenericInserted: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.GenericInserted),
				},
				GenericStrong: ansi.StylePrimitive{
					Bold: ptr(true),
				},
				GenericSubheading: ansi.StylePrimitive{
					Color: ptr(t.Tokens.Chroma.GenericSubheading),
				},
				Background: ansi.StylePrimitive{
					BackgroundColor: ptr(t.Tokens.Chroma.Background),
				},
			},
		},
		Table: ansi.StyleTable{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{},
			},
		},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n🠶 ",
		},
	}
}

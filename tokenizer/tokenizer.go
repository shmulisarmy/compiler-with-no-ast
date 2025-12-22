package tokenizer

const (
	TOKEN_EOF = iota
	TOKEN_IDENTIFIER
	TOKEN_NUMBER
	TOKEN_OPERATOR
	TOKEN_ASSIGN
	TOKEN_Punctuation
	TOKEN_STRING
)

type Token struct {
	Type  int
	Value string
}

type Tokenizer struct {
	input       string
	position    int
	currentChar byte
}

func NewTokenizer(input string) *Tokenizer {
	t := &Tokenizer{input: input, position: 0}
	if len(input) > 0 {
		t.currentChar = input[0]
	}
	return t
}

func (t *Tokenizer) advance() {
	t.position++
	if t.position >= len(t.input) {
		t.currentChar = 0
	} else {
		t.currentChar = t.input[t.position]
	}
}

func (t *Tokenizer) skipWhitespace() {
	for t.currentChar != 0 && (t.currentChar == ' ' || t.currentChar == '\t' || t.currentChar == '\n' || t.currentChar == '\r') {
		t.advance()
	}
}

func (t *Tokenizer) readNumber() string {
	result := ""
	for t.currentChar != 0 && ((t.currentChar >= '0' && t.currentChar <= '9') || t.currentChar == '.') {
		result += string(t.currentChar)
		t.advance()
	}
	return result
}

func (t *Tokenizer) readIdentifier() string {
	result := ""
	for t.currentChar != 0 && ((t.currentChar >= 'a' && t.currentChar <= 'z') || (t.currentChar >= 'A' && t.currentChar <= 'Z') || t.currentChar == '_' || (t.currentChar >= '0' && t.currentChar <= '9')) {
		result += string(t.currentChar)
		t.advance()
	}
	return result
}

func (t *Tokenizer) NextToken() Token {
	for t.currentChar != 0 {
		t.skipWhitespace()
		if t.position >= len(t.input) {
			return Token{Type: TOKEN_EOF, Value: ""}
		}

		if t.currentChar >= '0' && t.currentChar <= '9' {
			return Token{Type: TOKEN_NUMBER, Value: t.readNumber()}
		}

		if (t.currentChar >= 'a' && t.currentChar <= 'z') || (t.currentChar >= 'A' && t.currentChar <= 'Z') || t.currentChar == '_' {
			return Token{Type: TOKEN_IDENTIFIER, Value: t.readIdentifier()}
		}

		switch t.currentChar {
		case '+', '-', '*', '/', '=', '>', '<':
			c := t.currentChar
			t.advance()
			if t.currentChar == '=' {
				token := Token{Type: TOKEN_OPERATOR, Value: string([]rune{rune(c), rune('=')})}
				t.advance()
				return token
			}
			if c == '=' {
				return Token{Type: TOKEN_ASSIGN, Value: string(c)}
			}
			return Token{Type: TOKEN_OPERATOR, Value: string(c)}
		case '\'', '"':
			string_char := t.currentChar
			t.advance()
			start_index := t.position
			for t.currentChar != string_char {
				t.advance()
			}
			token := Token{Type: TOKEN_STRING, Value: t.input[start_index:t.position]}
			t.advance()
			return token

		case '(', ')', ';', ',', '.', '?', '{', '}':
			token := Token{Type: TOKEN_Punctuation, Value: string(t.currentChar)}
			t.advance()
			return token
		default:
			panic("Unknown character: " + string(t.currentChar))
		}
	}
	return Token{Type: TOKEN_EOF, Value: ""}
}

// Helper function to get token type name
func (t Token) String() string {
	switch t.Type {
	case TOKEN_EOF:
		return "EOF"
	case TOKEN_IDENTIFIER:
		return "IDENTIFIER(" + t.Value + ")"
	case TOKEN_NUMBER:
		return "NUMBER(" + t.Value + ")"
	case TOKEN_OPERATOR:
		return "OPERATOR(" + t.Value + ")"
	case TOKEN_ASSIGN:
		return "ASSIGN(" + t.Value + ")"
	case TOKEN_Punctuation:
		return "PUNCTUATION(" + t.Value + ")"
	default:
		return "UNKNOWN_TOKEN"
	}
}

// Tokenize function to get all tokens at once
func Tokenize(input string) []Token {
	t := NewTokenizer(input)
	var tokens []Token
	for {
		token := t.NextToken()
		tokens = append(tokens, token)
		if token.Type == TOKEN_EOF {
			break
		}
	}
	return tokens
}

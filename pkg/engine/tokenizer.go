package engine

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	util "github.com/codecrafters-io/interpreter-starter-go/pkg/util/engine_util"
)

type LexToken struct {
	TokenType string
	Lexeme    string
	Literal   string
}

func newToken(tokenType string, lexeme string, literal string) LexToken {
	return LexToken{tokenType, lexeme, literal}
}

func newTokenNoLit(tokenType string, lexeme string) LexToken {
	return LexToken{tokenType, lexeme, "null"}
}

type LexError struct {
	line    int
	message string
}

func (e LexError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.line, e.message)
}

const (
	TOKEN_EOF           = "EOF"
	TOKEN_LEFT_PAREN    = "LEFT_PAREN"
	TOKEN_RIGHT_PAREN   = "RIGHT_PAREN"
	TOKEN_LEFT_BRACE    = "LEFT_BRACE"
	TOKEN_RIGHT_BRACE   = "RIGHT_BRACE"
	TOKEN_LESS          = "LESS"
	TOKEN_GREATER       = "GREATER"
	TOKEN_BANG          = "BANG"
	TOKEN_EQUAL         = "EQUAL"
	TOKEN_SLASH         = "SLASH"
	TOKEN_STAR          = "STAR"
	TOKEN_DOT           = "DOT"
	TOKEN_COMMA         = "COMMA"
	TOKEN_PLUS          = "PLUS"
	TOKEN_MINUS         = "MINUS"
	TOKEN_SEMICOLON     = "SEMICOLON"
	TOKEN_BANG_EQUAL    = "BANG_EQUAL"
	TOKEN_EQUAL_EQUAL   = "EQUAL_EQUAL"
	TOKEN_LESS_EQUAL    = "LESS_EQUAL"
	TOKEN_GREATER_EQUAL = "GREATER_EQUAL"
	TOKEN_IDENTIFIER    = "IDENTIFIER"
	TOKEN_NUMBER        = "NUMBER"
	TOKEN_STRING        = "STRING"
	TOKEN_AND           = "AND"
	TOKEN_CLASS         = "CLASS"
	TOKEN_ELSE          = "ELSE"
	TOKEN_FALSE         = "FALSE"
	TOKEN_FOR           = "FOR"
	TOKEN_FUN           = "FUN"
	TOKEN_IF            = "IF"
	TOKEN_NIL           = "NIL"
	TOKEN_OR            = "OR"
	TOKEN_RETURN        = "RETURN"
	TOKEN_SUPER         = "SUPER"
	TOKEN_THIS          = "THIS"
	TOKEN_TRUE          = "TRUE"
	TOKEN_VAR           = "VAR"
	TOKEN_WHILE         = "WHILE"
	TOKEN_PRINT         = "PRINT"
	TOKEN_COMMENT       = "COMMENT"
)

var tokenPresets = map[string]LexToken{
	TOKEN_EOF:           newTokenNoLit(TOKEN_EOF, ""),
	TOKEN_LEFT_PAREN:    newTokenNoLit(TOKEN_LEFT_PAREN, "("),
	TOKEN_RIGHT_PAREN:   newTokenNoLit(TOKEN_RIGHT_PAREN, ")"),
	TOKEN_LEFT_BRACE:    newTokenNoLit(TOKEN_LEFT_BRACE, "{"),
	TOKEN_RIGHT_BRACE:   newTokenNoLit(TOKEN_RIGHT_BRACE, "}"),
	TOKEN_LESS:          newTokenNoLit(TOKEN_LESS, "<"),
	TOKEN_GREATER:       newTokenNoLit(TOKEN_GREATER, ">"),
	TOKEN_BANG:          newTokenNoLit(TOKEN_BANG, "!"),
	TOKEN_EQUAL:         newTokenNoLit(TOKEN_EQUAL, "="),
	TOKEN_SLASH:         newTokenNoLit(TOKEN_SLASH, "/"),
	TOKEN_STAR:          newTokenNoLit(TOKEN_STAR, "*"),
	TOKEN_DOT:           newTokenNoLit(TOKEN_DOT, "."),
	TOKEN_COMMA:         newTokenNoLit(TOKEN_COMMA, ","),
	TOKEN_PLUS:          newTokenNoLit(TOKEN_PLUS, "+"),
	TOKEN_MINUS:         newTokenNoLit(TOKEN_MINUS, "-"),
	TOKEN_SEMICOLON:     newTokenNoLit(TOKEN_SEMICOLON, ";"),
	TOKEN_BANG_EQUAL:    newTokenNoLit(TOKEN_BANG_EQUAL, "!="),
	TOKEN_EQUAL_EQUAL:   newTokenNoLit(TOKEN_EQUAL_EQUAL, "=="),
	TOKEN_LESS_EQUAL:    newTokenNoLit(TOKEN_LESS_EQUAL, "<="),
	TOKEN_GREATER_EQUAL: newTokenNoLit(TOKEN_GREATER_EQUAL, ">="),
	TOKEN_AND:           newTokenNoLit(TOKEN_AND, "and"),
	TOKEN_CLASS:         newTokenNoLit(TOKEN_CLASS, "class"),
	TOKEN_ELSE:          newTokenNoLit(TOKEN_ELSE, "else"),
	TOKEN_FALSE:         newTokenNoLit(TOKEN_FALSE, "false"),
	TOKEN_FOR:           newTokenNoLit(TOKEN_FOR, "for"),
	TOKEN_FUN:           newTokenNoLit(TOKEN_FUN, "fun"),
	TOKEN_IF:            newTokenNoLit(TOKEN_IF, "if"),
	TOKEN_NIL:           newTokenNoLit(TOKEN_NIL, "nil"),
	TOKEN_OR:            newTokenNoLit(TOKEN_OR, "or"),
	TOKEN_PRINT:         newTokenNoLit(TOKEN_PRINT, "print"),
	TOKEN_RETURN:        newTokenNoLit(TOKEN_RETURN, "return"),
	TOKEN_SUPER:         newTokenNoLit(TOKEN_SUPER, "super"),
	TOKEN_THIS:          newTokenNoLit(TOKEN_THIS, "this"),
	TOKEN_TRUE:          newTokenNoLit(TOKEN_TRUE, "true"),
	TOKEN_VAR:           newTokenNoLit(TOKEN_VAR, "var"),
	TOKEN_WHILE:         newTokenNoLit(TOKEN_WHILE, "while"),
}

func getPresetToken(tokenType string) LexToken {
	return tokenPresets[tokenType]
}

var runesTokenMap = map[string]string{
	"{":      TOKEN_LEFT_BRACE,
	"}":      TOKEN_RIGHT_BRACE,
	"(":      TOKEN_LEFT_PAREN,
	")":      TOKEN_RIGHT_PAREN,
	"*":      TOKEN_STAR,
	".":      TOKEN_DOT,
	",":      TOKEN_COMMA,
	"+":      TOKEN_PLUS,
	"-":      TOKEN_MINUS,
	";":      TOKEN_SEMICOLON,
	"/":      TOKEN_SLASH,
	"!":      TOKEN_BANG,
	"=":      TOKEN_EQUAL,
	"<":      TOKEN_LESS,
	">":      TOKEN_GREATER,
	"!=":     TOKEN_BANG_EQUAL,
	"==":     TOKEN_EQUAL_EQUAL,
	"<=":     TOKEN_LESS_EQUAL,
	">=":     TOKEN_GREATER_EQUAL,
	"//":     TOKEN_COMMENT,
	"and":    TOKEN_AND,
	"class":  TOKEN_CLASS,
	"else":   TOKEN_ELSE,
	"false":  TOKEN_FALSE,
	"for":    TOKEN_FOR,
	"fun":    TOKEN_FUN,
	"if":     TOKEN_IF,
	"nil":    TOKEN_NIL,
	"or":     TOKEN_OR,
	"print":  TOKEN_PRINT,
	"return": TOKEN_RETURN,
	"super":  TOKEN_SUPER,
	"this":   TOKEN_THIS,
	"true":   TOKEN_TRUE,
	"var":    TOKEN_VAR,
	"while":  TOKEN_WHILE,
}
var combineTokenRuneTriggers = []byte{
	'=',
	'/',
}

var i int
var line int

func Tokenize(fileContents []byte) ([]LexToken, []error) {
	errs := make([]error, 0)
	lexTokens := make([]LexToken, 0)

	line = 1
	for i = 0; i < len(fileContents); i++ {
		b := fileContents[i]

		if b == '\n' {
			line++
			continue
		}

		if util.IsWhitespace(b) {
			continue
		}

		if b == '"' {
			err := tokenizeString(fileContents, &lexTokens)
			if err != nil {
				errs = append(errs, err)
			}
			continue
		}

		if util.IsDigit(b) {
			err := tokenizeNumber(fileContents, &lexTokens)
			if err != nil {
				errs = append(errs, err)
			}
			continue
		}

		if _, ok := runesTokenMap[string(b)]; ok {
			tokenizeSimpleToken(fileContents, &lexTokens)
			continue
		}

		if util.IsAlpha(b) || b == '_' {
			tokenizeIdentifier(fileContents, &lexTokens)
			continue
		}

		message := fmt.Sprintf("Unexpected character: %s", string(b))
		errs = append(errs, LexError{line, message})
	}
	lexTokens = append(lexTokens, getPresetToken(TOKEN_EOF))
	return lexTokens, errs
}

func tokenizeString(fileContents []byte, tokens *[]LexToken) error {
	if fileContents[i] != '"' {
		message := "Tried to tokenize a string but the first character was not a double quote."
		return LexError{line, message}
	}

	stringBuffer := ""
	i++

	for ; i < len(fileContents); i++ {
		b := fileContents[i]
		if b == '"' {
			stringToken := newToken(
				TOKEN_STRING,
				"\""+stringBuffer+"\"",
				stringBuffer)

			*tokens = append(*tokens, stringToken)
			return nil
		} else if b == '\n' {
			i--
			message := "Unterminated string."
			return LexError{line, message}
		}
		stringBuffer += string(b)
	}
	message := "Unterminated string."
	return LexError{line, message}
}

func tokenizeNumber(fileContents []byte, tokens *[]LexToken) error {
	if !util.IsDigit(fileContents[i]) {
		message := "Tried to tokenize a number but the first character was not a digit."
		return LexError{line, message}
	}

	stringBuffer := string(fileContents[i])
	i++

	for ; i < len(fileContents); i++ {
		b := fileContents[i]
		if !util.IsDigit(b) && !util.IsFirstDot(b, stringBuffer) {
			i--
			break
		} else if util.IsFirstDot(b, stringBuffer) && i >= len(fileContents)-1 {
			i--
			break
		}
		stringBuffer += string(b)
	}

	if stringBuffer[len(stringBuffer)-1] == '.' {
		stringBuffer = stringBuffer[:len(stringBuffer)-1]
		i--
	}

	f, err := strconv.ParseFloat(stringBuffer, 64)
	if err != nil {
		message := fmt.Sprintf("Invalid number: %s", stringBuffer)
		return LexError{line, message}
	}

	formattedBuffer := strconv.FormatFloat(f, 'f', -1, 64)
	if !strings.Contains(formattedBuffer, ".") {
		formattedBuffer += ".0"
	}
	*tokens = append(*tokens, newToken(
		TOKEN_NUMBER,
		stringBuffer,
		formattedBuffer))

	return nil
}

func tokenizeSimpleToken(fileContents []byte, tokens *[]LexToken) {
	if !slices.Contains(combineTokenRuneTriggers, fileContents[i]) {
		tokenType := runesTokenMap[string(fileContents[i])]
		*tokens = append(*tokens, getPresetToken(tokenType))
		return
	}
	if i > 0 && i < len(fileContents) {
		previousByte := fileContents[i-1]
		if !util.IsWhitespace(previousByte) {
			lastToken := (*tokens)[len(*tokens)-1]
			combinedLexeme := lastToken.Lexeme + string(fileContents[i])
			if tokenType, ok := runesTokenMap[combinedLexeme]; ok {
				if tokenType == TOKEN_COMMENT {
					*tokens = (*tokens)[:len(*tokens)-1]
					for ; i < len(fileContents); i++ {
						if fileContents[i] == '\n' {
							i--
							return
						}
					}
				}
				*tokens = (*tokens)[:len(*tokens)-1]
				*tokens = append(*tokens, getPresetToken(tokenType))
				return
			}
		}
	}
	tokenType := runesTokenMap[string(fileContents[i])]
	*tokens = append(*tokens, getPresetToken(tokenType))
}

func tokenizeIdentifier(fileContents []byte, tokens *[]LexToken) {
	stringBuffer := string(fileContents[i])
	i++

	for ; i < len(fileContents); i++ {
		b := fileContents[i]
		if !util.IsAlpha(b) && !util.IsDigit(b) && b != '_' {
			i--
			break
		}
		stringBuffer += string(b)
	}

	if tokenType, ok := runesTokenMap[stringBuffer]; ok {
		*tokens = append(*tokens, getPresetToken(tokenType))
		return
	}
	*tokens = append(*tokens, newTokenNoLit(TOKEN_IDENTIFIER, stringBuffer))
}

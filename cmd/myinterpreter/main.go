package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	fileContents, cmd := readArgs()

	tokens, errors := tokenize(fileContents)

	errorsLen := len(errors)
	if errorsLen > 0 {
		for _, e := range errors {
			fmt.Fprintln(os.Stderr, e.Error())
		}
	}

	if cmd == "tokenize" {
		for _, t := range tokens {
			fmt.Printf("%s %s %s\n", t.tokenType, t.lexeme, t.literal)
		}
	}

	if cmd == "parse" {
		parse(tokens)
	}

	if errorsLen > 0 {
		os.Exit(65)
	}

}

func readArgs() ([]byte, string) {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" && command != "parse" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	return fileContents, command
}

type LexError struct {
	line    int
	message string
}

func (e LexError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.line, e.message)
}

type LexToken struct {
	tokenType string
	lexeme    string
	literal   string
}

func newToken(tokenType string, lexeme string, literal string) LexToken {
	return LexToken{tokenType, lexeme, literal}
}

func newTokenNoLit(tokenType string, lexeme string) LexToken {
	return LexToken{tokenType, lexeme, "null"}
}

var singleCharTokens = map[byte]LexToken{
	'{': newTokenNoLit("LEFT_BRACE", "{"),
	'}': newTokenNoLit("RIGHT_BRACE", "}"),
	'(': newTokenNoLit("LEFT_PAREN", "("),
	')': newTokenNoLit("RIGHT_PAREN", ")"),
	'*': newTokenNoLit("STAR", "*"),
	'.': newTokenNoLit("DOT", "."),
	',': newTokenNoLit("COMMA", ","),
	'+': newTokenNoLit("PLUS", "+"),
	'-': newTokenNoLit("MINUS", "-"),
	';': newTokenNoLit("SEMICOLON", ";"),
	'/': newTokenNoLit("SLASH", "/"),
	'!': newTokenNoLit("BANG", "!"),
	'=': newTokenNoLit("EQUAL", "="),
	'<': newTokenNoLit("LESS", "<"),
	'>': newTokenNoLit("GREATER", ">"),
}

var dualCharTokenTriggers = []string{
	"<",
	">",
	"!",
	"=",
	"/",
}

var dualCharTokens = map[string]LexToken{
	"!=": newTokenNoLit("BANG_EQUAL", "!="),
	"==": newTokenNoLit("EQUAL_EQUAL", "=="),
	"<=": newTokenNoLit("LESS_EQUAL", "<="),
	">=": newTokenNoLit("GREATER_EQUAL", ">="),
	"//": newTokenNoLit("COMMENT", "//"),
}

var reservedWords = map[string]LexToken{
	"and":    newTokenNoLit("AND", "and"),
	"class":  newTokenNoLit("CLASS", "class"),
	"else":   newTokenNoLit("ELSE", "else"),
	"false":  newTokenNoLit("FALSE", "false"),
	"for":    newTokenNoLit("FOR", "for"),
	"fun":    newTokenNoLit("FUN", "fun"),
	"if":     newTokenNoLit("IF", "if"),
	"nil":    newTokenNoLit("NIL", "nil"),
	"or":     newTokenNoLit("OR", "or"),
	"return": newTokenNoLit("RETURN", "return"),
	"super":  newTokenNoLit("SUPER", "super"),
	"this":   newTokenNoLit("THIS", "this"),
	"true":   newTokenNoLit("TRUE", "true"),
	"var":    newTokenNoLit("VAR", "var"),
	"while":  newTokenNoLit("WHILE", "while"),
	"eof":    newTokenNoLit("EOF", "eof"),
	"print":  newTokenNoLit("PRINT", "print"),
}

var ignoreChars = []byte{' ', '\t', '\r', '	'}

var content = []byte{}
var contentN = -1
var currentLine = 1

func nextByte() (byte, bool) {
	if contentN >= len(content)-1 {
		return 0, false
	}
	contentN += 1
	result := content[contentN]
	return result, true
}

func tickBack() {
	contentN -= 1
}

func parse(tokens []LexToken) {
	for _, t := range tokens {
		if t.tokenType == "TRUE" || t.tokenType == "FALSE" {
			fmt.Printf("%s\n", t.lexeme)
		}
	}
}

func tokenize(input []byte) ([]LexToken, []error) {
	content = input
	tokens := []LexToken{}
	errors := []error{}

	for b, ok := nextByte(); ok; b, ok = nextByte() {

		if slices.Contains(ignoreChars, b) {
			continue
		}

		if b == '\n' {
			currentLine += 1
			continue
		}

		if lexToken, ok := singleCharTokens[b]; ok {
			handleSingleCharToken(lexToken, &tokens)
			continue
		}

		if b == '"' {
			err := handleStringToken(&tokens)
			if err != nil {
				errors = append(errors, err)
			}
			continue
		}

		if b >= '0' && b <= '9' {
			err := handleNumberToken(&tokens)
			if err != nil {
				errors = append(errors, err)
			}
			continue
		}

		if b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b == '_' {
			handleIdentifierToken(&tokens)
			continue
		}

		// doesn't match any token
		message := fmt.Sprintf("Unexpected character: %s", string(b))
		errors = append(errors, LexError{currentLine, message})
	}
	tokens = append(tokens, newTokenNoLit("EOF", ""))
	return tokens, errors
}

func handleIdentifierToken(tokens *[]LexToken) {
	tickBack()
	stringBuffer := ""

	for b, ok := nextByte(); ok; b, ok = nextByte() {
		if b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9' || b == '_' {
			stringBuffer += string(b)
			continue
		} else {
			tickBack()
			break
		}
	}

	if lexToken, ok := reservedWords[stringBuffer]; ok {
		*tokens = append(*tokens, lexToken)
		return
	}
	*tokens = append(*tokens, newToken("IDENTIFIER", stringBuffer, "null"))
}

func handleNumberToken(tokens *[]LexToken) error {
	tickBack()
	stringBuffer := ""

	for b, ok := nextByte(); ok; b, ok = nextByte() {
		if b >= '0' && b <= '9' {
			stringBuffer += string(b)
			continue
		} else if b == '.' && !strings.Contains(stringBuffer, ".") {
			stringBuffer += string(b)
			continue
		} else {
			tickBack()
			break
		}
	}

	if stringBuffer[len(stringBuffer)-1] == '.' {
		stringBuffer = stringBuffer[:len(stringBuffer)-1]
		tickBack()
	}

	f, err := strconv.ParseFloat(stringBuffer, 64)
	if err != nil {
		message := fmt.Sprintf("Invalid number: %s", stringBuffer)
		return LexError{currentLine, message}
	}

	formattedBuffer := strconv.FormatFloat(f, 'f', -1, 64)
	if !strings.Contains(formattedBuffer, ".") {
		formattedBuffer += ".0"
	}
	*tokens = append(*tokens, newToken("NUMBER", stringBuffer, formattedBuffer))
	return nil
}

func handleStringToken(tokens *[]LexToken) error {
	stringBuffer := ""

	for b, ok := nextByte(); ok; b, ok = nextByte() {
		if b == '"' {
			*tokens = append(*tokens, newToken("STRING", "\""+stringBuffer+"\"", stringBuffer))
			return nil
		} else if b == '\n' {
			tickBack()
			message := "Unterminated string."
			return LexError{currentLine, message}
		}
		stringBuffer += string(b)
	}
	message := "Unterminated string."
	return LexError{currentLine, message}
}

func handleSingleCharToken(lexToken LexToken, tokens *[]LexToken) {
	if !slices.Contains(dualCharTokenTriggers, lexToken.lexeme) {
		*tokens = append(*tokens, lexToken)
		return
	}

	b, ok := nextByte()
	if !ok {
		*tokens = append(*tokens, lexToken)
		return
	}

	lexemeCombined := lexToken.lexeme + string(b)
	dualLexToken, ok := dualCharTokens[lexemeCombined]
	if !ok {
		*tokens = append(*tokens, lexToken)
		tickBack()
		return
	}

	if dualLexToken.tokenType != "COMMENT" {
		*tokens = append(*tokens, dualLexToken)
		return
	}

	for cb, ok := nextByte(); ok; cb, ok = nextByte() {
		if cb == '\n' {
			tickBack()
			return
		}
	}
}

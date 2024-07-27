package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type LexError struct {
	line    int
	message string
}

func (e LexError) String() string {
	return fmt.Sprintf("[line %v] Error: %s", e.line, e.message)
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

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	tokens, errors := tokenize(fileContents)

	errorsLen := len(errors)
	if errorsLen > 0 {
		for _, e := range errors {
			fmt.Fprintln(os.Stderr, e.String())
		}
	}

	for _, t := range tokens {
		fmt.Printf("%s %s %s\n", t.tokenType, t.lexeme, t.literal)
	}

	if errorsLen > 0 {
		os.Exit(65)
	}

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

var dualCharTokensTriggers = []string{
	"EQUAL",
	"SLASH",
}

var dualCharTokens = map[string]LexToken{
	"!=": newTokenNoLit("BANG_EQUAL", "!="),
	"==": newTokenNoLit("EQUAL_EQUAL", "=="),
	"<=": newTokenNoLit("LESS_EQUAL", "<="),
	">=": newTokenNoLit("GREATER_EQUAL", ">="),
	"//": newTokenNoLit("COMMENT", "//"),
}

var ignoreChars = []byte{' ', '\t', '\r', '	'}

func tokenize(content []byte) ([]LexToken, []LexError) {
	tokens := []LexToken{}
	errors := []LexError{}
	currentLine := 1
	breakContinuity := false
	commentActive := false
	stringBuffer := ""
	stringActive := false
	numberBuffer := ""
	numberActive := false
	if len(content) > 0 {
		for _, b := range content {

			// is a digit
			if b >= '0' && b <= '9' {
				if !numberActive {
					numberActive = true
				}
				numberBuffer += string(b)
				continue
			}

			if numberActive {
				if b == '.' {
					numberBuffer += string(b)
					continue
				} else if b < '0' || b > '9' {
					numberActive = false
					if !strings.Contains(numberBuffer, ".") {
						numberBuffer += ".0"
					}
					tokens = append(tokens, newToken("NUMBER", numberBuffer, numberBuffer))
					numberBuffer = ""
				}
			}

			if b == '"' {
				if stringActive {
					stringActive = !stringActive
					tokens = append(tokens, newToken("STRING", "\""+stringBuffer+"\"", stringBuffer))
					stringBuffer = ""
				} else {
					stringActive = true
				}
				continue
			}

			if b == '\n' {
				if stringActive {
					message := "Unterminated string."
					errors = append(errors, LexError{currentLine, message})
					stringActive = false
					stringBuffer = ""
				}
				commentActive = false
				breakContinuity = true
				currentLine += 1
				continue
			}

			if stringActive {
				stringBuffer += string(b)
				continue
			}

			if commentActive || slices.Contains(ignoreChars, b) {
				breakContinuity = true
				continue
			}

			if lexToken, ok := singleCharTokens[b]; ok {
				isTokenContinous := slices.Contains(dualCharTokensTriggers, lexToken.tokenType)
				if isTokenContinous && len(tokens) > 0 {
					prev := tokens[len(tokens)-1]
					lexemeCombined := prev.lexeme + lexToken.lexeme
					if dualLexToken, ok := dualCharTokens[lexemeCombined]; ok && !breakContinuity {
						if dualLexToken.tokenType == "COMMENT" {
							commentActive = true
							tokens = tokens[:len(tokens)-1]
							continue
						}
						tokens[len(tokens)-1] = dualLexToken
					} else {
						tokens = append(tokens, lexToken)
					}
				} else {
					tokens = append(tokens, lexToken)
				}
			} else {
				message := fmt.Sprintf("Unexpected character: %s", string(b))
				errors = append(errors, LexError{currentLine, message})
				breakContinuity = true
				continue
			}
			breakContinuity = false
		}
	}
	if numberActive {
		if !strings.Contains(numberBuffer, ".") {
			numberBuffer += ".0"
		}
		tokens = append(tokens, newToken("NUMBER", numberBuffer, numberBuffer))
	}
	if stringActive {
		message := "Unterminated string."
		errors = append(errors, LexError{currentLine, message})
	}
	tokens = append(tokens, newTokenNoLit("EOF", ""))

	return tokens, errors
}

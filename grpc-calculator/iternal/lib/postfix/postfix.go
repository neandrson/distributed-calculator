package postfix

import (
	"fmt"
	"strconv"
	"strings"

	shuntingYard "github.com/mgenware/go-shunting-yard"
)

func Postfix(expr string) ([]string, error) {
	infixTokens, err := shuntingYard.Scan(expr)
	if err != nil {
		return nil, err
	}

	postfixTokens, err := shuntingYard.Parse(infixTokens)
	if err != nil {
		return nil, err
	}

	s := make([]string, 0)
	for _, char := range postfixTokens {
		switch val := char.Value.(type) {
		case int:
			s = append(s, fmt.Sprint(val))
		case string:
			s = append(s, val)
		}
	}
	if len(s) >= 5 {
		for i := 0; i < len(s)-5; i++ {
			if isNumeric(s[i]) && isNumeric(s[i+1]) && strings.ContainsAny(s[i+2], "+-/*^") && isNumeric(s[i+3]) && strings.ContainsAny(s[i+4], "+-/*^") && isNumeric(s[i+5]) {
				char := s[i+4]
				s[i+4] = s[i+5]
				s[i+5] = char
			}
		}
	}
	return s, nil
}
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
func SaveSlice(in []string) string {
	s := ""

	for _, cha := range in {
		s = s + cha + ":"
	}

	return s
}
func ParsSlice(in string) []string {
	out := make([]string, 0)

	for _, char := range in {

		if string(char) != ":" {
			out = append(out, string(char))

		}

	}

	return out
}

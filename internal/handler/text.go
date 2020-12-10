package handler

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	regexpForWhiteSpace = regexp.MustCompile(`[\s　]+`)
	regexpForMention    = regexp.MustCompile(`(?P<mention><@U[0-9A-Z]+>)`)
)

func extractMentionIDsAndTokens(text string) ([]string, []string) {
	text = sanitizeNotPrintableChars(text)
	text = normalizeText(text)
	words := strings.Split(text, " ")
	ids := make([]string, 0, len(words))
	tokens := make([]string, 0, len(words))
	for _, word := range words {
		if strings.HasPrefix(word, "<@U") && strings.HasSuffix(word, ">") {
			ids = append(ids, string([]rune(word)[2:len(word)-1]))
		} else {
			tokens = append(tokens, strings.TrimSpace(word))
		}
	}

	return ids, tokens
}

func sanitizeNotPrintableChars(str string) string {
	runes := make([]rune, 0, utf8.RuneCountInString(str))
	for _, r := range str {
		if (0x20 <= r && r <= 0x7e) || (0xa1 <= r && r <= 0xdf) || (0xff < r) {
			runes = append(runes, r)
		} else {
			runes = append(runes, ' ')
		}
	}
	return string(runes)
}

func normalizeText(text string) string {
	text = regexpForMention.ReplaceAllString(text, " ${mention} ")
	text = regexpForWhiteSpace.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)
	return text
}

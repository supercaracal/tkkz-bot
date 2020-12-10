package handler

import (
	"testing"
)

func TestExtractMentionIDsAndTokens(t *testing.T) {
	cases := []struct {
		paramText      string
		expectedIds    []string
		expectedTokens []string
	}{
		{"foobar", []string{}, []string{"foobar"}},
		{"<@U0000AAAA> foobar", []string{"U0000AAAA"}, []string{"foobar"}},
		{"<@U0000AAAA> <@U0000BBBB> foobar", []string{"U0000AAAA", "U0000BBBB"}, []string{"foobar"}},
		{"　\t<@U0000AAAA>\r\n <@U0000BBBB> foobar 　", []string{"U0000AAAA", "U0000BBBB"}, []string{"foobar"}},
		{"　foo\t<@U0000AAAA>\r\nbar <@U0000BBBB> baz 　zap ", []string{"U0000AAAA", "U0000BBBB"}, []string{"foo", "bar", "baz", "zap"}},
		{"<abc> foobar", []string{}, []string{"<abc>", "foobar"}},
		{"foo<@U0000AAAA>bar", []string{"U0000AAAA"}, []string{"foo", "bar"}},
		{"foo<@U0000AAAA><@U0000BBBB>bar", []string{"U0000AAAA", "U0000BBBB"}, []string{"foo", "bar"}},
		{"foo<@U0000AAAA><@U0000BBBB><@U0000CCCC>bar", []string{"U0000AAAA", "U0000BBBB", "U0000CCCC"}, []string{"foo", "bar"}},
		{"foo<@U0000AAAA>bar<@U0000BBBB>baz", []string{"U0000AAAA", "U0000BBBB"}, []string{"foo", "bar", "baz"}},
		{"<@U0000AAAA> \xc2\xa0foobar", []string{"U0000AAAA"}, []string{"foobar"}},
		{string([]rune{'<', '@', 'U', '1', '>', ' ', 'f', 'o', 'o', 0xa0, 'b', 'a', 'r'}), []string{"U1"}, []string{"foo", "bar"}},
		{string([]rune{'<', '@', 'U', '1', '>', ' ', 'あ', 0xa0, 'い'}), []string{"U1"}, []string{"あ", "い"}},
	}

	for n, c := range cases {
		actualIds, actualTokens := extractMentionIDsAndTokens(c.paramText)
		if len(c.expectedIds) != len(actualIds) {
			t.Errorf("%d: expected=%d, actual=%d", n, len(c.expectedIds), len(actualIds))
			continue
		}
		for i, id := range actualIds {
			if c.expectedIds[i] != id {
				t.Errorf("%d: expected=%s, actual=%s", n, c.expectedIds[i], id)
				break
			}
		}
		if len(c.expectedTokens) != len(actualTokens) {
			t.Errorf("%d: expected=%d, actual=%d", n, len(c.expectedTokens), len(actualTokens))
			continue
		}
		for i, token := range actualTokens {
			if c.expectedTokens[i] != token {
				t.Errorf("%d: expected=%s, actual=%s", n, c.expectedTokens[i], token)
				break
			}
		}
	}
}

func TestSanitizeNotPrintableChars(t *testing.T) {
	cases := []struct {
		str  string
		want string
	}{
		{"pqr", "pqr"},
		{"ｱｲｳ", "ｱｲｳ"},
		{"あいう", "あいう"},
		{string([]rune{'p', 0x1f, 'q', 0x7f, 'r'}), "p q r"},
		{string([]rune{'p', 0xa0, 'q', 0xa0, 'r'}), "p q r"},
		{string([]rune{'ｱ', 0xa0, 'ｲ', 0xe0, 'ｳ'}), "ｱ ｲ ｳ"},
		{string([]rune{'ｱ', 0xa0, 'ｲ', 0xa0, 'ｳ'}), "ｱ ｲ ｳ"},
		{string([]rune{'あ', 0xa0, 'い', 0xa0, 'う'}), "あ い う"},
	}

	for n, c := range cases {
		if got := sanitizeNotPrintableChars(c.str); got != c.want {
			t.Errorf("%d: want=%s, got=%s", n, c.want, got)
		}
	}
}

func TestNormalizeText(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{"a\tb　c", "a b c"},
		{"foo<@U0000AAAA>bar", "foo <@U0000AAAA> bar"},
		{" a  bc ", "a bc"},
		{" あ  いう ", "あ いう"},
	}

	for n, c := range cases {
		if got := normalizeText(c.text); got != c.want {
			t.Errorf("%d: want=%s, got=%s", n, c.want, got)
		}
	}
}

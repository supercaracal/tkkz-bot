package main

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
	}

	for n, c := range cases {
		actualIds, actualTokens := extractMentionIDsAndTokens(c.paramText)
		if len(c.expectedIds) != len(actualIds) {
			t.Errorf("%d: expected=%d, actual=%d", n, len(c.expectedIds), len(actualIds))
		}
		for i, id := range actualIds {
			if c.expectedIds[i] != id {
				t.Errorf("%d: expected=%s, actual=%s", n, c.expectedIds[i], id)
				break
			}
		}
		if len(c.expectedTokens) != len(actualTokens) {
			t.Errorf("%d: expected=%d, actual=%d", n, len(c.expectedTokens), len(actualTokens))
		}
		for i, token := range actualTokens {
			if c.expectedTokens[i] != token {
				t.Errorf("%d: expected=%s, actual=%s", n, c.expectedTokens[i], token)
				break
			}
		}
	}
}

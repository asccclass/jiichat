package main

import (
	"strings"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func MarkdownToHTML(input string) (string, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),		
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithHardWraps(),html.WithXHTML(),html.WithUnsafe()),
	)
	var buf strings.Builder
	if err := md.Convert([]byte(input), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
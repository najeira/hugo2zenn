package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	tagReplacer = strings.NewReplacer(
		" ", "", "-", "", "_", "",
	)
)

func convertHugoItems(dir string, items []*HugoItem) {
	for _, item := range items {
		writeZennItem(dir, item)
	}
}

func writeZennItem(dir string, item *HugoItem) {
	body := convertHugoItem(item)

	date := item.Date.Format("2006-01-02")
	name := fmt.Sprintf("%s-%s.md", date, item.Slug)
	fn := filepath.Join(dir, name)
	f, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := io.Copy(f, strings.NewReader(body)); err != nil {
		panic(err)
	}
}

func convertHugoItem(item *HugoItem) string {
	var out strings.Builder
	out.WriteString("---\n")
	writeZennMetaString(&out, "title", item.Title)
	writeZennMetaString(&out, "slug", item.Slug)
	writeZennMetaString(&out, "emoji", "👻")
	writeZennMetaString(&out, "type", "tech")
	writeZennMetaStrings(&out, "topics", item.Tags)
	out.WriteString("published: true\n")

	out.WriteString("published_date: \"")
	out.WriteString(item.Date.Format("2006-01-02"))
	out.WriteString("\"\n")

	out.WriteString("---\n")

	if !strings.HasPrefix(item.Content, "\n") {
		out.WriteString("\n")
	}
	out.WriteString(item.Content)
	return out.String()
}

func writeZennMetaString(out *strings.Builder, name string, body string) {
	out.WriteString(name)
	out.WriteString(": ")
	out.WriteString("\"")
	out.WriteString(body)
	out.WriteString("\"\n")
}

func writeZennMetaStrings(out *strings.Builder, name string, tags []string) {
	out.WriteString(name)
	out.WriteString(": [")
	for i, tag := range tags {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString("\"")

		// 記号やスペースを使えないので変換しておく
		tag = tagReplacer.Replace(tag)
		out.WriteString(tag)
		out.WriteString("\"")
	}
	out.WriteString("]\n")
}

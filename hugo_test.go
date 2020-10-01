package main

import (
	"strings"
	"testing"
	"time"
)

func Test_parseHugoMeta(t *testing.T) {
	testCases := []struct {
		line string
		name string
		body string
	}{
		{
			line: "title = \"タイトルのテストです\"",
			name: "title",
			body: "\"タイトルのテストです\"",
		},
		{
			line: "tags = [\"Go\", \"Flutter\"]",
			name: "tags",
			body: "[\"Go\", \"Flutter\"]",
		},
	}
	for _, tc := range testCases {
		n, b := parseHugoMeta(tc.line)
		if n != tc.name {
			t.Errorf("got %s, expect %s", n, tc.name)
		}
		if b != tc.body {
			t.Errorf("got %s, expect %s", b, tc.body)
		}
	}
}

func Test_readHugoItem(t *testing.T) {
	testCases := []struct {
		body string
		item HugoItem
	}{
		{
			body: `
+++ 
date = "2019-03-15"
title = "Flutter StrutStyleで日本語と英語のTextの高さを揃える"
slug = "qiita-467936cffe857ebf35e9" 
tags = ["Flutter"]
categories = []
+++

こんにちは！こんにちは！
`,
			item: HugoItem{
				Title:   "Flutter StrutStyleで日本語と英語のTextの高さを揃える",
				Slug:    "qiita-467936cffe857ebf35e9",
				Content: "",
				Tags:    []string{"Flutter"},
				Date:    time.Date(2019, time.March, 15, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	for _, tc := range testCases {
		r := strings.NewReader(tc.body)
		item := readHugoItem(r)
		if item.Title != tc.item.Title {
			t.Errorf("title: got %s, expect %s", item.Title, tc.item.Title)
		}
		if item.Slug != tc.item.Slug {
			t.Errorf("slug: got %s, expect %s", item.Slug, tc.item.Slug)
		}
		if len(item.Tags) != len(tc.item.Tags) || item.Tags[0] != tc.item.Tags[0] {
			t.Errorf("tags: got %s, expect %s", item.Tags, tc.item.Tags)
		}
		if item.Date.String() != tc.item.Date.String() {
			t.Errorf("date: got %s, expect %s", item.Date.String(), tc.item.Date.String())
		}
	}
}

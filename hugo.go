package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HugoItem struct {
	Title      string
	Slug       string
	Content    string
	Tags       []string
	Categories []string
	Date       time.Time
}

func readHugoItems(dir string) []*HugoItem {
	var items []*HugoItem
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		item := readHugoItem(f)
		log.Printf("%s %s", path, item.Title)
		items = append(items, item)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return items
}

func readHugoItem(r io.Reader) *HugoItem {
	var (
		startMeta bool
		endMeta   bool
		content   strings.Builder
		item      HugoItem
	)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if !startMeta {
			if strings.HasPrefix(line, "+++") {
				startMeta = true
			}
		} else if !endMeta {
			if strings.HasPrefix(line, "+++") {
				endMeta = true
			} else {
				name, body := parseHugoMeta(line)
				//log.Println(name, body)
				switch name {
				case "title":
					item.Title = decodeHugoMetaString(body)
				case "slug":
					item.Slug = decodeHugoMetaString(body)
				case "date":
					item.Date = decodeHugoMetaTime(body)
				case "tags":
					item.Tags = decodeHugoMetaStrings(body)
				case "categories":
					item.Categories = decodeHugoMetaStrings(body)
				}
			}
		} else {
			content.WriteString(line)
			content.WriteRune('\n')
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	item.Content = content.String()
	return &item
}

func parseHugoMeta(line string) (string, string) {
	var (
		separated bool
		name      strings.Builder
		body      strings.Builder
	)

	r := strings.NewReader(line)
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			n := name.String()
			n = strings.TrimSpace(n)
			n = strings.ToLower(n)

			b := body.String()
			b = strings.TrimSpace(b)
			return n, b
		}

		if separated {
			body.WriteRune(ch)
		} else {
			if ch == '=' {
				separated = true
			} else {
				name.WriteRune(ch)
			}
		}
	}
}

func decodeHugoMetaString(body string) string {
	return strings.Trim(body, "\"")
}

func decodeHugoMetaStrings(body string) []string {
	var ret []string
	if err := json.Unmarshal([]byte(body), &ret); err != nil {
		panic(err)
	}
	return ret
}

func decodeHugoMetaTime(body string) time.Time {
	body = decodeHugoMetaString(body)
	t, err := time.Parse("2006-01-02", body)
	if err != nil {
		panic(err)
	}
	return t
}

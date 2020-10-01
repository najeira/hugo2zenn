package main

import (
	"flag"
	"log"
)

var (
	flagHugoDir = flag.String("hugo", "/Users/najeira/Projects/najeira.com/content/posts", "hugo dir")
	flagZennDir = flag.String("zenn", "/Users/najeira/Projects/zenn-dev/articles", "zenn dir")
)

func main() {
	log.SetFlags(0)
	flag.Parse()
	items := readHugoItems(*flagHugoDir)
	convertHugoItems(*flagZennDir, items)
}

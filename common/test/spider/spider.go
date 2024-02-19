package spider

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func HttpParse() {
	uri := "https://www.hudiemi.com/hudiedaquan/hudie_238.html"
	// 创建一个自定义的 http.Client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	// 发送 GET 请求
	resp, err := client.Get(uri)
	if err != nil {
		log.Println(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	printHTML(doc, 0)
}

func printHTML(n *html.Node, depth int) {
	if n.Type == html.ElementNode {
		fmt.Printf("%*s<%s>", depth, "", n.Data)
	} else if n.Type == html.TextNode {
		data := strings.TrimSpace(n.Data)
		fmt.Printf("%*s%s", depth, "", AllSpaceReplacer.Replace(data))
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		printHTML(c, depth+1)
	}

	if n.Type == html.ElementNode {
		fmt.Printf("%*s</%s>\n", depth, "", n.Data)
	}
}

var AllSpaceReplacer *strings.Replacer

func init() {
	commonSpaceChar := []string{" ", "\t", "\n", "\r", "\v", "\f"}
	var replacePairs = make([]string, 0)
	for _, v := range commonSpaceChar {
		replacePairs = append(replacePairs, v)
		replacePairs = append(replacePairs, "")
	}
	AllSpaceReplacer = strings.NewReplacer(replacePairs...)
}

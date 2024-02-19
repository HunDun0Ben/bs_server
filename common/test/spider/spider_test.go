package spider_test

import (
	"demo/common/test/spider"
	"fmt"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	fmt.Printf("aa = %s\n", strings.TrimSpace("\n"))
	spider.HttpParse()
}

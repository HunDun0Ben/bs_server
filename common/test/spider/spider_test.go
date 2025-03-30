package spider_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/HunDun0Ben/bs_server/common/test/spider"
)

func Test(t *testing.T) {
	fmt.Printf("aa = %s\n", strings.TrimSpace("\n"))
	spider.HTTPParse()
}

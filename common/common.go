package common

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/muesli/cache2go"
)

type AnyMaps map[string]interface{}

// 符合规则的文章链接
type Article struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	SiteURL string `json:"site_url,omitempty"`
}

var Articles = cache2go.Cache("Articles")

var RootDir, _ = RunPath()
var DataDir = filepath.Join(RootDir, "data")

func RunPath() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return os.Getwd()
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}

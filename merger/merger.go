package merger

import (
	"fmt"
	"regexp"
	"strings"
)

type DataMerger interface {
	Merge(contentType string, body string, dataMap map[string]string) string
}

type dataMergerImpl struct {
	patterns map[string]*fieldPattern
}

type fieldPattern struct {
	plain *regexp.Regexp
	html  *regexp.Regexp
}

func New(fields []string) DataMerger {
	patterns := make(map[string]*fieldPattern)

	for _, fieldName := range fields {
		m1, err1 := regexp.Compile(`<<\s*` + fieldName + `\s*>>`)
		m2, err2 := regexp.Compile(`&lt;&lt;\s*` + fieldName + `\s*&gt;&gt;`)

		if err1 != nil {
			fmt.Printf("%s: ingore %v\n", fieldName, err1)
		} else if err2 != nil {
			fmt.Printf("%s: ingore %v\n", fieldName, err2)
		} else {
			patterns[fieldName] = &fieldPattern{
				plain: m1,
				html:  m2,
			}
		}
	}
	return &dataMergerImpl{
		patterns: patterns,
	}
}

func (impl *dataMergerImpl) Merge(contentType string, body string, dataMap map[string]string) string {
	html := strings.Contains(contentType, "html")

	content := body
	for fieldName, p := range impl.patterns {
		if html {
			content = p.html.ReplaceAllString(content, dataMap[fieldName])
		} else {
			content = p.plain.ReplaceAllString(content, dataMap[fieldName])
		}
	}
	return content
}

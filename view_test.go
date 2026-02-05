package goview

import (
	"bytes"
	"html/template"
	"strings"
	"testing"
	"time"
)

var cases = []struct {
	Name string
	Data M
	Out  string
}{
	{
		Name: "echo.tpl",
		Data: M{"name": "GoView"},
		Out:  "$GoView",
	},
	{
		Name: "default.tpl",
		Data: M{"v": nil},
		Out:  "X",
	},
	{
		Name: "dict.tpl",
		Data: M{},
		Out:  "1-2",
	},
	{
		Name: "list.tpl",
		Data: M{},
		Out:  "2",
	},
	{
		Name: "urlquery.tpl",
		Data: M{"q": "a b"},
		Out:  "<a href=\"/search?q=a&#43;b\">search</a>",
	},
	{
		Name: "date.tpl",
		Data: M{"ts": time.Unix(0, 0).UTC()},
		Out:  "1970-01-01",
	},
	{
		Name: "json.tpl",
		Data: M{"v": struct {
			A int `json:"a"`
		}{A: 1}},
		Out: "<script>var x={\"a\":1};</script>",
	},
	{
		Name: "escape.tpl",
		Data: M{"s": "<a>"},
		Out:  "&lt;a&gt;",
	},
	{
		Name: "safehtml.tpl",
		Data: M{"h": "<b>x</b>"},
		Out:  "<b>x</b>",
	},
	{
		Name: "include",
		Data: M{"name": "GoView"},
		Out:  "<v>IncGoView</v>",
	},
	{
		Name: "index",
		Data: M{},
		Out:  "<v>Index</v>",
	},
	{
		Name: "sum",
		Data: M{
			"sum": func(a int, b int) int {
				return a + b
			},
			"a": 1,
			"b": 2,
		},
		Out: "<v>3</v>",
	},
}

func TestViewEngine_RenderWriter(t *testing.T) {
	gv := New(&Config{
		Root:      "_examples/test",
		Extension: ".tpl",
		Master:    "layouts/master",
		Partials:  []string{},
		Funcs: template.FuncMap{
			"echo": func(v string) string {
				return "$" + v
			},
		},
		DisableCache: true,
	})

	for _, v := range cases {
		buff := new(bytes.Buffer)
		err := gv.RenderWriter(buff, v.Name, v.Data)
		if err != nil {
			t.Errorf("name: %v, data: %v, error: %v", v.Name, v.Data, err)
			continue
		}
		val := strings.TrimSpace(buff.String())
		if val != v.Out {
			t.Errorf("actual: %v, expect: %v", val, v.Out)
		}
	}
}

func TestViewEngine_RenderWriter_BuiltinFuncOverride(t *testing.T) {
	gvr := New(&Config{
		Root:      "_examples/test",
		Extension: ".tpl",
		Master:    "",
		Partials:  []string{},
		Funcs: template.FuncMap{
			"default": func(d any, v any) any {
				return "OVERRIDE"
			},
		},
		DisableCache: true,
	})

	buff := new(bytes.Buffer)
	err := gvr.RenderWriter(buff, "override_default.tpl", M{"v": nil})
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	if got := strings.TrimSpace(buff.String()); got != "OVERRIDE" {
		t.Fatalf("actual: %v, expect: %v", got, "OVERRIDE")
	}
}

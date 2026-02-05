package goview

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"reflect"
	"time"
)

func builtinFuncMap() template.FuncMap {
	return template.FuncMap{
		"default": func(d any, v any) any {
			if isEmptyValue(v) {
				return d
			}
			return v
		},
		"list": func(values ...any) []any {
			return values
		},
		"dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict expects an even number of arguments")
			}
			m := make(map[string]any, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				k, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				m[k] = values[i+1]
			}
			return m, nil
		},
		"urlquery": func(v any) template.URL {
			return template.URL(url.QueryEscape(fmt.Sprint(v)))
		},
		"date": func(v any, layout string) (string, error) {
			switch t := v.(type) {
			case time.Time:
				return t.Format(layout), nil
			case *time.Time:
				if t == nil {
					return "", nil
				}
				return t.Format(layout), nil
			case int64:
				return time.Unix(t, 0).Format(layout), nil
			case int:
				return time.Unix(int64(t), 0).Format(layout), nil
			default:
				return "", fmt.Errorf("date expects time.Time, *time.Time, int64, or int")
			}
		},
		"json": func(v any) (template.JS, error) {
			b, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return template.JS(b), nil
		},
		"escape": func(s any) template.HTML {
			return template.HTML(template.HTMLEscapeString(fmt.Sprint(s)))
		},
		"safeHTML": func(s any) template.HTML {
			return template.HTML(fmt.Sprint(s))
		},
	}
}

func isEmptyValue(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return true
	}
	for rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return true
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Bool:
		return !rv.Bool()
	case reflect.String:
		return rv.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Slice, reflect.Map, reflect.Array:
		return rv.Len() == 0
	case reflect.Struct:
		return rv.IsZero()
	default:
		return rv.IsZero()
	}
}

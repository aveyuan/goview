package ginview

import (
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/aveyuan/goview"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

const templateEngineKey = "foolin-goview-ginview"

// ViewEngine view engine for gin
type ViewEngine struct {
	*goview.ViewEngine
}

// ViewRender view render implement gin interface
type ViewRender struct {
	Engine *ViewEngine
	Name   string
	Data   interface{}
}

// New new view engine for gin
func New(config *goview.Config) *ViewEngine {
	return Wrap(goview.New(config))
}

// Wrap wrap view engine for goview.ViewEngine
func Wrap(engine *goview.ViewEngine) *ViewEngine {
	return &ViewEngine{
		ViewEngine: engine,
	}
}

// Default new default engine
func Default() *ViewEngine {
	return New(goview.DefaultConfig)
}

// Instance implement gin interface
func (e *ViewEngine) Instance(name string, data interface{}) render.Render {
	return ViewRender{
		Engine: e,
		Name:   name,
		Data:   data,
	}
}

// HTML render html
func (e *ViewEngine) HTML(ctx *gin.Context, code int, name string, data interface{}) {
	instance := e.Instance(name, data)
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Render(code, instance)
}

// Render (YAML) marshals the given interface object and writes data with custom ContentType.
func (v ViewRender) Render(w http.ResponseWriter) error {
	err := v.Engine.RenderWriter(w, v.Name, v.Data)
	if err == nil {
		return nil
	}

	fmt.Fprintf(gin.DefaultErrorWriter, "[goview] render %q error: %v\n", v.Name, err)

	if rw, ok := w.(gin.ResponseWriter); ok {
		if !rw.Written() {
			if rw.Status() < http.StatusBadRequest {
				rw.WriteHeader(http.StatusInternalServerError)
			}
		}
	}

	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = goview.HTMLContentType
	}

	escapedErr := html.EscapeString(err.Error())
	_, _ = io.WriteString(w, "<!doctype html><html><head><meta charset=\"utf-8\"><title>Render Error</title></head><body><h1>Render Error</h1><pre>")
	_, _ = io.WriteString(w, escapedErr)
	_, _ = io.WriteString(w, "</pre></body></html>")

	return err
}

// WriteContentType write html content type
func (v ViewRender) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = goview.HTMLContentType
	}
}

// NewMiddleware gin middleware for func `gintemplate.HTML()`
func NewMiddleware(config *goview.Config) gin.HandlerFunc {
	return Middleware(New(config))
}

// Middleware gin middleware wrapper
func Middleware(e *ViewEngine) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := c.Get(templateEngineKey); !ok {
			c.Set(templateEngineKey, e)
		}
		c.Next()
	}
}

// HTML html render for template
// You should use helper func `Middleware()` to set the supplied
// TemplateEngine and make `HTML()` work validly.
func HTML(ctx *gin.Context, code int, name string, data interface{}) {
	if val, ok := ctx.Get(templateEngineKey); ok {
		if e, ok := val.(*ViewEngine); ok {
			e.HTML(ctx, code, name, data)
			return
		}
	}
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.HTML(code, name, data)
}

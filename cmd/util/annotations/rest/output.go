package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

const (
	handlerPackage = `package %s
`
	handlerImports = `
import (
	"github.com/xiusin/pine"
)
`
	handlerTemplate = `
func (c *{{ .Receiver }}) RegisterRoute(b pine.IRouterWrapper) {
	{{range .Handlers}}b.{{ .Method }}("{{ .Route }}", "{{ .Handler }}"){{end}}
}
`
)

type outputKey struct {
	pkg string
	dir string
}

type handlerTemplateMethods struct {
	Method  string
	Route   string
	Handler string
}
type handlerTemplateParams struct {
	Receiver string
	Path     string
	Handlers []handlerTemplateMethods
}

type output struct {
	files map[outputKey]*bytes.Buffer
	tmp   *template.Template
}

func newOutput() (*output, error) {
	tmp, err := template.New("handler").Parse(handlerTemplate)
	if err != nil {
		return nil, err
	}
	return &output{
		files: make(map[outputKey]*bytes.Buffer),
		tmp:   tmp,
	}, nil

}

func (o *output) append(meta handlerMetadata, handlers handlerMapping) error {
	outKey := outputKey{
		pkg: meta.pkg,
		dir: meta.dir,
	}
	b := o.upsertBuffer(outKey)

	htm := []handlerTemplateMethods{}
	for m, h := range handlers.mapping {
		ss := strings.Split(h, sp)
		htm = append(htm, handlerTemplateMethods{Method: httpMethodConstant(m), Handler: ss[0], Route: ss[1]})
	}
	toTemplate := handlerTemplateParams{
		Receiver: meta.structName,
		Path:     handlers.path,
		Handlers: htm,
	}

	err := o.tmp.Execute(b, toTemplate)
	if err != nil {
		return err
	}
	return nil
}

func (o *output) get() map[string][]byte {
	out := map[string][]byte{}
	for k, buf := range o.files {
		out[k.dir+"/"+".handlers.gen.go"] = buf.Bytes()
	}
	return out
}

func (o *output) upsertBuffer(k outputKey) *bytes.Buffer {
	if b, ok := o.files[k]; ok {
		return b
	}
	b := bytes.NewBufferString(fmt.Sprintf(handlerPackage, k.pkg) + handlerImports + "\n")
	o.files[k] = b
	return b
}

func httpMethodConstant(m string) string {
	return map[string]string{
		http.MethodGet:     "GET",
		http.MethodHead:    "HEAD",
		http.MethodPost:    "POST",
		http.MethodPut:     "PUT",
		http.MethodPatch:   "PATCH",
		http.MethodDelete:  "DELETE",
		http.MethodConnect: "CONNECT",
		http.MethodOptions: "OPTIONS",
		http.MethodTrace:   "TRACE",
	}[m]
}

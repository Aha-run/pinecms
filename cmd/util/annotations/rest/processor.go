package rest

import (
	"fmt"
	"go/ast"
	"net/http"

	annotation "github.com/xiusin/go-annotation/pkg"
)

type Rest struct {
	Method string `annotation:"name=method,default=GET"`
	Path   string `annotation:"name=path,default=/"`
	Route  string `annotation:"name=route"`
}

const sp = `|||-|||`

func init() {
	annotation.Register[Rest](&Processor{
		mapping: make(map[handlerMetadata]handlerMapping),
	})
}

var _ annotation.AnnotationProcessor = (*Processor)(nil)

type Processor struct {
	mapping map[handlerMetadata]handlerMapping
}

func (p *Processor) Process(node annotation.Node) error {
	annotations := annotation.FindAnnotations[Rest](node.Annotations())
	if len(annotations) == 0 {
		return nil
	}

	if len(annotations) > 1 {
		return fmt.Errorf("expected 1 rest annotation, but got: %d", len(annotations))
	}

	n := node.ASTNode()
	switch nt := n.(type) {
	case *ast.TypeSpec:
		return p.processStructure(annotations[0], node, nt)
	case *ast.FuncDecl:
		return p.processMethod(annotations[0], node, nt)
	default:
		return fmt.Errorf("unexpected node type %T - %t", n, n)
	}
}

func (p *Processor) Version() string {
	return "0.0.1"
}

func (p *Processor) Name() string {
	return "Rest"
}

type handlerMetadata struct {
	pkg        string
	structName string
	dir        string
	fileName   string
}

type handlerMapping struct {
	path    string
	mapping map[string]string
}

func (p *Processor) processStructure(rest Rest, node annotation.Node, s *ast.TypeSpec) error {
	key := handlerMetadata{
		pkg:        node.Meta().PackageName(),
		structName: s.Name.Name,
		dir:        node.Meta().Dir(),
		fileName:   node.Meta().FileName(),
	}

	_, ok := p.mapping[key]
	if !ok {
		p.mapping[key] = handlerMapping{
			path:    rest.Path,
			mapping: make(map[string]string),
		}
	}

	return nil
}

func (p *Processor) processMethod(rest Rest, node annotation.Node, f *ast.FuncDecl) error {
	if !p.validateHTTPMethod(rest.Method) {
		return fmt.Errorf("invalid HTTP method: %s", rest.Method)
	}

	recvName := MethodReceiver(f)
	if recvName == "" {
		return fmt.Errorf("expected method receiver, but got empty for %s", f.Name.Name)
	}

	key := handlerMetadata{
		pkg:        node.Meta().PackageName(),
		structName: recvName,
		dir:        node.Meta().Dir(),
		fileName:   node.Meta().FileName(),
	}

	v, ok := p.mapping[key]
	if !ok {
		return fmt.Errorf("no mapping for %s", key)
	}

	if len(rest.Route) == 0 {
		rest.Route = "/" + f.Name.Name
	}

	v.mapping[rest.Method] = f.Name.Name + sp + rest.Route

	return nil
}

func (p *Processor) validateHTTPMethod(m string) bool {
	_, ok := map[string]struct{}{
		http.MethodGet:     {},
		http.MethodHead:    {},
		http.MethodPost:    {},
		http.MethodPut:     {},
		http.MethodPatch:   {},
		http.MethodDelete:  {},
		http.MethodConnect: {},
		http.MethodOptions: {},
		http.MethodTrace:   {},
	}[m]
	return ok
}

func (p *Processor) Output() map[string][]byte {
	o, err := newOutput()
	if err != nil {
		panic(err)
	}
	for k, v := range p.mapping {
		err = o.append(k, v)
		if err != nil {
			panic(err)
		}
	}

	return o.get()
}

func MethodReceiver(decl *ast.FuncDecl) string {
	if decl.Recv == nil {
		return ""
	}

	for _, v := range decl.Recv.List {
		switch rv := v.Type.(type) {
		case *ast.StarExpr:
			return rv.X.(*ast.Ident).Name
		case *ast.UnaryExpr:
			return rv.X.(*ast.Ident).Name
		}
	}
	return ""
}

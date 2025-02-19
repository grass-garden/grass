package router

import (
	"net/http"
	"sync"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

var _ http.Handler = (*Router)(nil)

type Router struct {
	once sync.Once

	mux         *http.ServeMux
	pattern     string
	middlewares []Middleware
	routes      []Route

	doc                *v3.Document
	serializers        map[string]Serializer
	errorProcessor     ErrorProcessor
	methodToStatusCode MethodToStatusCode
	contentType        string

	enableAutoSlash bool
}

type (
	Option                        func(*Router)
	Middleware                    func(*ContextAny) error
	Handler[I, O any, Ctx ctx[I]] func(Ctx) (O, error)
)

func New() *Router {
	return &Router{
		once: sync.Once{},

		pattern:     "",
		mux:         http.NewServeMux(),
		middlewares: make([]Middleware, 0),
		routes:      make([]Route, 0),

		doc:                defaultSchema(),
		serializers:        defaultSerializers(),
		contentType:        contentTypeJson,
		errorProcessor:     defaultErrorProcessor,
		methodToStatusCode: defaultMethodToStatusCode,

		enableAutoSlash: false,
	}
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	r.once.Do(func() {
		for _, ro := range r.routes {
			r.mux.Handle(ro.MuxPattern(), ro)
		}
	})

	r.mux.ServeHTTP(res, req)
}

func (r *Router) Schema() *v3.Document {
	for _, ro := range r.routes {
		item, ok := r.doc.Paths.PathItems.Get(ro.Pattern())
		if !ok {
			item = &v3.PathItem{}
		}

		switch ro.Method() {
		case http.MethodGet:
			item.Get = ro.Operation()
		case http.MethodHead:
			item.Head = ro.Operation()
		case http.MethodPost:
			item.Post = ro.Operation()
		case http.MethodPut:
			item.Put = ro.Operation()
		case http.MethodPatch:
			item.Patch = ro.Operation()
		case http.MethodDelete:
			item.Delete = ro.Operation()
		default:
			continue
		}

		r.doc.Paths.PathItems.Set(ro.Pattern(), item)
	}

	return r.doc
}

func Get[Input, Output any, Ctx ctx[Input]](
	r *Router,
	pattern string,
	handler Handler[Input, Output, Ctx],
) *route[Input, Output, Ctx] {
	route := newRoute(r, http.MethodGet, pattern, handler)
	r.routes = append(r.routes, route)
	return route
}

func Head[Input, Output any, Ctx ctx[Input]](
	r *Router,
	pattern string,
	handler Handler[Input, Output, Ctx],
) *route[Input, Output, Ctx] {
	route := newRoute(r, http.MethodHead, pattern, handler)
	r.routes = append(r.routes, route)
	return route
}

func Post[Input, Output any, Ctx ctx[Input]](
	r *Router,
	pattern string,
	handler Handler[Input, Output, Ctx],
) *route[Input, Output, Ctx] {
	route := newRoute(r, http.MethodPost, pattern, handler)
	r.routes = append(r.routes, route)
	return route
}

func Put[Input, Output any, Ctx ctx[Input]](
	r *Router,
	pattern string,
	handler Handler[Input, Output, Ctx],
) *route[Input, Output, Ctx] {
	route := newRoute(r, http.MethodPut, pattern, handler)
	r.routes = append(r.routes, route)
	return route
}

func Patch[Input, Output any, Ctx ctx[Input]](
	r *Router,
	pattern string,
	handler Handler[Input, Output, Ctx],
) *route[Input, Output, Ctx] {
	route := newRoute(r, http.MethodPatch, pattern, handler)
	r.routes = append(r.routes, route)
	return route
}

func Delete[Input, Output any, Ctx ctx[Input]](
	r *Router,
	pattern string,
	handler Handler[Input, Output, Ctx],
) *route[Input, Output, Ctx] {
	route := newRoute(r, http.MethodDelete, pattern, handler)
	r.routes = append(r.routes, route)
	return route
}

func Any[Input, Output any, Ctx ctx[Input]](
	r *Router,
	pattern string,
	handler Handler[Input, Output, Ctx],
) *route[Input, Output, Ctx] {
	route := newRoute(r, "", pattern, handler)
	r.routes = append(r.routes, route)
	return route
}

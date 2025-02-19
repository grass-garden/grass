package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gobeam/stringy"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

const (
	xContentType    = "Content-Type"
	contentTypeJson = "application/json"
)

var (
	_ Route              = (*route[any, any, ctx[any]])(nil)
	_ http.Handler       = (*route[any, any, ctx[any]])(nil)
	_ MethodToStatusCode = defaultMethodToStatusCode
)

type Route interface {
	http.Handler
	Use(...Middleware) Route
	Operation() *v3.Operation
	Method() string
	Pattern() string
	MuxPattern() string
}

type route[Input, Output any, Ctx ctx[Input]] struct {
	method      string
	pattern     string
	middlewares []Middleware
	handler     Handler[Input, Output, Ctx]
	router      *Router

	contentType    string
	serializer     Serializer
	errorProcessor ErrorProcessor

	operationId string
	summary     string
	description string
	statusCode  int
}

func newRoute[Input, Output any, Ctx ctx[Input]](
	router *Router,
	method string,
	pattern string,
	handler Handler[Input, Output, Ctx],
) *route[Input, Output, Ctx] {
	if router.pattern != "" {
		if pattern != "/" {
			pattern = router.pattern + pattern
		} else {
			pattern = router.pattern
		}
	}

	if pattern[0] != '/' {
		pattern = "/" + pattern
	}

	if router.enableAutoSlash && pattern[len(pattern)-1] != '/' {
		pattern = pattern + "/"
	}

	var output Output
	var serializer Serializer
	contentType := router.contentType
	if v, ok := any(output).(Serializer); ok {
		serializer = v
		contentType = v.ContentType()
	} else if v, ok := router.serializers[contentType]; ok {
		serializer = v
	} else {
		serializer = JSONSerializer{}
		contentType = contentTypeJson
	}

	errorProcessor := router.errorProcessor
	if errorProcessor == nil {
		errorProcessor = defaultErrorProcessor
	}

	// begin openapi
	chars := []string{"/", " ", "{", " by ", "}", " ", "*", ""}
	operationId := stringy.New(method + " " + pattern).KebabCase(chars...).ToLower()
	summary := strings.ReplaceAll(operationId, "-", " ")
	description := summary

	return &route[Input, Output, Ctx]{
		method:      method,
		pattern:     pattern,
		handler:     handler,
		router:      router,
		middlewares: router.middlewares,

		serializer:     serializer,
		contentType:    contentType,
		errorProcessor: errorProcessor,

		summary:     summary,
		description: description,
		operationId: operationId,
		statusCode:  router.methodToStatusCode(method),
	}
}

func (r *route[Input, Output, Ctx]) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctxAny := &ContextAny{
		req:          req,
		res:          res,
		statusCode:   r.statusCode,
		isNextCalled: true,
	}

	ctx := newContext[Input, Ctx](ctxAny)

	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(error); ok {
				r.handleError(ctx, e)
			} else {
				r.handleError(ctx, fmt.Errorf("panic recovered: %v", err))
			}
		}
	}()

	for _, middleware := range r.middlewares {
		ctxAny.isNextCalled = false
		if err := middleware(ctxAny); err != nil || !ctxAny.isNextCalled {
			r.handleError(ctx, err)
			return
		}
	}

	output, err := r.handler(ctx)
	if err != nil {
		r.handleError(ctx, err)
		return
	}

	ctx.ResponseWriter().WriteHeader(ctxAny.statusCode)
	ctx.SetHeader(xContentType, r.contentType)
	_ = r.serializer.Marshal(ctx.ResponseWriter(), output)
}

func (r *route[Input, Output, Ctx]) Use(middlewares ...Middleware) Route {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

func (r *route[Input, Output, Ctx]) Method() string {
	return r.method
}

func (r *route[Input, Output, Ctx]) Pattern() string {
	return r.pattern
}

func (r *route[Input, Output, Ctx]) MuxPattern() string {
	if r.method == "" {
		return r.pattern
	}
	return r.method + " " + r.pattern
}

func (r *route[Input, Output, Ctx]) handleError(ctx Ctx, err error) {
	statusCode := http.StatusInternalServerError
	err = r.router.errorProcessor(err)
	if v, ok := err.(Error); ok {
		statusCode = v.StatusCode()
	}

	res := ctx.ResponseWriter()
	res.WriteHeader(statusCode)
	ctx.SetHeader(xContentType, r.contentType)
	_ = r.serializer.Marshal(res, err)
}

type MethodToStatusCode func(string) int

func defaultMethodToStatusCode(method string) int {
	switch method {
	case http.MethodPost:
		return http.StatusCreated
	case http.MethodPut, http.MethodPatch, http.MethodDelete:
		return http.StatusAccepted
	default:
		return http.StatusOK
	}
}

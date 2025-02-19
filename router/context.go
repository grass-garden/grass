package router

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	_ ctx[any]        = (*Context[any])(nil)
	_ ctx[any]        = (*ContextAny)(nil)
	_ context.Context = (*ContextAny)(nil)
)

type ctx[T any] interface {
	Next() error

	Request() *http.Request
	ResponseWriter() http.ResponseWriter

	Body() T
	GetBody() (T, error)
	BodyRaw() io.ReadCloser

	Header(key string) string
	SetHeader(key, value string)
	AddHeader(key, value string)

	ResponseHeader(key string) string

	PathParam(key string) string
	SetPathParam(key, value string)

	QueryParam(key string) string
	SetQueryParam(key, value string)

	Status() int
	SetStatus(code int)
}

type ContextAny struct {
	res          http.ResponseWriter
	req          *http.Request
	statusCode   int
	isNextCalled bool
}

type Context[Body any] struct {
	*ContextAny
	body *Body
}

func newContext[Body any, Ctx ctx[Body]](ctxAny *ContextAny) Ctx {
	var c Ctx
	switch any(c).(type) {
	case *ContextAny:
		return any(ctxAny).(Ctx)
	case *Context[Body]:
		return any(&Context[Body]{ContextAny: ctxAny}).(Ctx)
	default:
		panic("invalid context struct")
	}
}

// Deadline implements context.Context.
func (ctx *ContextAny) Deadline() (deadline time.Time, ok bool) {
	return ctx.req.Context().Deadline()
}

// Done implements context.Context.
func (ctx *ContextAny) Done() <-chan struct{} {
	return ctx.req.Context().Done()
}

// Err implements context.Context.
func (ctx *ContextAny) Err() error {
	return ctx.req.Context().Err()
}

// Value implements context.Context.
func (ctx *ContextAny) Value(key any) any {
	return ctx.req.Context().Value(key)
}

func (ctx *ContextAny) Header(key string) string {
	return ctx.Request().Header.Get(key)
}

func (ctx *ContextAny) SetHeader(key, value string) {
	ctx.ResponseWriter().Header().Set(key, value)
}

func (ctx *ContextAny) AddHeader(key, value string) {
	ctx.ResponseWriter().Header().Add(key, value)
}

func (ctx *ContextAny) ResponseHeader(key string) string {
	return ctx.ResponseWriter().Header().Get(key)
}

func (ctx *ContextAny) PathParam(key string) string {
	return ctx.Request().PathValue(key)
}

func (ctx *ContextAny) SetPathParam(key, value string) {
	ctx.Request().SetPathValue(key, value)
}

func (ctx *ContextAny) QueryParam(key string) string {
	return ctx.Request().URL.Query().Get(key)
}

func (ctx *ContextAny) SetQueryParam(key, value string) {
	query := ctx.Request().URL.Query()
	query.Set(key, value)
	ctx.Request().URL.RawQuery = query.Encode()
}

func (ctx *ContextAny) Status() int {
	return ctx.statusCode
}

func (ctx *ContextAny) SetStatus(code int) {
	ctx.statusCode = code
}

func (ctx *ContextAny) Body() any {
	panic("*ContextAny.Body should not be called, use *ContextAny.BodyRaw")
}

func (ctx *ContextAny) GetBody() (any, error) {
	panic("*ContextAny.GetBody should not be called, use *ContextAny.BodyRaw")
}

func (ctx *ContextAny) BodyRaw() io.ReadCloser {
	return ctx.Request().Body
}

func (ctx *ContextAny) Next() error {
	ctx.isNextCalled = true
	return nil
}

func (ctx *ContextAny) Request() *http.Request {
	return ctx.req
}

func (ctx *ContextAny) ResponseWriter() http.ResponseWriter {
	return ctx.res
}

func (ctx *Context[Body]) Body() Body {
	body, err := ctx.GetBody()
	if err != nil {
		panic(err)
	}
	return body
}

func (ctx *Context[Body]) GetBody() (Body, error) {
	if ctx.body != nil {
		return *ctx.body, nil
	}

	body := new(Body)
	if ctx.req.ContentLength > 0 && ctx.req.Body != http.NoBody {
		if err := json.NewDecoder(ctx.req.Body).Decode(body); err != nil {
			return *body, fmt.Errorf("could not read incoming request")
		}
	}

	// if err := validator.Default().Struct(body); err != nil {
	// 	// _ = err.(validator.ValidationErrors)
	// 	return *body, UnprocessableEntityError{Err: err}
	// }

	ctx.body = body
	return *ctx.body, nil
}

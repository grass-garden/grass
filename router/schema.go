package router

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	typetostring "github.com/samber/go-type-to-string"
	"go.grass.garden/utils"
)

func (r *route[Input, Output, Ctx]) Operation() *v3.Operation {
	operation := &v3.Operation{
		OperationId: r.operationId,
		Summary:     r.summary,
		Description: r.description,
	}

	// Input
	inputType := reflect.TypeOf((*Input)(nil)).Elem()
	operation.RequestBody = &v3.RequestBody{
		Required:    utils.ToPointer(true),
		Description: http.StatusText(r.statusCode),
		Content: orderedmap.FromPairs(
			orderedmap.NewPair(
				r.contentType, &v3.MediaType{
					Schema: walk(r.router.doc, operation, inputType),
				},
			),
		),
	}

	// Output
	outputType := reflect.TypeOf((*Output)(nil)).Elem()
	operation.Responses = &v3.Responses{
		Codes: orderedmap.FromPairs(
			orderedmap.NewPair(strconv.Itoa(r.statusCode), &v3.Response{
				Content: orderedmap.FromPairs(orderedmap.NewPair(r.contentType, &v3.MediaType{
					Schema: walk(r.router.doc, operation, outputType),
				})),
			}),
		),
		Default: &v3.Response{
			Content: orderedmap.FromPairs(orderedmap.NewPair(r.contentType, &v3.MediaType{
				Schema: walk(r.router.doc, operation, reflect.TypeOf((*HTTPError)(nil)).Elem()),
			})),
		},
	}

	return operation
}

func walk(doc *v3.Document, op *v3.Operation, t reflect.Type) *base.SchemaProxy {
	s := &base.Schema{}
	proxy := base.CreateSchemaProxy(s)
	s.Properties = orderedmap.New[string, *base.SchemaProxy]()

	switch t.Kind() {
	case reflect.Bool:
		s.Type = append(s.Type, "boolean")
	case reflect.String:
		s.Type = append(s.Type, "string")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s.Type = append(s.Type, "integer")
		if t.Kind() != reflect.Int {
			s.Format = t.Kind().String()
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s.Type = append(s.Type, "integer")
		s.Format = t.Kind().String()
		s.Minimum = utils.ToPointer(0.0)
	case reflect.Float32:
		s.Type = append(s.Type, "number")
		s.Format = "float"
	case reflect.Float64:
		s.Type = append(s.Type, "number")
		s.Format = "double"
	case reflect.Array:
		s.Type = append(s.Type, "array")
		s.MinItems = utils.ToPointer(int64(t.Len()))
		s.MaxItems = utils.ToPointer(int64(t.Len()))
		s.Items = &base.DynamicValue[*base.SchemaProxy, bool]{
			A: walk(doc, op, t.Elem()),
		}
	case reflect.Slice:
		s.Type = append(s.Type, "array")
		s.Items = &base.DynamicValue[*base.SchemaProxy, bool]{
			A: walk(doc, op, t.Elem()),
		}
	case reflect.Map:
		s.Type = append(s.Type, "object")
		s.AdditionalProperties = &base.DynamicValue[*base.SchemaProxy, bool]{
			A: walk(doc, op, t.Elem()),
		}
	case reflect.Ptr:
		return walk(doc, op, t.Elem())
	case reflect.Interface:
		s.Type = append(s.Type, "object")
		s.AdditionalProperties = &base.DynamicValue[*base.SchemaProxy, bool]{
			N: 1,
			B: true,
		}
	case reflect.Struct:
		sft := typeToString(t)
		_, present := doc.Components.Schemas.Get(sft)
		if present {
			return base.CreateSchemaProxyRef(componentSchemaRef(t))
		}

		s.Type = append(s.Type, "object")
		for i := 0; i < t.NumField(); i++ {
			if f := t.Field(i); f.IsExported() {
				nsp := walk(doc, op, f.Type)
				params := structPropToParams(f, nsp)
				if len(params) > 0 {
					op.Parameters = append(op.Parameters, params...)
				} else if name, skip := propName(f); !skip {
					s.Properties.Set(name, nsp)
				}
			}
		}

		if s.Properties.Len() == 0 {
			return nil
		}

		doc.Components.Schemas.Set(sft, proxy)
		return base.CreateSchemaProxyRef(componentSchemaRef(t))
	}

	return proxy
}

func structPropToParams(sf reflect.StructField, schema *base.SchemaProxy) (params []*v3.Parameter) {
	if v := sf.Tag.Get("header"); v != "" {
		params = append(params, &v3.Parameter{
			Name:   v,
			In:     "header",
			Schema: schema,
		})
	}
	if v := sf.Tag.Get("path"); v != "" {
		params = append(params, &v3.Parameter{
			Name:   v,
			In:     "path",
			Schema: schema,
		})
	}
	if v := sf.Tag.Get("query"); v != "" {
		params = append(params, &v3.Parameter{
			Name:            v,
			In:              "query",
			Schema:          schema,
			AllowEmptyValue: true,
		})
	}
	return params
}

func propName(sf reflect.StructField) (string, bool) {
	tags := []string{"json"}
	for _, tag := range tags {
		if name := sf.Tag.Get(tag); name != "" {
			before, _, _ := strings.Cut(name, ",")
			return before, before == "-"
		}
	}
	return sf.Name, false
}

func componentSchemaRef(t reflect.Type) string {
	return "#/components/schemas/" + typeToString(t)
}

func typeToString(t reflect.Type) string {
	return typetostring.GetReflectType(t)
}

func defaultSchema() *v3.Document {
	return &v3.Document{
		Version: "3.1.0",
		Info: &base.Info{
			Title:   "Openapi Schema",
			Version: "0.1.0",
		},
		Paths: &v3.Paths{
			PathItems: orderedmap.New[string, *v3.PathItem](),
		},
		Components: &v3.Components{
			Schemas: orderedmap.New[string, *base.SchemaProxy](),
		},
	}
}

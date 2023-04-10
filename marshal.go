package palm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/DWVoid/calm"
)

type IMarshalErr interface {
	Path() string
	Reason() string
}

type _MarshalErr struct {
	path   []string
	reason string
}

func (o *_MarshalErr) Reason() string { return o.reason }

func (o *_MarshalErr) Path() string {
	if len(o.path) == 0 {
		return "<root>"
	}
	b := strings.Builder{}
	for i := len(o.path) - 1; i >= 0; i-- {
		b.WriteString(o.path[i])
		if i != 0 {
			b.WriteRune('.')
		}
	}
	return b.String()
}

func _Type2String(t int32) string {
	switch t {
	case TypeNull:
		return "<nil>"
	case TypeBool:
		return "<bool>"
	case TypeNum:
		return "<number>"
	case TypeStr:
		return "<string>"
	case TypeMap:
		return "<object>"
	case TypeArr:
		return "<array>"
	}
	return "<error_type>"
}

func NewMarshalErrType(want, actual int32) IMarshalErr {
	return NewMarshalErr(fmt.Sprintf("type error: expect %s got %s", _Type2String(want), _Type2String(actual)))
}

func NewMarshalErr(reason string) IMarshalErr {
	return &_MarshalErr{path: []string{}, reason: reason}
}

func WrapMarshalErr(err IMarshalErr, path string) IMarshalErr {
	e := err.(*_MarshalErr)
	e.path = append(e.path, path)
	return e
}

// IMarshalSchema The interface for a marshal-able type
//
// Note:
//
//	There has always been an interest to speed up encode/decode by skipping reflection.
//	However, this cannot be achieved without exposing golang's internal api.
//	I will stick with reflect.Value for this for the foreseeable future.
//	If anyone is interested in doing so, maybe checkout Masaaki Goshima(https://github.com/goccy)'s work
//	and implement said mechanism on your own fork
type IMarshalSchema interface {
	Zero(value reflect.Value)
	Read(value reflect.Value) ITree
	Write(value reflect.Value, tree ITree) IMarshalErr
}

func _TypeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func To[T any](tree ITree, ctx IMarshalSchemaContext) calm.ResultT[T] {
	return calm.RunT(func() T {
		var result T
		schema := ctx.With().Resolve(_TypeOf[T]())
		err := schema.Write(reflect.ValueOf(&result).Elem(), tree)
		if err != nil {
			calm.ThrowDetail(calm.ERequest, "schema", fmt.Sprintf("%s: %s", err.Path(), err.Reason()))
		}
		return result
	})
}

func From[T any](obj T, ctx IMarshalSchemaContext) calm.ResultT[ITree] {
	return calm.RunT(func() ITree {
		schema := ctx.With().Resolve(_TypeOf[T]())
		return schema.Read(reflect.ValueOf(&obj).Elem())
	})
}

func ToDefault[T any](tree ITree) calm.ResultT[T] {
	return To[T](tree, DefaultSchemaContext())
}

func FromDefault[T any](obj T) calm.ResultT[ITree] {
	return From(obj, DefaultSchemaContext())
}

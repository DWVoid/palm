package palm

import (
	"fmt"
	"reflect"

	"github.com/DWVoid/calm"
)

type _ArraySchema struct {
	Elem  IMarshalSchema
	Count int64
}

func (ths *_ArraySchema) Zero(value reflect.Value) {
	for i := 0; i < int(ths.Count); i++ {
		ths.Elem.Zero(value.Index(i))
	}
}

func (ths *_ArraySchema) Write(value reflect.Value, tree ITree) IMarshalErr {
	if !IsArr(tree) {
		return NewMarshalErrType(TypeArr, tree.Type())
	}
	arr := tree.AsArr().Get()
	if int64(len(arr)) != ths.Count {
		return NewMarshalErr(fmt.Sprintf("array length: expect %d got %d", ths.Count, len(arr)))
	}
	for i, node := range arr {
		if err := ths.Elem.Write(value.Index(i), node); err != nil {
			return WrapMarshalErr(err, fmt.Sprintf("[%d]", i))
		}
	}
	return nil
}

func (ths *_ArraySchema) Read(value reflect.Value) ITree {
	arr := make([]ITree, value.Len())
	for i := 0; i < value.Len(); i++ {
		arr[i] = ths.Elem.Read(value.Index(i))
	}
	return NewArr(arr)
}

func _CreateArraySchema(schema reflect.Type, ctx IMarshalSchemaResolutionContext) calm.ResultT[IMarshalSchema] {
	return calm.RunT(func() IMarshalSchema {
		return &_ArraySchema{Elem: ctx.Resolve(schema.Elem()), Count: int64(schema.Len())}
	})
}

type _SliceSchema struct {
	Elem IMarshalSchema
	Type reflect.Type
}

func (ths *_SliceSchema) Zero(value reflect.Value) {
	value.Set(reflect.MakeSlice(ths.Type, 0, 0))
}

func (ths *_SliceSchema) Write(value reflect.Value, tree ITree) IMarshalErr {
	if !IsArr(tree) {
		return NewMarshalErrType(TypeArr, tree.Type())
	}
	arr := tree.AsArr().Get()
	value.Set(reflect.MakeSlice(ths.Type, len(arr), len(arr)))
	for i, node := range arr {
		if err := ths.Elem.Write(value.Index(i), node); err != nil {
			return WrapMarshalErr(err, fmt.Sprintf("[%d]", i))
		}
	}
	return nil
}

func (ths *_SliceSchema) Read(value reflect.Value) ITree {
	arr := make([]ITree, value.Len())
	for i := 0; i < value.Len(); i++ {
		arr[i] = ths.Elem.Read(value.Index(i))
	}
	return NewArr(arr)
}

func _CreateSliceSchema(schema reflect.Type, ctx IMarshalSchemaResolutionContext) calm.ResultT[IMarshalSchema] {
	return calm.RunT(func() IMarshalSchema {
		return &_SliceSchema{Elem: ctx.Resolve(schema.Elem()), Type: schema}
	})
}

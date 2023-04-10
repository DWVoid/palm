package palm

import (
	"reflect"

	"github.com/DWVoid/calm"
)

type _PointerSchema struct {
	Elem IMarshalSchema
	Type reflect.Type
}

func (ths *_PointerSchema) Zero(value reflect.Value) {
	value.SetZero()
}

func (ths *_PointerSchema) Write(value reflect.Value, tree ITree) IMarshalErr {
	if IsNull(tree) {
		value.SetZero()
		return nil
	}
	value.Set(reflect.New(ths.Type))
	return ths.Elem.Write(value.Elem(), tree)
}

func (ths *_PointerSchema) Read(value reflect.Value) ITree {
	if value.IsZero() {
		return NewNull()
	}
	return ths.Elem.Read(value.Elem())
}

func _CreatePointerSchema(schema reflect.Type, ctx IMarshalSchemaResolutionContext) calm.ResultT[IMarshalSchema] {
	return calm.RunT(func() IMarshalSchema {
		return &_PointerSchema{Elem: ctx.Resolve(schema.Elem()), Type: schema.Elem()}
	})
}

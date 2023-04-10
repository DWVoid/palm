package palm

import (
	"reflect"

	"github.com/DWVoid/calm"
)

type _RoStrTree string

func (o _RoStrTree) Type() int32 { return TypeStr }

func (o _RoStrTree) AsBool() (x calm.ResultT[bool])            { return }
func (o _RoStrTree) AsNum() (x calm.ResultT[float64])          { return }
func (o _RoStrTree) AsMap() (x calm.ResultT[map[string]ITree]) { return }
func (o _RoStrTree) AsArr() (x calm.ResultT[[]ITree])          { return }
func (o _RoStrTree) AsStr() calm.ResultT[string]               { return calm.ValResultT(string(o)) }

func (o _RoStrTree) SetBool(bool)            {}
func (o _RoStrTree) SetNum(float64)          {}
func (o _RoStrTree) SetStr(string)           {}
func (o _RoStrTree) SetMap(map[string]ITree) {}
func (o _RoStrTree) SetArr([]ITree)          {}

type _MapSchema struct {
	Key, Value       IMarshalSchema
	TM, TKey, TValue reflect.Type
}

func (ths *_MapSchema) Zero(value reflect.Value) {
	value.SetZero()
}

func (ths *_MapSchema) Write(value reflect.Value, tree ITree) IMarshalErr {
	if !IsMap(tree) {
		return NewMarshalErrType(TypeMap, tree.Type())
	}
	dict := tree.AsMap().Get()
	value.Set(reflect.MakeMap(ths.TM))
	for k, v := range dict {
		rK := reflect.New(ths.TKey).Elem()
		rV := reflect.New(ths.TValue).Elem()
		if err := ths.Key.Write(rK, _RoStrTree(k)); err != nil {
			return WrapMarshalErr(err, k)
		}
		if err := ths.Value.Write(rV, v); err != nil {
			return WrapMarshalErr(err, k)
		}
		value.SetMapIndex(rK, rV)
	}
	return nil
}

func (ths *_MapSchema) Read(value reflect.Value) ITree {
	dict := make(map[string]ITree)
	iter := value.MapRange()
	for iter.Next() {
		dict[ths.Key.Read(iter.Key()).AsStr().Get()] = ths.Value.Read(iter.Value())
	}
	return NewMap(dict)
}

func _CreateMapSchema(schema reflect.Type, ctx IMarshalSchemaResolutionContext) calm.ResultT[IMarshalSchema] {
	return calm.RunT(func() IMarshalSchema {
		return &_MapSchema{
			Key:   ctx.Resolve(schema.Key()),
			Value: ctx.Resolve(schema.Elem()),
			TM:    schema, TKey: schema.Key(), TValue: schema.Elem(),
		}
	})
}

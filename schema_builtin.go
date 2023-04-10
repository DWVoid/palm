package palm

import (
	"reflect"
)

type _Int interface {
	int | int8 | int16 | int32 | int64
}

type _UInt interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type _Real interface {
	float32 | float64
}

type _IntSchema[T _Int] struct{}

func (ths *_IntSchema[T]) Zero(value reflect.Value) {
	value.SetZero()
}

func (ths *_IntSchema[T]) Write(value reflect.Value, tree ITree) IMarshalErr {
	if !IsNum(tree) {
		return NewMarshalErrType(TypeNum, tree.Type())
	}
	value.SetInt(int64(tree.AsNum().Get()))
	return nil
}

func (ths *_IntSchema[T]) Read(value reflect.Value) ITree {
	return NewNum(float64(value.Int()))
}

type _UIntSchema[T _UInt] struct{}

func (ths *_UIntSchema[T]) Zero(value reflect.Value) {
	value.SetZero()
}

func (ths *_UIntSchema[T]) Read(value reflect.Value) ITree {
	return NewNum(float64(value.Uint()))
}

func (ths *_UIntSchema[T]) Write(value reflect.Value, tree ITree) IMarshalErr {
	if !IsNum(tree) {
		return NewMarshalErrType(TypeNum, tree.Type())
	}
	value.SetUint(uint64(tree.AsNum().Get()))
	return nil
}

type _RealSchema[T _Real] struct{}

func (ths *_RealSchema[T]) Zero(value reflect.Value) {
	value.SetZero()
}

func (ths *_RealSchema[T]) Write(value reflect.Value, tree ITree) IMarshalErr {
	if !IsNum(tree) {
		return NewMarshalErrType(TypeNum, tree.Type())
	}
	value.SetFloat(tree.AsNum().Get())
	return nil
}

func (ths *_RealSchema[T]) Read(value reflect.Value) ITree {
	return NewNum(value.Float())
}

type _BoolSchema struct{}

func (ths *_BoolSchema) Zero(value reflect.Value) {
	value.SetZero()
}

func (ths *_BoolSchema) Write(value reflect.Value, tree ITree) IMarshalErr {
	if !IsBool(tree) {
		return NewMarshalErrType(TypeBool, tree.Type())
	}
	value.SetBool(tree.AsBool().Get())
	return nil
}

func (ths *_BoolSchema) Read(value reflect.Value) ITree {
	return NewBool(value.Bool())
}

type _StringSchema struct{}

func (ths *_StringSchema) Zero(value reflect.Value) {
	value.SetZero()
}

func (ths *_StringSchema) Write(value reflect.Value, tree ITree) IMarshalErr {
	if !IsStr(tree) {
		return NewMarshalErrType(TypeStr, tree.Type())
	}
	value.SetString(tree.AsStr().Get())
	return nil
}

func (ths *_StringSchema) Read(value reflect.Value) ITree {
	return NewStr(value.String())
}

type _ITreeSchema struct{}

func (ths *_ITreeSchema) Zero(value reflect.Value) {
	value.Set(reflect.ValueOf(NewNull()))
}

func (ths *_ITreeSchema) Write(value reflect.Value, tree ITree) IMarshalErr {
	value.Set(reflect.ValueOf(tree))
	return nil
}

func (ths *_ITreeSchema) Read(value reflect.Value) ITree {
	return value.Interface().(ITree)
}

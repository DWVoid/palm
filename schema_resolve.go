package palm

import (
	"reflect"

	"github.com/DWVoid/calm"
	"github.com/DWVoid/cart/conc"
)

var (
	_Default _MarshalSchemaContext
)

type _MarshalSchemaResolution struct {
	Schema    *calm.ResultT[IMarshalSchema]
	MinDepend IMarshalSchemaContext
}

type IMarshalSchemaContext interface {
	With() IMarshalSchemaResolutionContext
	Supply(t reflect.Type, schema IMarshalSchema)
	set(t reflect.Type, v calm.ResultT[IMarshalSchema])
	resolve(t reflect.Type) _MarshalSchemaResolution
	depth() int32
}

type IMarshalSchemaResolutionContext interface {
	Resolve(t reflect.Type) IMarshalSchema
}

func DefaultSchemaContext() IMarshalSchemaContext {
	return &_Default
}

func MakeSchemaContext(parent IMarshalSchemaContext) IMarshalSchemaContext {
	return &_MarshalSchemaContext{
		Parent: parent,
		Depth:  parent.depth() + 1,
		Table:  conc.Map[reflect.Type, calm.ResultT[IMarshalSchema]]{},
	}
}

type _MarshalSchemaContext struct {
	Depth  int32
	Parent IMarshalSchemaContext
	Table  conc.Map[reflect.Type, calm.ResultT[IMarshalSchema]]
}

func (ths *_MarshalSchemaContext) With() IMarshalSchemaResolutionContext {
	return &_MarshalSchemaResolutionContext{context: ths}
}

func (ths *_MarshalSchemaContext) Supply(t reflect.Type, schema IMarshalSchema) {
	ths.set(t, calm.ValResultT(schema))
}

func (ths *_MarshalSchemaContext) set(t reflect.Type, v calm.ResultT[IMarshalSchema]) {
	// schema resolution on the same type should always produce the same result, so just overwrite is fine
	ths.Table.Store(t, v)
}

func (ths *_MarshalSchemaContext) resolve(t reflect.Type) _MarshalSchemaResolution {
	load, has := ths.Table.Load(t)
	if has {
		return _MarshalSchemaResolution{Schema: &load, MinDepend: ths}
	}
	if ths.Parent != nil {
		return ths.Parent.resolve(t)
	}
	return _MarshalSchemaResolution{}
}

func (ths *_MarshalSchemaContext) depth() int32 {
	if ths.Parent != nil {
		return ths.Parent.depth() + 1
	}
	return 0
}

type _MarshalSchemaResolutionContext struct {
	context   IMarshalSchemaContext
	minDepend IMarshalSchemaContext
	stack     []IMarshalSchemaContext
}

func (ths *_MarshalSchemaResolutionContext) Push() {
	ths.stack = append(ths.stack, ths.minDepend)
}

func (ths *_MarshalSchemaResolutionContext) Pop(t reflect.Type, res calm.ResultT[IMarshalSchema]) IMarshalSchema {
	// add result to corresponding schema context
	ths.minDepend.set(t, res)
	// pop and merge the stack
	ths.UpdateMinDepend(ths.stack[len(ths.stack)-1])
	ths.stack = ths.stack[:len(ths.stack)-1]
	// return res
	return res.Get()
}

func (ths *_MarshalSchemaResolutionContext) UpdateMinDepend(min IMarshalSchemaContext) {
	if min == nil {
		return
	}
	if ths.minDepend == min {
		return
	}
	if ths.minDepend != nil {
		if ths.minDepend.depth() <= min.depth() {
			return
		}
	}
	ths.minDepend = min
}

func (ths *_MarshalSchemaResolutionContext) _CreateSchema(t reflect.Type) calm.ResultT[IMarshalSchema] {
	// supported built-in types are predefined during init
	switch t.Kind() {
	case reflect.Map:
		return _CreateMapSchema(t, ths)
	case reflect.Pointer:
		return _CreatePointerSchema(t, ths)
	case reflect.Array:
		return _CreateArraySchema(t, ths)
	case reflect.Slice:
		return _CreateSliceSchema(t, ths)
	case reflect.Struct:
		return _CreateStructSchema(t, ths)
	}
	return calm.ErrResultT[IMarshalSchema](calm.ErrClean(calm.EInternal, "unable to construct schema"))
}

func (ths *_MarshalSchemaResolutionContext) Resolve(t reflect.Type) (result IMarshalSchema) {
	result = nil
	load := ths.context.resolve(t)
	if load.Schema != nil {
		if load.MinDepend == ths.context {
			ths.UpdateMinDepend(load.MinDepend)
			return load.Schema.Get() // same context, use cached result
		}
		// parent context, if err do retry
		result = load.Schema.Fold(func(calm.Error) IMarshalSchema { return nil })
		if result != nil {
			ths.UpdateMinDepend(load.MinDepend)
			return
		}
	}
	ths.Push()
	return ths.Pop(t, ths._CreateSchema(t))
}

func init() {
	_Default = _MarshalSchemaContext{Parent: nil, Depth: 0, Table: conc.Map[reflect.Type, calm.ResultT[IMarshalSchema]]{}}
	_Default.Supply(_TypeOf[int](), &_IntSchema[int]{})
	_Default.Supply(_TypeOf[int8](), &_IntSchema[int8]{})
	_Default.Supply(_TypeOf[int16](), &_IntSchema[int16]{})
	_Default.Supply(_TypeOf[int32](), &_IntSchema[int32]{})
	_Default.Supply(_TypeOf[int64](), &_IntSchema[int64]{})
	_Default.Supply(_TypeOf[uint](), &_UIntSchema[uint]{})
	_Default.Supply(_TypeOf[uint8](), &_UIntSchema[uint8]{})
	_Default.Supply(_TypeOf[uint16](), &_UIntSchema[uint16]{})
	_Default.Supply(_TypeOf[uint32](), &_UIntSchema[uint32]{})
	_Default.Supply(_TypeOf[uint64](), &_UIntSchema[uint64]{})
	_Default.Supply(_TypeOf[float32](), &_RealSchema[float32]{})
	_Default.Supply(_TypeOf[float64](), &_RealSchema[float64]{})
	_Default.Supply(_TypeOf[bool](), &_BoolSchema{})
	_Default.Supply(_TypeOf[string](), &_StringSchema{})
	_Default.Supply(_TypeOf[ITree](), &_ITreeSchema{})
}

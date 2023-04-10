package palm

import (
	"fmt"

	"github.com/DWVoid/calm"
)

const (
	TypeNull = 0
	TypeBool = 1
	TypeNum  = 2
	TypeStr  = 3
	TypeMap  = 4
	TypeArr  = 5
)

type ITree interface {
	Type() int32
	AsBool() calm.ResultT[bool]
	AsNum() calm.ResultT[float64]
	AsStr() calm.ResultT[string]
	AsMap() calm.ResultT[map[string]ITree]
	AsArr() calm.ResultT[[]ITree]
	SetBool(bool)
	SetNum(float64)
	SetStr(string)
	SetMap(map[string]ITree)
	SetArr([]ITree)
}

func IsNull(tree ITree) bool { return tree.Type() == TypeNull }
func IsBool(tree ITree) bool { return tree.Type() == TypeBool }
func IsNum(tree ITree) bool  { return tree.Type() == TypeNum }
func IsStr(tree ITree) bool  { return tree.Type() == TypeStr }
func IsMap(tree ITree) bool  { return tree.Type() == TypeMap }
func IsArr(tree ITree) bool  { return tree.Type() == TypeArr }

type _TreeNull struct {
}

func (n *_TreeNull) Type() int32 { return TypeNull }

func _ErrT[T any]() calm.Error {
	return calm.ErrClean(
		calm.EInternal,
		fmt.Sprintf("ITree ECode Mismatch：want %s", _TypeOf[T]()),
	)
}

func _TypeErrT[T any]() calm.ResultT[T] { return calm.ErrResultT[T](_ErrT[T]()) }
func _TypeErr[T any]() calm.Result      { return calm.ErrResult(_ErrT[T]()) }

func (n *_TreeNull) AsBool() calm.ResultT[bool]            { return _TypeErrT[bool]() }
func (n *_TreeNull) AsNum() calm.ResultT[float64]          { return _TypeErrT[float64]() }
func (n *_TreeNull) AsStr() calm.ResultT[string]           { return _TypeErrT[string]() }
func (n *_TreeNull) AsMap() calm.ResultT[map[string]ITree] { return _TypeErrT[map[string]ITree]() }
func (n *_TreeNull) AsArr() calm.ResultT[[]ITree]          { return _TypeErrT[[]ITree]() }

func (n *_TreeNull) SetBool(bool)            { _TypeErr[bool]() }
func (n *_TreeNull) SetNum(float64)          { _TypeErr[float64]() }
func (n *_TreeNull) SetStr(string)           { _TypeErr[string]() }
func (n *_TreeNull) SetMap(map[string]ITree) { _TypeErr[map[string]ITree]() }
func (n *_TreeNull) SetArr([]ITree)          { _TypeErr[[]ITree]() }

type (
	_TreeBool struct {
		_TreeNull
		v bool
	}
	_TreeNum struct {
		_TreeNull
		v float64
	}
	_TreeStr struct {
		_TreeNull
		v string
	}
	_TreeMap struct {
		_TreeNull
		v map[string]ITree
	}
	_TreeArr struct {
		_TreeNull
		v []ITree
	}
)

func (n *_TreeBool) Type() int32                { return TypeBool }
func (n *_TreeBool) AsBool() calm.ResultT[bool] { return calm.ValResultT(n.v) }
func (n *_TreeBool) SetBool(v bool)             { n.v = v }

func (n *_TreeNum) Type() int32                  { return TypeNum }
func (n *_TreeNum) AsNum() calm.ResultT[float64] { return calm.ValResultT(n.v) }
func (n *_TreeNum) SetNum(v float64)             { n.v = v }

func (n *_TreeStr) Type() int32                 { return TypeStr }
func (n *_TreeStr) AsStr() calm.ResultT[string] { return calm.ValResultT(n.v) }
func (n *_TreeStr) SetStr(v string)             { n.v = v }

func (n *_TreeMap) Type() int32                           { return TypeMap }
func (n *_TreeMap) AsMap() calm.ResultT[map[string]ITree] { return calm.ValResultT(n.v) }
func (n *_TreeMap) SetMap(v map[string]ITree)             { n.v = v }

func (n *_TreeArr) Type() int32                  { return TypeArr }
func (n *_TreeArr) AsArr() calm.ResultT[[]ITree] { return calm.ValResultT(n.v) }
func (n *_TreeArr) SetArr(v []ITree)             { n.v = v }

func NewNull() ITree                  { return &_TreeNull{} }
func NewBool(v bool) ITree            { return &_TreeBool{v: v} }
func NewNum(v float64) ITree          { return &_TreeNum{v: v} }
func NewStr(v string) ITree           { return &_TreeStr{v: v} }
func NewMap(v map[string]ITree) ITree { return &_TreeMap{v: v} }
func NewArr(v []ITree) ITree          { return &_TreeArr{v: v} }

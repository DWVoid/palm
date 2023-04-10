package palm

import "github.com/DWVoid/cart/algo"

func FromAny(v any) ITree {
	switch r := v.(type) {
	case bool:
		return NewBool(r)
	case int:
		return NewNum(float64(r))
	case int8:
		return NewNum(float64(r))
	case int16:
		return NewNum(float64(r))
	case int32:
		return NewNum(float64(r))
	case int64:
		return NewNum(float64(r))
	case uint:
		return NewNum(float64(r))
	case uint8:
		return NewNum(float64(r))
	case uint16:
		return NewNum(float64(r))
	case uint32:
		return NewNum(float64(r))
	case uint64:
		return NewNum(float64(r))
	case float32:
		return NewNum(float64(r))
	case float64:
		return NewNum(r)
	case string:
		return NewStr(r)
	case []interface{}:
		return NewArr(algo.Map(r, FromAny))
	case map[string]interface{}:
		return NewMap(algo.MapMap(r, func(a string, b any) (string, ITree) { return a, FromAny(b) }))
	}
	return NewNull()
}

func ToAny(v ITree) any {
	switch v.Type() {
	case TypeBool:
		return v.AsBool().Get()
	case TypeNum:
		return v.AsNum().Get()
	case TypeStr:
		return v.AsStr().Get()
	case TypeArr:
		return algo.Map(v.AsArr().Get(), ToAny)
	case TypeMap:
		return algo.MapMap(v.AsMap().Get(), func(a string, b ITree) (string, any) { return a, ToAny(b) })
	}
	return nil
}

package palm

import (
	"net/url"
	"strconv"

	"github.com/DWVoid/calm"
)

func FromUrlValues(v url.Values) ITree {
	m := make(map[string]ITree)
	for k, vs := range v {
		f := vs[0]
		if f != "" {
			m[k] = NewStr(f)
		}
	}
	return NewMap(m)
}

func ToUrlValues(v ITree) calm.ResultT[url.Values] {
	return calm.RunT(func() (result url.Values) {
		calm.Assert(IsMap(v), "url encoding only support maps").Get()
		result = url.Values{}
		for k, s := range v.AsMap().Get() {
			switch s.Type() {
			case TypeBool:
				result.Set(k, strconv.FormatBool(s.AsBool().Get()))
			case TypeNum:
				result.Set(k, strconv.FormatFloat(s.AsNum().Get(), 'f', -1, 64))
			case TypeStr:
				result.Set(k, s.AsStr().Get())
			default:
				calm.ThrowClean(calm.EInternal, "url encoding does not support nested types")
			}
		}
		return
	})
}

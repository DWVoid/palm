package palm

import (
	"reflect"
	"strings"

	"github.com/DWVoid/calm"
)

type _StructSchema struct {
	Fields []_StructFieldOpt
}

type _StructFieldOpt struct {
	Id      int            // field index in the struct
	Name    string         // serialized name of this field
	InOpt   bool           // if value in tree is null, set value in struct to zero
	OutOpt  bool           // if value in struct is zero, set value in tree to zero
	Explode bool           // explode the containing fields to the current scope. only applicable to anno structs
	Schema  IMarshalSchema // mapping schema of element
}

func (ths *_StructSchema) Zero(value reflect.Value) {
	for i, field := range ths.Fields {
		field.Schema.Zero(value.Field(i))
	}
}

func (ths *_StructSchema) Write(value reflect.Value, tree ITree) IMarshalErr {
	if !IsMap(tree) {
		return NewMarshalErrType(TypeMap, tree.Type())
	}
	dict := tree.AsMap().Get()
	for _, field := range ths.Fields {
		if !field.Explode {
			item, has := dict[field.Name]
			if !has {
				if !field.InOpt {
					return NewMarshalErr("missing required field: " + field.Name)
				}
				field.Schema.Zero(value.Field(field.Id))
				continue
			}
			if err := field.Schema.Write(value.Field(field.Id), item); err != nil {
				return WrapMarshalErr(err, field.Name)
			}
		} else {
			if err := field.Schema.Write(value.Field(field.Id), tree); err != nil {
				return WrapMarshalErr(err, field.Name)
			}
		}
	}
	return nil
}

func (ths *_StructSchema) ReadTo(value reflect.Value, dict map[string]ITree) {
	for _, field := range ths.Fields {
		fValue := value.Field(field.Id)
		if !field.Explode {
			if field.OutOpt && fValue.IsZero() {
				continue
			}
			fTree := field.Schema.Read(fValue)
			dict[field.Name] = fTree
		} else {
			field.Schema.(*_StructSchema).ReadTo(fValue, dict)
		}
	}
}

func (ths *_StructSchema) Read(value reflect.Value) ITree {
	dict := make(map[string]ITree)
	ths.ReadTo(value, dict)
	return NewMap(dict)
}

func (ths *_StructSchema) RecursiveCheckCollision(set map[string]bool) string {
	for _, field := range ths.Fields {
		if set[field.Name] {
			return field.Name
		} else {
			if field.Explode {
				field.Schema.(*_StructSchema).RecursiveCheckCollision(set)
			} else {
				set[field.Name] = true
			}
		}
	}
	return ""
}

func _CreateStructSchema(schema reflect.Type, ctx IMarshalSchemaResolutionContext) calm.ResultT[IMarshalSchema] {
	return calm.RunT(func() IMarshalSchema {
		result := _StructSchema{}
		for i := 0; i < schema.NumField(); i++ {
			rField := schema.Field(i)
			rFTags := strings.Split(rField.Tag.Get("palm"), ",")
			field := _StructFieldOpt{Id: i, Name: rFTags[0]}
			if field.Name == "" {
				field.Name = rField.Name
			}
			if rField.Anonymous && rField.Type.Kind() == reflect.Struct {
				field.Explode = true
			}
			skip := !rField.IsExported()
			for _, tag := range rFTags[1:] {
				switch strings.TrimSpace(tag) {
				case "skip":
					skip = true
				case "opt":
					field.InOpt = true
					field.OutOpt = true
				case "in_opt":
					field.InOpt = true
				case "out_opt":
					field.OutOpt = true
				case "no_explode":
					field.Explode = false
				}
			}
			if !skip {
				field.Schema = ctx.Resolve(rField.Type)
				if field.Explode {
					if _, ok := field.Schema.(*_StructSchema); !ok {
						calm.ThrowClean(calm.EInternal, "exploding a custom struct marshaller")
					}
				}
				result.Fields = append(result.Fields, field)
			}
		}
		// now as all fields are resolved, we recursively scan for collisions of exploded fields
		if check := result.RecursiveCheckCollision(make(map[string]bool)); check != "" {
			calm.ThrowClean(calm.EInternal, "collision of exploded field: "+check)
		}
		return &result
	})
}

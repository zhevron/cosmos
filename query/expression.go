package query

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Expression interface {
	String() string
}

type equal struct {
	Field string
	Value interface{}
}

func Equal(field string, value interface{}) Expression {
	return equal{Field: field, Value: value}
}

func (e equal) String() string {
	return e.Field + " = " + valueToString(e.Value)
}

type notEqual struct {
	Field string
	Value interface{}
}

func NotEqual(field string, value interface{}) Expression {
	return notEqual{Field: field, Value: value}
}

func (e notEqual) String() string {
	return e.Field + " != " + valueToString(e.Value)
}

type less struct {
	Field string
	Value interface{}
}

func Less(field string, value interface{}) Expression {
	return less{Field: field, Value: value}
}

func (e less) String() string {
	return e.Field + " < " + valueToString(e.Value)
}

type lessOrEqual struct {
	Field string
	Value interface{}
}

func LessOrEqual(field string, value interface{}) Expression {
	return lessOrEqual{Field: field, Value: value}
}

func (e lessOrEqual) String() string {
	return e.Field + " <= " + valueToString(e.Value)
}

type greater struct {
	Field string
	Value interface{}
}

func Greater(field string, value interface{}) Expression {
	return greater{Field: field, Value: value}
}

func (e greater) String() string {
	return e.Field + " > " + valueToString(e.Value)
}

type greaterOrEqual struct {
	Field string
	Value interface{}
}

func GreaterOrEqual(field string, value interface{}) Expression {
	return greaterOrEqual{Field: field, Value: value}
}

func (e greaterOrEqual) String() string {
	return e.Field + " >= " + valueToString(e.Value)
}

type in struct {
	Field  string
	Values []interface{}
}

func In(field string, values interface{}) Expression {
	rv := reflect.ValueOf(values)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		panic("non-array/slice value passed to query.In")
	}

	valuesLen := rv.Len()
	interfaceValues := make([]interface{}, valuesLen)
	for i := 0; i < valuesLen; i++ {
		interfaceValues[i] = rv.Index(i).Interface()
	}

	return in{Field: field, Values: interfaceValues}
}

func (e in) String() string {
	values := make([]string, len(e.Values))
	for i, v := range e.Values {
		values[i] = valueToString(v)
	}

	return e.Field + " IN (" + strings.Join(values, ", ") + ")"
}

type And []Expression

func (e And) String() string {
	exprs := make([]string, len(e))
	for i, ex := range e {
		exprs[i] = ex.String()
	}

	return "(" + strings.Join(exprs, " AND ") + ")"
}

type Or []Expression

func (e Or) String() string {
	exprs := make([]string, len(e))
	for i, ex := range e {
		exprs[i] = ex.String()
	}

	return "(" + strings.Join(exprs, " OR ") + ")"
}

type isNull struct {
	Field string
}

func IsNull(field string) Expression {
	return isNull{Field: field}
}

func (e isNull) String() string {
	return "IS_NULL(" + e.Field + ")"
}

type isNotNull struct {
	Field string
}

func IsNotNull(field string) Expression {
	return isNotNull{Field: field}
}

func (e isNotNull) String() string {
	return "IS_NULL(" + e.Field + ") = false"
}

type arrayContainsOption func(e *arrayContains)

func ContainsPartial(e *arrayContains) {
	e.Partial = true
}

type arrayContains struct {
	Field   string
	Value   interface{}
	Partial bool
}

func ArrayContains(field string, value interface{}, opts ...arrayContainsOption) Expression {
	e := &arrayContains{
		Field:   field,
		Value:   value,
		Partial: false,
	}

	for _, o := range opts {
		o(e)
	}

	return *e
}

func (e arrayContains) String() string {
	return "ARRAY_CONTAINS(" + e.Field + ", " + valueToString(e.Value) + ", " + valueToString(e.Partial) + ")"
}

type arrayNotContains struct {
	arrayContains
}

func ArrayNotContains(field string, value interface{}, opts ...arrayContainsOption) Expression {
	return arrayNotContains{
		arrayContains: ArrayContains(field, value, opts...).(arrayContains),
	}
}

func (e arrayNotContains) String() string {
	return e.arrayContains.String() + " = false"
}

func valueToString(value interface{}) string {
	switch v := value.(type) {
	case int:
		return strconv.FormatInt(int64(v), 10)

	case int8:
		return strconv.FormatInt(int64(v), 10)

	case int16:
		return strconv.FormatInt(int64(v), 10)

	case int32:
		return strconv.FormatInt(int64(v), 10)

	case int64:
		return strconv.FormatInt(v, 10)

	case uint:
		return strconv.FormatUint(uint64(v), 10)

	case uint8:
		return strconv.FormatUint(uint64(v), 10)

	case uint16:
		return strconv.FormatUint(uint64(v), 10)

	case uint32:
		return strconv.FormatUint(uint64(v), 10)

	case uint64:
		return strconv.FormatUint(v, 10)

	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)

	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)

	case bool:
		return strconv.FormatBool(v)

	case string:
		if strings.HasPrefix(v, "@") {
			return v
		}
		return "'" + v + "'"

	default:
		rv := reflect.ValueOf(v)
		if !rv.IsValid() || rv.IsNil() {
			return "null"
		}

		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}

		switch rv.Kind() {
		case reflect.Array, reflect.Slice:
			arr := make([]string, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				arr[i] = valueToString(rv.Index(i))
			}
			return "[" + strings.Join(arr, ",") + "]"

		case reflect.Map:
			var arr []string
			iter := rv.MapRange()
			for iter.Next() {
				arr = append(arr, valueToString(iter.Key().Interface())+": "+valueToString(iter.Value().Interface()))
			}
			return "{" + strings.Join(arr, ",") + "}"

		default:
			if b, err := json.Marshal(rv.Interface()); err == nil {
				return string(b)
			}
		}

		return fmt.Sprintf("'%v'", v)
	}
}

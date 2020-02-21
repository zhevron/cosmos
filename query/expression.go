package query

import (
	"fmt"
	"strconv"
	"strings"
)

type Expression interface {
	String() string
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
		return fmt.Sprintf("'%v'", v)
	}
}

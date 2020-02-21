package query

import (
	"fmt"
	"strconv"
	"strings"
)

type Expression interface {
	String() string
}

type IsNull struct {
	Field string
}

func (e IsNull) String() string {
	return "IS_NULL(" + e.Field + ")"
}

type IsNotNull struct {
	Field string
}

func (e IsNotNull) String() string {
	return "IS_NULL(" + e.Field + ") = false"
}

type Equal struct {
	Field string
	Value interface{}
}

func (e Equal) String() string {
	return e.Field + "=" + valueToString(e.Value)
}

type NotEqual struct {
	Field string
	Value interface{}
}

func (e NotEqual) String() string {
	return e.Field + "!=" + valueToString(e.Value)
}

type Less struct {
	Field string
	Value interface{}
}

func (e Less) String() string {
	return e.Field + "<" + valueToString(e.Value)
}

type LessOrEqual struct {
	Field string
	Value interface{}
}

func (e LessOrEqual) String() string {
	return e.Field + "<=" + valueToString(e.Value)
}

type Greater struct {
	Field string
	Value interface{}
}

func (e Greater) String() string {
	return e.Field + ">" + valueToString(e.Value)
}

type GreaterOrEqual struct {
	Field string
	Value interface{}
}

func (e GreaterOrEqual) String() string {
	return e.Field + ">=" + valueToString(e.Value)
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
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(v.(int64), 10)

	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(v.(uint64), 10)

	case float32, float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)

	case bool:
		return strconv.FormatBool(v)

	case string:
		return "'" + v + "'"

	default:
		return fmt.Sprintf("'%v'", v)
	}
}

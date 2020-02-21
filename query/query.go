package query

import (
	"strings"
)

type Order string

const (
	Ascending  Order = "ASC"
	Descending Order = "DESC"
)

type Query struct {
	fields []string
	from   string
	where  Expression
	order  *struct {
		field     string
		direction Order
	}
}

func Select(fields ...string) Query {
	if len(fields) == 0 {
		fields = []string{"*"}
	}

	return Query{
		fields: fields,
		from:   "c",
	}
}

func (q Query) Where(expr Expression) Query {
	q.where = expr

	return q
}

func (q Query) OrderBy(field string, direction Order) Query {
	q.order = &struct {
		field     string
		direction Order
	}{
		field:     field,
		direction: direction,
	}

	return q
}

func (q Query) String() string {
	query := "SELECT " + strings.Join(q.fields, ",") + " FROM " + q.from // nolint:gosec

	if q.where != nil {
		query += " WHERE " + q.where.String() // nolint:gosec
	}

	if q.order != nil {
		query += " ORDER BY " + q.order.field + " " + string(q.order.direction) // nolint:gosec
	}

	return query
}

type QueryParameter struct {
	Name  string
	Value interface{}
}

func (p QueryParameter) ValueAsString() string {
	return valueToString(p.Value)
}

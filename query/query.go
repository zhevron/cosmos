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
	joins  []string
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

func (q Query) Join(alias string, source string) Query {
	join := "JOIN " + alias + " IN " + source
	if alias == "" {
		join = "JOIN " + source
	}

	q.joins = append(q.joins, join)

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

	if len(q.joins) > 0 {
		query += " " + strings.Join(q.joins, " ")
	}

	if q.where != nil {
		query += " WHERE " + q.where.String() // nolint:gosec
	}

	if q.order != nil {
		query += " ORDER BY " + q.order.field + " " + string(q.order.direction) // nolint:gosec
	}

	return query
}

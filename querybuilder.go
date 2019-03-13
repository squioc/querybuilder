package builder

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

/**
*
* QueryBuilder - Helps to generate sql query
*
**/

var ErrEmptyBaseQuery = errors.New("The base query cannot be empty")

// criterion represents a key-value pair for the where clause
type criterion struct {
	name  string
	value interface{}
}

// QueryBuilder represents a struct that help to generate sql query
type QueryBuilder struct {
	// baseQuery is the starting query from which the final query will be built
	baseQuery string

	// criteria is a list of filter to add to the query
	criteria []criterion
}

// QueryBuilderOption is a type to define function that manipulate the `QueryBuiler`
type QueryBuilderOption func(*QueryBuilder)

// NewQueryBuilder creates a new query builder
func NewQueryBuilder(baseQuery string, opts ...QueryBuilderOption) (*QueryBuilder, error) {

	// checks the base query
	if len(baseQuery) == 0 {
		return nil, ErrEmptyBaseQuery
	}

	// initializes the query builder
	qb := &QueryBuilder{
		baseQuery: baseQuery,
	}

	// applies options on the query builder
	for _, opt := range opts {
		opt(qb)
	}

	return qb, nil
}

// WithKVCriteria allows to initialize the `QueryBuilder` with a list of criterion
func WithKVCriteria(kvs map[string]interface{}) QueryBuilderOption {
	return func(qb *QueryBuilder) {
		for name, value := range kvs {
			qb.AppendCriterion(name, value)
		}
	}
}

// AppendCriterion adds a new `criterion` in the `QueryBuilder`
func (qb *QueryBuilder) AppendCriterion(name string, value interface{}) {
	kv := criterion{
		name:  name,
		value: value,
	}
	qb.criteria = append(qb.criteria, kv)
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
	queryString := bytes.NewBufferString(qb.baseQuery)

	var filters []string
	var values []interface{}
	for index, criterion := range qb.criteria {
		filters = append(filters, fmt.Sprintf("%s = $%d", criterion.name, index+1))
		values = append(values, criterion.value)
	}

	if len(filters) > 0 {
		whereClause := strings.Join(filters, " and ")

		queryString.WriteString(" where ")
		queryString.WriteString(whereClause)
	}

	return queryString.String(), values
}

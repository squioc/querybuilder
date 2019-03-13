package builder_test

import (
	"github.com/squioc/querybuilder"
	"testing"
)

type queryBuilderChecker func(testing.TB, *builder.QueryBuilder, error)

func CheckQueryBuilderNoError() queryBuilderChecker {
	return func(t testing.TB, qb *builder.QueryBuilder, err error) {
		if err != nil {
			t.Errorf("Unexpected error. Got: %s", err)
		} else {
			if qb == nil {
				t.Errorf("Expected QueryBuilder. Got nil")
			}
		}
	}
}

func CheckQueryBuilderError(expectedErr error) queryBuilderChecker {
	return func(t testing.TB, qb *builder.QueryBuilder, err error) {
		if err == nil {
			t.Errorf("Expected error. Got nil")
		}
		if err != expectedErr {
			t.Errorf("Unexpected error. Expected: %s. Got: %s", expectedErr, err)
		}
	}
}

type queryBuilderResultChecker func(testing.TB, string, []interface{})

func CheckQueryBuilderResult(expectedQuery string, expectedValues []interface{}) queryBuilderResultChecker {
	return func(t testing.TB, actualQuery string, actualValues []interface{}) {
		if actualQuery != expectedQuery {
			t.Errorf("Unexpected query string. Expected: %s. Got: %s", expectedQuery, actualQuery)
		}
		if len(actualValues) != len(expectedValues) {
			t.Errorf("Mismatching values. Expected length: %d. Actual length: %d", len(expectedValues), len(actualValues))
		} else {
			for index, actualValue := range actualValues {
				expectedValue := expectedValues[index]
				if actualValue != expectedValue {
					t.Errorf("Mismatching value at %d position. Expected: %s. Got: %s", index, expectedValue, actualValue)
				}
			}
		}
	}
}

func TestQueryBuilder(t *testing.T) {
	tests := []struct {
		name         string
		baseQuery    string
		kvCriteria   map[string]interface{}
		checkBuilder queryBuilderChecker
		checkResult  queryBuilderResultChecker
	}{
		{
			name:         "base query without criteria",
			baseQuery:    "select * from channels",
			kvCriteria:   map[string]interface{}{},
			checkBuilder: CheckQueryBuilderNoError(),
			checkResult:  CheckQueryBuilderResult("select * from channels", []interface{}{}),
		},
		{
			name:      "base query with criteria",
			baseQuery: "select * from channels",
			kvCriteria: map[string]interface{}{
				"name":          "test",
				"nb_partitions": 6,
			},
			checkBuilder: CheckQueryBuilderNoError(),
			checkResult: CheckQueryBuilderResult(
				"select * from channels where name = $1 and nb_partitions = $2",
				[]interface{}{"test", 6}),
		},
		{
			name:         "empty base query",
			baseQuery:    "",
			kvCriteria:   map[string]interface{}{},
			checkBuilder: CheckQueryBuilderError(builder.ErrEmptyBaseQuery),
			checkResult:  nil,
		},
		{
			name:      "empty base query with criteria",
			baseQuery: "",
			kvCriteria: map[string]interface{}{
				"name":          "test",
				"nb_partitions": 6,
			},
			checkBuilder: CheckQueryBuilderError(builder.ErrEmptyBaseQuery),
			checkResult:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			options := []builder.QueryBuilderOption{}

			if test.kvCriteria != nil && len(test.kvCriteria) > 0 {
				options = append(options, builder.WithKVCriteria(test.kvCriteria))
			}

			qb, err := builder.NewQueryBuilder(test.baseQuery, options...)
			test.checkBuilder(t, qb, err)
			if qb != nil && test.checkResult != nil {
				query, values := qb.Build()
				test.checkResult(t, query, values)
			}
		})
	}
}

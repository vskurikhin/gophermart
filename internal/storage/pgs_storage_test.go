/*
 * This file was last modified at 2024-04-21 01:01 by Victor N. Skurikhin.
 * pgs_storage_test.go
 * $Id$
 */

package storage

import (
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gophermart/internal/utils"
	"testing"
)

//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestPgsStoragePositive(t *testing.T) {

	storage := NewPgsStorage().(*PgsStorage)
	assert.NotNil(t, storage)

	var tests = []struct {
		name string
		sql  string
		args []interface{}
	}{
		{"Test positive #1", "SELECT version()", []interface{}{}},
		{"Test positive #2", "SELECT $1 + $2", []interface{}{1, 2}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := storage.sqlRow(test.name, test.sql, test.args...)
			assert.Nil(t, err)
		})
	}
}

//goland:noinspection GoImportUsedAsName,GoUnhandledErrorResult
func TestDaoSaveGet(t *testing.T) {

	store := NewPgsStorage().(*PgsStorage)
	store.sqlRow("create",
		`CREATE TABLE IF NOT EXISTS testtest (id int, t1 text, t2 text, UNIQUE(id), UNIQUE(t1), UNIQUE(t1, t2))`,
	)

	var id int
	var first *testtest
	var tests = []struct {
		name string
		repo interface{}
		init func() (interface{}, error)
		get  func() (interface{}, error)
	}{
		{
			name: "Test #1 Insert Save(testtest) GetByID(testtest)",
			init: func() (interface{}, error) {
				test := &testtest{id: 42, t1: "test"}
				e, err := test.save(store)
				assert.Nil(t, err)
				assert.NotNil(t, e)
				id = e.id
				return e, err
			},
			get: func() (interface{}, error) {
				var err error
				first, err = getTestByID(store, id)
				return first, err
			},
		},
		{
			name: "Test #2 Upsert Save(testtest) GetBiString(testtest)",
			init: func() (interface{}, error) {
				first.t1 = "testtest"
				first.t2 = "test"
				e, err := first.save(store)
				assert.Nil(t, err)
				assert.NotNil(t, e)
				id = e.id
				return e, err
			},
			get: func() (interface{}, error) {
				return getTestByString(store, "testtest")
			},
		},
		{
			name: "Test #3 Upsert Save(testtest) GetBiString(testtest)",
			init: func() (interface{}, error) {
				first.t1 = "testtest"
				first.t2 = "testtest"
				e, err := first.save(store)
				assert.Nil(t, err)
				assert.NotNil(t, e)
				id = e.id
				return e, err
			},
			get: func() (interface{}, error) {
				return getTestByStr1Str2(store, "testtest", "testtest")
			},
		},
		{
			name: "Test #4 Upsert Save(testtest) GetBiString(testtest)",
			init: func() (interface{}, error) {
				arr := [][]interface{}{
					{11, "t1_11", "t2"},
					{12, "t1_12", "t2"},
					{13, "t1_13", "t2"},
				}
				s := store.WithContext(utils.NewIDContext())
				a := make(TxArgs, 0)

				sql := "INSERT INTO testtest (id, t1, t2) VALUES ($1, $2, $3)"
				for _, i := range arr {
					a.Append(sql, i...)
				}
				err := s.Transaction(a)
				rows, _ := s.GetAllForString("SELECT * FROM testtest WHERE t2 = $1", "t2")
				result := make([]testtest, 0)
				for rows.Next() {
					tmp, _ := extractTestTest(rows)
					result = append(result, *tmp)
				}
				return result, err
			},
			get: func() (interface{}, error) {
				rows, err := store.GetAllForString("SELECT * FROM testtest WHERE t2 = $1", "t2")
				result := make([]testtest, 0)
				for rows.Next() {
					tmp, _ := extractTestTest(rows)
					result = append(result, *tmp)
				}
				return result, err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test1, err := test.init()
			assert.Nil(t, err)
			test2, err := test.get()
			assert.Nil(t, err)
			assert.Equal(t, test1, test2)
		})
	}
	store.sqlRow("drop", "DROP TABLE IF EXISTS testtest")
}

type testtest struct {
	id int
	t1 string
	t2 string
}

func getTestByID(st Storage, id int) (*testtest, error) {
	row, err := st.GetByID(`SELECT * FROM testtest WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	id, t1, t2, err := extractTestTestTuple(row)
	if err != nil {
		return nil, err
	}
	return &testtest{id: id, t1: t1, t2: t2}, nil
}

func getTestByString(st Storage, test1 string) (*testtest, error) {
	row, err := st.GetByString(
		`SELECT * FROM testtest WHERE t1 = $1`,
		test1,
	)
	if err != nil {
		return nil, err
	}
	return extractTestTest(row)
}

func getTestByStr1Str2(st Storage, test1, test2 string) (*testtest, error) {
	row, err := st.GetByStr1Str2(
		`SELECT * FROM testtest WHERE t1 = $1 AND t2 = $2`,
		test1, test2,
	)
	if err != nil {
		return nil, err
	}
	return extractTestTest(row)
}

func (s *testtest) save(st Storage) (*testtest, error) {
	row, err := st.Save(
		`INSERT INTO testtest (id, t1, t2) VALUES ($1, $2, $3)
	         ON CONFLICT (id)
	         DO UPDATE SET t1 = $2, t2 = $3
			 RETURNING *`,
		s.id, s.t1, s.t2,
	)
	if err != nil {
		return nil, err
	}
	return extractTestTest(row)
}

func extractTestTest(row pgx.Row) (*testtest, error) {
	id, t1, t2, err := extractTestTestTuple(row)
	if err != nil {
		return nil, err
	}
	return &testtest{id: id, t1: t1, t2: t2}, nil
}

func extractTestTestTuple(row pgx.Row) (int, string, string, error) {
	var i int
	var t1, t2 string
	err := row.Scan(&i, &t1, &t2)
	return i, t1, t2, err
}

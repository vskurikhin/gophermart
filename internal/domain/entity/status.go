/*
 * This file was last modified at 2024-04-12 13:27 by Victor N. Skurikhin.
 * status.go
 * $Id$
 */

package entity

import (
	"github.com/vskurikhin/gophermart/internal/storage"
	"time"
)

type Status struct {
	id        int
	status    string
	createdAt time.Time
	updateAt  time.Time
}

func (s *Status) Id() int {
	return s.id
}

func (s *Status) Status() string {
	return s.status
}

func (s *Status) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Status) UpdateAt() time.Time {
	return s.updateAt
}

func (s *Status) Save(st storage.Storage) error {

	_, err := st.Save(
		`INSERT INTO status (status) VALUES ($1)
	         ON CONFLICT (id)
	         DO UPDATE SET status = $1
			  RETURNING *`,
		s.status,
	)
	return err
}

func FuncGetAllStatuss() func(storage.Storage) ([]*Status, error) {
	return func(s storage.Storage) ([]*Status, error) {
		result := make([]*Status, 0)
		return result, nil
	}
}

func FuncGetStatus() func(storage.Storage, int) (*Status, error) {
	return func(s storage.Storage, id int) (*Status, error) {

		row, err := s.GetById("SELECT status FROM status WHERE id = $1", id)
		if err != nil {
			return nil, err
		}
		var status string
		err = row.Scan(&status)

		if err != nil {
			return nil, err
		}
		return &Status{id: id, status: status}, nil
	}
}

/*
 * This file was last modified at 2024-04-20 17:09 by Victor N. Skurikhin.
 * accrual.go
 * $Id$
 */

package model

import (
	"github.com/mailru/easyjson"
	"io"
)

type Accrual struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual Float  `json:"accrual"`
}

func (v *Accrual) GetAccrual() float64 {
	return float64(v.Accrual)
}

func UnmarshalAccrualFromReader(reader io.Reader) (*Accrual, error) {

	accrual := new(Accrual)

	if err := easyjson.UnmarshalFromReader(reader, accrual); err != nil {
		return nil, err
	}
	return accrual, nil
}

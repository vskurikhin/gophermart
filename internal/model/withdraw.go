/*
 * This file was last modified at 2024-05-07 17:36 by Victor N. Skurikhin.
 * withdraw.go
 * $Id$
 */

package model

import (
	"github.com/mailru/easyjson"
	"io"
	"math/big"
)

type Withdraw struct {
	Order       string `json:"order" example:"2377225624"`
	Sum         Float  `json:"sum" example:"751"`
	ProcessedAt *Time  `json:"processed_at,omitempty" example:"2020-12-09T16:09:57+03:00"`
}

func NewWithdraw(order string, sum *big.Float) *Withdraw {
	f, _ := sum.Float64()
	return &Withdraw{Order: order, Sum: Float(f)}
}

func (w *Withdraw) GetSum() float64 {
	return float64(w.Sum)
}

func UnmarshalWithdrawFromReader(reader io.Reader) (*Withdraw, error) {

	withdraw := new(Withdraw)

	if err := easyjson.UnmarshalFromReader(reader, withdraw); err != nil {
		return nil, err
	}
	return withdraw, nil
}

func (w *Withdraw) MarshalToWriter(writer io.Writer) error {

	if _, err := easyjson.MarshalToWriter(w, writer); err != nil {
		return err
	}
	return nil
}

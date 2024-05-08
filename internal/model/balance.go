/*
 * This file was last modified at 2024-05-07 17:40 by Victor N. Skurikhin.
 * balance.go
 * $Id$
 */

package model

import (
	"github.com/mailru/easyjson"
	"io"
	"math/big"
)

type Balance struct {
	Current   Float `json:"current" example:"500.5"`
	Withdrawn Float `json:"withdrawn" example:"42"`
}

func NewBalanceBigFloat(current big.Float, withdrawn big.Float) *Balance {
	c, _ := current.Float64()
	w, _ := withdrawn.Float64()
	return &Balance{Current: Float(c), Withdrawn: Float(w)}
}

func (b *Balance) MarshalToWriter(writer io.Writer) error {

	if _, err := easyjson.MarshalToWriter(b, writer); err != nil {
		return err
	}
	return nil
}

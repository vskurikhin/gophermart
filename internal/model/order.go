/*
 * This file was last modified at 2024-04-19 16:34 by Victor N. Skurikhin.
 * order.go
 * $Id$
 */

package model

import (
	"github.com/mailru/easyjson"
	"io"
)

type Order struct {
	Number     string `json:"number"`
	Accrual    int    `json:"accrual,omitempty"`
	UploadedAt Time   `json:"uploaded_at"`
}

func (o *Order) MarshalToWriter(writer io.Writer) error {

	if _, err := easyjson.MarshalToWriter(o, writer); err != nil {
		return err
	}
	return nil
}

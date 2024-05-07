/*
 * This file was last modified at 2024-05-07 17:07 by Victor N. Skurikhin.
 * order.go
 * $Id$
 */

package model

import (
	"github.com/mailru/easyjson"
	"io"
)

type Order struct {
	Number     string `json:"number" example:"9278923470"`
	Status     string `json:"status" example:"PROCESSED"`
	Accrual    *Float `json:"accrual,omitempty" example:"500"`
	UploadedAt Time   `json:"uploaded_at" example:"2020-12-10T15:15:45+03:00"`
}

func (o *Order) MarshalToWriter(writer io.Writer) error {

	if _, err := easyjson.MarshalToWriter(o, writer); err != nil {
		return err
	}
	return nil
}

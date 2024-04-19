/*
 * This file was last modified at 2024-04-20 00:34 by Victor N. Skurikhin.
 * withdrawals.go
 * $Id$
 */

package model

import (
	"github.com/mailru/easyjson"
	"io"
)

// easyjson:json
type Withdrawals []Withdraw

func (v Withdrawals) MarshalToWriter(writer io.Writer) error {

	if _, err := easyjson.MarshalToWriter(v, writer); err != nil {
		return err
	}
	return nil
}

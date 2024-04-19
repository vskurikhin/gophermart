/*
 * This file was last modified at 2024-04-19 16:34 by Victor N. Skurikhin.
 * orders.go
 * $Id$
 */

package model

import (
	"github.com/mailru/easyjson"
	"io"
)

// easyjson:json
type Orders []Order

func (v Orders) MarshalToWriter(writer io.Writer) error {

	if _, err := easyjson.MarshalToWriter(v, writer); err != nil {
		return err
	}
	return nil
}

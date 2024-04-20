/*
 * This file was last modified at 2024-04-20 17:08 by Victor N. Skurikhin.
 * empty.go
 * $Id$
 */

package model

import "net/http"

type Empty struct {
}

//goland:noinspection GoUnusedParameter
func (e *Empty) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

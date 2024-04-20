/*
 * This file was last modified at 2024-04-20 17:09 by Victor N. Skurikhin.
 * empty_list.go
 * $Id$
 */

package model

import "net/http"

type EmptyList []struct{}

func NewEmptyList() EmptyList {
	return make([]struct{}, 0)
}

//goland:noinspection GoUnusedParameter
func (e EmptyList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

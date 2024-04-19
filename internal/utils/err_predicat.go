/*
 * This file was last modified at 2024-04-19 18:11 by Victor N. Skurikhin.
 * err_predicat.go
 * $Id$
 */

package utils

func IsErrNoRowsInResultSet(err error) bool {
	return err != nil && err.Error() == "no rows in result set"
}

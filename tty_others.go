// +build appengine js

package logftext

import "errors"

func enableSeqTTY(fd uintptr, flag bool) error {
	return errors.New("default not a terminal")
}

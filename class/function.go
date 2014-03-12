package class

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
)

func EncryptPassword(str string) (pwd string, err error) {
	m := md5.New()
	_, err = io.WriteString(m, str)
	p1 := fmt.Sprintf("%x", m.Sum(nil))
	s := sha1.New()
	_, err = io.WriteString(s, str)
	p2 := fmt.Sprintf("%x", s.Sum(nil))
	pwd = p1 + p2
	return
}

package CHelperSecurity

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

type Security struct {}

func NewSecurity() interface{} {return new(Security)}

func (*Security) Sum256(str string) string {
	data := sha256.Sum256([]byte(str))
	return hex.EncodeToString(data[:])
}

func (*Security) Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

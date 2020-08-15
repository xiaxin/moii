package crypto

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func Md5(string string) string {
	crypto := md5.New()

	io.WriteString(crypto, string)

	sum := crypto.Sum(nil)

	return hex.EncodeToString(sum)
}

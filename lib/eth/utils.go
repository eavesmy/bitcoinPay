package eth

import "encoding/hex"

func StringPadding64(str string) string {

	if str[:2] == "0x" {
		str = str[2:]
	}

	reduce := 64 - len(str)

	for i := 0; i < reduce; i++ {
		str = "0" + str
	}
	return str
}

func Str2Hex(str string) string {
	return hex.EncodeToString([]byte(str))
}

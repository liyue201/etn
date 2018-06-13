package rest

const (
	OK         = 0      //成功
	BadRequest = 400001 // 参数不合法，请检查参数
)

var statusText = map[int]string{
	OK:         "Success",
	BadRequest: "参数不合法，请检查参数",
}

func StatusText(code int) string {
	return statusText[code]
}

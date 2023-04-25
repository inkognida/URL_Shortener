package utils

const (
	charset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	linkLength    = 10
	charsetLength = int64(len(charset))
)

// GenerateLink генерирует уникальную ссылку
func GenerateLink(shortId int64) string {
	res := make([]byte, 10)
	for i := 0; i < linkLength; i++ {
		if shortId > 0 {
			res[i] = charset[shortId%charsetLength]
			shortId /= charsetLength
		} else {
			res[i] = charset[shortId%charsetLength]
		}
	}
	return string(res)
}

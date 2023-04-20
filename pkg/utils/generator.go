package utils

const (
	charset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	linkLength    = 10
	charsetLength = int64(len(charset))
)

func GenerateLink(shortId int64) string {
	if shortId < 0 {
		shortId = -shortId
	}
	var b [11]byte
	for i := linkLength; shortId > 0 && i > 0; i-- {
		shortId, b[i] = shortId/charsetLength, charset[shortId%charsetLength]
	}
	return string(b[:])
}

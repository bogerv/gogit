package masking

// Uid desensit user id
func Uid(uid string) string {
	uidLen := len(uid)
	if uidLen < 3 {
		return uid + generateStar(4)
	}
	return uid[0:2] + generateStar(4) + uid[uidLen-1:uidLen]
}

func generateStar(num int) string {
	res := ""
	for i := 0; i < num; i++ {
		res += "*"
	}

	return res
}

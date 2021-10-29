package masking

import (
	"fmt"
	"strings"
	"testing"
)

func TestUid(t *testing.T) {
	uid := []string{"12345475", "1", "12", "123"}
	result := []string{"12****5", "1****", "12****", "12****3"}
	for i := 0; i < len(uid); i++ {
		u := Uid(uid[i])
		fmt.Printf("%s, %s\n", u, result[i])
		if !strings.EqualFold(u, result[i]) {
			t.Error()
		}
	}
}

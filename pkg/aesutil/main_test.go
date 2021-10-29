package aesutil

import (
	"fmt"
	"testing"
)

func TestAesCFBDecrypt(t *testing.T) {
	res, err := AesCFBDecrypt("4ee047cc9417073c11e3202170599e4e42c66d56865ee6019de74344ecedb656745fe817d41dd27c0092ad591287f1450faa6d057e29565843cba1dd3820f5569b5870dda4370c8e58c1d970b9356a36")
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	fmt.Println(res)
	if res != "b64bd26d973c03c5e0752729386deaf45c8cd136a8dc8bb6ce2e0dcd7a9e0867" {
		t.Error(res)
	}

	res, err = AesCFBDecrypt("bd00e42e39e94ac807b94516dca7c2a12b7efec69b0d4602e0b893b2b64f726bb884afc368b588e222fcd66d4b269620253d75ac2f87040332f93054006ad44fc6c0d9822df2dce2323acbd31f406fc8")
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	fmt.Println(res)
	if res != "e4b0d39fa3b3a3665ef7c9ac6da1c3c65a74a1255abee5466ff2673bf89d7f9c" {
		t.Error(res)
	}
}

func TestAesCFBEncrypt(t *testing.T) {
	res, err := AesCFBEncrypt("b64bd26d973c03c5e0752729386deaf45c8cd136a8dc8bb6ce2e0dcd7a9e0867")
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	fmt.Println(res)

	res, err = AesCFBEncrypt("e4b0d39fa3b3a3665ef7c9ac6da1c3c65a74a1255abee5466ff2673bf89d7f9c")
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	fmt.Println(res)
}

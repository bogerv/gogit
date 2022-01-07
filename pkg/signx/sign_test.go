package signx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const Value = "12312312&#*$^&*  dsjvsdfjsdlfj 1231"

func TestMD5(t *testing.T) {
	assert.Equal(t, MD5(Value), "fadb29a787dd84b4d6398b7f6cea3cf6")
}

func TestSha1(t *testing.T) {
	assert.Equal(t, Sha1(Value), "60c2d1f0e21dd69b5149d405c42f50924081d143")
}

func TestSha224(t *testing.T) {
	assert.Equal(t, Sha224(Value), "fbc2cc7c57a68681fcb491a42cfe7bad521e3a27e7d391367b415338")
}

func TestSha256(t *testing.T) {
	assert.Equal(t, Sha256(Value), "701c41aa070a9394dadfabcb2fd06dc4d6f579269fb89d028b30c8fa770c0e3a")
}

func TestSha384(t *testing.T) {
	assert.Equal(t, Sha384(Value), "136d6503669b4288f7a0fb5c00cdd8485bb52d04d96c6eb6c70788d807f0cb35096a3ac28c400e243337d22533d94c22")
}

func TestSha512(t *testing.T) {
	assert.Equal(t, Sha512(Value), "fb4c4e58182021e98268194834478f25aedf326cd82b8b2935b2de42af2dd2c982dbf134325110a3aa5162197f61480cbf594cae629d6de66e413c02a2a61481")
}

package utils

import "testing"

func TestGenInviteCode(t *testing.T) {
	res := make(map[string]int)
	for i := 1000000; i < 9000000; i++ {
		code := GenInviteCode(int64(i))
		_, ok := res[code]
		if ok {
			t.Fatal(i, code)
		}
		id, err := DecodeInviteCode(code)
		if err != nil {
			t.Fatal(err)
		}
		if int64(i) != id {
			t.Fatal(i, code, id)
		}
	}
}
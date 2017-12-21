package dec64

import (
	"testing"
)

func testOneDec(t *testing.T, s string, ref int64) {
	d, err := Parse(s)
	if err != nil {
		t.Error(err)
	}
	if ref != int64(d) {
		t.Errorf("Result is %d should be %d", d, ref)
	}
	if s != d.String() {
		t.Errorf("String is %s should be %s", d.String(), s)
	}
}

func TestParse(t *testing.T) {
	testOneDec(t, "1", 256)
	testOneDec(t, "-1", -256)
	testOneDec(t, "345", 345*256)
	testOneDec(t, ".6789", 1738236)
	testOneDec(t, "2.05", 52734)
}

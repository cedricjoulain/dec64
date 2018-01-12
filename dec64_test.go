package dec64

import (
	"math"
	"strconv"
	"testing"
)

func testOneDec(t *testing.T, s string, ref int64) {
	d, err := Parse(s)
	if err != nil {
		t.Error(err)
	}
	if ref != int64(d) {
		t.Errorf("%s Result is %d should be %d", s, d, ref)
	}
	if s != d.String() {
		t.Errorf("String is %s should be %s", d.String(), s)
	}
	f, _ := strconv.ParseFloat(s, 64)
	if math.Abs(f - Float64(d)) > 0.000000000000001 {
		t.Errorf("Float64 is %g should be %g", Float64(d), f)
	}
}

func TestParse(t *testing.T) {
	testOneDec(t, "1", 256)
	testOneDec(t, "-1", -256)
	testOneDec(t, "345", 345*256)
	testOneDec(t, ".6789", 1738236)
	testOneDec(t, "2.05", 52734)
	testOneDec(t, "300201", 76851456)
	testOneDec(t, "100", 258)
	testOneDec(t, ".09", 9 * 256 + 254)
	testOneDec(t, ".007", 7 * 256 + 253)
	testOneDec(t, "0", 0)
}

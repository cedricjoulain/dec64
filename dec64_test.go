package dec64

import (
	"math"
	"strconv"
	"testing"
)

func testOneDec(t *testing.T, s, refs string, ref int64) {
	d, err := Parse(s)
	if err != nil {
		t.Error(err)
	}
	if ref != int64(d) {
		t.Errorf("%s Result is %d should be %d", s, d, ref)
	}
	if refs != d.String() {
		t.Errorf("String is %s should be %s", d.String(), refs)
	}
	f, _ := strconv.ParseFloat(s, 64)
	if math.Abs(f-Float64(d)) > 0.000000000000001 {
		t.Errorf("Float64 is %g should be %g", Float64(d), f)
	}
}

func TestParse(t *testing.T) {
	testOneDec(t, "1", "1", 256)
	testOneDec(t, "-1", "-1", -256)
	testOneDec(t, "345", "345", 345*256)
	testOneDec(t, ".6789", ".6789", 1738236)
	testOneDec(t, "2.05", "2.05", 52734)
	testOneDec(t, "300201", "300201", 76851456)
	testOneDec(t, "6151000", "6151000", 6151*256+3)
	testOneDec(t, "100", "100", 258)
	testOneDec(t, ".09", ".09", 9*256+254)
	testOneDec(t, ".007", ".007", 7*256+256-3)
	testOneDec(t, "0", "0", 0)
	testOneDec(t, "1E-8", ".00000001", 1*256+256-8)
	testOneDec(t, "3E-8", ".00000003", 3*256+256-8)
	testOneDec(t, "1.2E-7", ".00000012", 12*256+256-8)
	testOneDec(t, "1.000506e+06", "1000506", 1000506*256)
	testOneDec(t, "6.151e+06", "6151000", 6151*256+3)
}

func testOneFloat(t *testing.T, f float64, ref int64) {
	d, err := FromFloat64(f)
	if err != nil {
		t.Error(err)
	}
	if ref != int64(d) {
		t.Errorf("%g Result is %d should be %d", f, d, ref)
	}
	// Less accurante than from string
	if math.Abs(f-Float64(d)) > 0.000000001 {
		t.Errorf("Float64 is %g should be %g", Float64(d), f)
	}
}

func TestFromFloat(t *testing.T) {
	testOneFloat(t, 0.09112614, 2332829432)
	testOneFloat(t, 6.0, 6*256)
	testOneFloat(t, 24.7165, 63274492)
	testOneFloat(t, 58803.0596245, 150535832638969)
	testOneFloat(t, 6.35e-06, 162808)
	testOneFloat(t, 912550.0000003, 2336128000001017)
	testOneFloat(t, 999.99999986, 25599999996664)
	testOneFloat(t, 0.00023224, 5945592)
}

func testOneInt64(t *testing.T, i int64, ref int64) {
	d, err := FromInt64(i)
	if err != nil {
		t.Error(err)
	}
	if ref != int64(d) {
		t.Errorf("%d Result is %d should be %d", i, d, ref)
	}
	if Int64(d) != i {
		t.Errorf("Int64 is %d should be %d", Int64(d), i)
	}
}

func TestFromInt64(t *testing.T) {
	testOneInt64(t, 0, 0)
	testOneInt64(t, 1, 256)
	testOneInt64(t, -1, -256)
	testOneInt64(t, 37, 37*256)
	testOneInt64(t, 2, 2*256)
}

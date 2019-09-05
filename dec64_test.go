package dec64

import (
	"bytes"
	"fmt"
	"io"
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

func testOneDecPrecision(t *testing.T, s, refs string, ref, n int64) {
	d, err := Parse(s)
	if err != nil {
		t.Error(err)
	}
	d = Round(d, n)
	if ref != int64(d) {
		t.Errorf("%s Result is %d should be %d", s, d, ref)
	}
	if refs != d.String() {
		t.Errorf("String is %s should be %s", d.String(), refs)
	}
}

func TestParse(t *testing.T) {
	testOneDec(t, "+1", "1", 256)
	testOneDec(t, "1", "1", 256)
	testOneDec(t, "-1", "-1", -256)
	testOneDec(t, "345", "345", 345*256)
	testOneDec(t, ".6789", "0.6789", 1738236)
	testOneDec(t, "2.05", "2.05", 52734)
	testOneDec(t, "300201", "300201", 76851456)
	testOneDec(t, "6151000", "6151000", 6151*256+3)
	testOneDec(t, "100", "100", 258)
	testOneDec(t, ".09", "0.09", 9*256+254)
	testOneDec(t, ".007", "0.007", 7*256+256-3)
	testOneDec(t, "0", "0", 0)
	testOneDec(t, "1E-8", "0.00000001", 1*256+256-8)
	testOneDec(t, "3E-8", "0.00000003", 3*256+256-8)
	testOneDec(t, "1.2E-7", "0.00000012", 12*256+256-8)
	testOneDec(t, "1.000506e+06", "1000506", 1000506*256)
	testOneDec(t, "6.151e+06", "6151000", 6151*256+3)
	testOneDec(t, "7003.69", "7003.69", 700369*256+254)
	testOneDec(t, "7003.1", "7003.1", 70031*256+255)
	testOneDec(t, "0.0041210199999999995", "0.00412102",
		4121020000000000*256+256-18)
	testOneDec(t, "3.1997721799999996", "3.1997721799999996",
		31997721799999996*256+256-16)
	testOneDec(t, "0.0180000000000000002", "0.018", 4608000000000000238)
	testOneDec(t, "-0.0041210199999999995", "-0.00412102",
		-4121020000000000*256+256-18)
	testOneDec(t, "123456789012345678", "123456789012345680",
		12345678901234568*256+1)
	testOneDec(t, "1234567890123456780", "1234567890123456800",
		12345678901234568*256+2)
	testOneDec(t, "-123456789012345678", "-123456789012345680",
		-12345678901234568*256+1)
	testOneDec(t, "-1234567890123456780", "-1234567890123456800",
		-12345678901234568*256+2)
	testOneDecPrecision(t, "79.068001", "79.068", 20241661, -3)
	testOneDecPrecision(t, "-79.068001", "-79.068", -20241155, -3)
	testOneDecPrecision(t, "79.084999", "79.085", 20246013, -3)
	testOneDecPrecision(t, "-79.084999", "-79.085", -20245507, -3)
	// check same behavior as math.Round
	testOneDecPrecision(t, "4.5", "5", 1280, 0)
	testOneDecPrecision(t, "4.51", "5", 1280, 0)
	testOneDecPrecision(t, "4.49", "4", 1024, 0)
	testOneDecPrecision(t, "-4.5", "-5", -1280, 0)
	testOneDecPrecision(t, "-4.51", "-5", -1280, 0)
	testOneDecPrecision(t, "-4.49", "-4", -1024, 0)
	testOneDecPrecision(t, "0", "0", 0, 0)
	testOneDecPrecision(t, "125420.000", "125000", 32003, 3)
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

func testOneInt64(t *testing.T, i, ref int64) {
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

func testOneEqual(t *testing.T, a, b Dec64, ref bool) {
	if a.Equal(b) != ref {
		t.Errorf("%s Equal %s should be %t", a.String(), b.String(), ref)
	}
}

func TestEqual(t *testing.T) {
	d, _ := FromInt64(1)
	testOneEqual(t, d, Dec64(256), true)
	// 10
	testOneEqual(t, Dec64(2560), Dec64(257), true)
	testOneEqual(t, Dec64(256), Dec64(512), false)
	// All 0 in fact
	testOneEqual(t, Dec64(2), Dec64(3), true)
	testOneEqual(t, Empty, Dec64(123456), true)
	testOneEqual(t, Dec64(123456), Empty, true)
	testOneEqual(t, NotAvailable, Dec64(123456), true)
	testOneEqual(t, Dec64(123456), NotAvailable, true)
}

func testOneIsInt(t *testing.T, s string, ref bool) {
	d, err := Parse(s)
	if err != nil {
		t.Error(err)
	}
	if d.IsInt() != ref {
		t.Errorf("%s IsInt is %t should be %t", s, d.IsInt(), ref)
	}
}

func TestIsInt(t *testing.T) {
	testOneIsInt(t, "0.0002", false)
	testOneIsInt(t, "20", true)
	testOneIsInt(t, "-0.000452", false)
	testOneIsInt(t, "-15487920", true)
	testOneIsInt(t, "-255.00000", true)
	testOneIsInt(t, "255.00000", true)
}

// example of list of traded volumes for BTC on 20180511
var sVolumes = []string{
	"0.06447466",
	"0.0244",
	"0.04082425",
	"0.18808132",
	"0.00410578",
	"0.02433487",
	"0.14442756",
	"0.003216",
	"0.00096851",
	"0.00153346",
	"0.00680321",
	"0.00204112",
	"0.08067412",
	"0",
	"0.02662033",
	"0.04217344",
	"1",
}

var sVBench = []string{
	"0.06447466",
	"2.44",
	"4082425",
	"0.18808132",
	"0.00410578",
	"2433400",
	"0.14442756",
	"0.003216",
	"0.00096851",
	"15.3346",
	"-0.00680321",
	"0.00204112",
	"0.08067412",
	"0",
	"1",
	"10000",
	"200",
	"85.236",
}

func BenchmarkString(b *testing.B) {
	var (
		err error
		s   string
	)
	dVolumes := make([]Dec64, len(sVBench))
	for i, v := range sVBench {
		dVolumes[i], err = Parse(v)
		if err != nil {
			b.Errorf(err.Error())
			return
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, d64 := range dVolumes {
			s = d64.String()
		}
	}
	b.SetBytes(8 * int64(len(dVolumes)))
	if s != sVBench[len(sVBench)-1] {
		b.Errorf("wrong d64 to string")
	}
}

func TestLists(t *testing.T) {
	var (
		err error
	)
	dVolumes := make([]Dec64, len(sVolumes))
	for i, v := range sVolumes {
		dVolumes[i], err = Parse(v)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
	}
	Homogenize(dVolumes)
	for i, ref := range sVolumes {
		v := dVolumes[i].String()
		if v != ref {
			t.Errorf("sVolumes[%d] is %s should %s", i, v, ref)
		}
	}
}

func TestOverflow(t *testing.T) {
	ref := int64(123578)
	s := fmt.Sprintf("%d.458", ref)
	d, err := Parse(s)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if Int64(d) != ref {
		t.Errorf("Int64(%s) is %d should %d", d.String(), Int64(d), ref)
		return
	}
	// Huge and small
	values := make([]Dec64, 2)
	values[0], err = Parse("0.000000000000001")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	values[1], err = Parse("100000000000000")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	// values could be modified
	refs := make([]Dec64, 2)
	for i, v := range values {
		refs[i] = v
	}
	Homogenize(values)
	// Now check it's ok
	for i, v := range values {
		if !v.Equal(refs[i]) {
			t.Errorf("values[%d] is %s should %s", i, v.String(), refs[i].String())
			return
		}
	}
}

// Test errors and border line behavior to ensure full coverage
func TestBorders(t *testing.T) {
	testOneDec(t, "", "0", 1)

	testParseError(t, "3.2.5", "Only one dot allowed")
	testParseError(t, "toto", "Unable to parse dec64 from toto")
	testParseError(t, "3.7ea", "Unable to handle a in exponent")
	// to big or too small
	toBig := "10"
	toSmall := "."
	for i := 0; i < 127; i++ {
		toBig += "0"
		toSmall += "0"
	}
	toSmall += "1"
	testParseError(t, toBig, fmt.Sprintf("%s is too big for dec64", toBig))
	testParseError(t, toSmall, fmt.Sprintf("%s is too small for dec64", toSmall))
}

func testParseError(t *testing.T, value, ref string) {
	_, err := Parse(value)
	if err == nil {
		t.Errorf("Parsing %s should be in error", value)
		return
	}
	if err.Error() != ref {
		t.Errorf("Parsing %s error is \n%s should be \n%s", value, err.Error(), ref)
	}
}

func TestListRW(t *testing.T) {
	var buf bytes.Buffer
	values := []Dec64{Dec64(1 * 256), Dec64(-1 * 256), Dec64(100 * 256)}

	err := ListToWriter(&buf, values)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	read, err := ListFromReader(&buf)
	if err != nil {
		if err != io.EOF {
			t.Errorf(err.Error())
			return
		}
	}
	if len(read) != len(values) {
		t.Errorf("len is %d should be %d", len(read), len(values))
		return
	}
	for i, value := range values {
		if !value.Equal(read[i]) {
			t.Errorf(
				"read[%d] is %s should be %s", i, read[i].String(), value.String())
		}
	}
}

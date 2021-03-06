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
		t.Errorf("%s Result is %d (%d*256+%d) should be %d", s, d, d/256, d%256, ref)
	}
	if refs != d.String() {
		t.Errorf("String is %s should be %s", d.String(), refs)
	}
	json, err := d.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	if refs != string(json) {
		t.Errorf("JSON is %s should be %s", string(json), refs)
	}
	f, _ := strconv.ParseFloat(refs, 64)
	if math.Abs(f-Float64(d)) > Epsilon {
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
	testOneDec(t, " +1", "1", 256)
	testOneDec(t, " 1", "1", 256)
	testOneDec(t, "  -1", "-1", -256)
	testOneDec(t, "345", "345", 345*256)
	testOneDec(t, ".6789", "0.6789", 1738236)
	testOneDec(t, "2.05", "2.05", 52734)
	testOneDec(t, "300201", "300201", 76851456)
	testOneDec(t, "6151000", "6151000", 6151*256+3)
	testOneDec(t, "100", "100", 1*256+2)
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
	testOneDec(t, "19831697.120000001043", "19831697.120000001", 19831697120000001*256+247)
	testOneDec(t, "-19831697.120000001043", "-19831697.120000001", -19831697120000001*256+247)
	testOneDec(t, "26414364.620000001043", "26414364.620000001", 26414364620000001*256+247)
	testOneDec(t, "-26414364.620000001043", "-26414364.620000001", -26414364620000001*256+247)
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

	// huge
	testOneDec(t, "297011256.990000009537", "297011256.99000001", 7603488178944000504)
	// rounding issue
	testOneDec(t, "136.33999999999997", "136.33999999999997", 13633999999999997*256+242)
	// bugged!
	testOneDec(t, "6.0196843678775026e-05", "0.00006019684367877503", 6019684367877503*256+236)
	testOneDec(t, "-6.0196843678775026e-05", "-0.00006019684367877503", -6019684367877503*256+236)
	testOneDec(t, "4.2617752851088906e-05", "0.00004261775285108891", 4261775285108891*256+236)
	testOneDec(t, "4.8371493066708666e-05", "0.00004837149306670867", 4837149306670867*256+236)
	testOneDec(t, "4.7550471428916355e-05", "0.00004755047142891636", 4755047142891636*256+236)
	testOneDec(t, "-4.7550471428916355e-05", "-0.00004755047142891636", -4755047142891636*256+236)
	testOneDec(t, "6.3340047378297835e-06", "0.000006334004737829784", 6334004737829784*256+235)
	testOneDec(t, "4.2847524025179855e-05", "0.00004284752402517986", 4284752402517986*256+236)
	testOneDec(t, "6.0915251656078004e-05", "0.000060915251656078", 6091525165607800*256+236)
}

func testOneFloat(t *testing.T, f float64, ref int64) {
	d, err := FromFloat64(f)
	if err != nil {
		t.Error(err)
	}
	if ref != int64(d) {
		t.Errorf("%g Result is %d (%d*256+%d) should be %d", f, d, d/256, d%256, ref)
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
		t.Errorf("%d Result is %d (%d*256+%d) should be %d", i, d, d/256, d%256, ref)
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

func testOneSignum(t *testing.T, s string, ref int) {
	d, err := Parse(s)
	if err != nil {
		t.Error(err)
	}
	if Signum(d) != ref {
		t.Errorf("%s Signum is %d should be %d", s, Signum(d), ref)
	}
}

func TestSignum(t *testing.T) {
	testOneSignum(t, "0.0002", 1)
	testOneSignum(t, "20", 1)
	testOneSignum(t, "-0.000452", -1)
	testOneSignum(t, "-15487920", -1)
	testOneSignum(t, "-255.00000", -1)
	testOneSignum(t, "255.00000", 1)
	testOneSignum(t, "0", 0)
	testOneSignum(t, "0.0", 0)
	testOneSignum(t, "0.000000", 0)
}
func testOneRound(t *testing.T, n int64, v, ref string) {
	dv, err := Parse(v)
	if err != nil {
		t.Error(err)
		return
	}
	dref, err := Parse(ref)
	if err != nil {
		t.Error(err)
		return
	}
	dv = Round(dv, n)
	if !dref.Equal(dv) {
		t.Errorf("Round(%s %d) = %s should be %s", v, n, dv.String(), ref)
	}
}

func TestRound(t *testing.T) {
	testOneRound(t, 0, "0", "0")
	testOneRound(t, 0, "0.4", "0")
	testOneRound(t, 0, "0.5", "1")
	testOneRound(t, 0, "0.9", "1")
	testOneRound(t, 0, "1.4", "1")
	testOneRound(t, 0, "1.5", "2")
	testOneRound(t, 0, "1.9", "2")
	testOneRound(t, 0, "-0.4", "0")
	testOneRound(t, 0, "-0.5", "-1")
	testOneRound(t, 0, "-0.9", "-1")
	testOneRound(t, 0, "-1.4", "-1")
	testOneRound(t, 0, "-1.5", "-2")
	testOneRound(t, 0, "-1.9", "-2")
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
			b.Error(err)
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
			t.Error(err)
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
		t.Error(err)
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
		t.Error(err)
		return
	}
	values[1], err = Parse("100000000000000")
	if err != nil {
		t.Error(err)
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
	testOneDec(t, "", "null", 1)

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
	for i := 112; i < 128; i++ {
		toBig += "0"
	}
	toSmall += "1"
	testParseError(t, toBig, fmt.Sprintf("%s is too big for dec64", toBig))
	testParseError(t, toSmall, fmt.Sprintf("%s is too small for dec64", toSmall))
}

func testParseError(t *testing.T, value, ref string) {
	d, err := Parse(value)
	if err == nil {
		t.Errorf("Result is %d (%d*256+%d)", d, d/256, d%256)
		t.Errorf("Parsing %s should be in error", value)
		return
	}
	if err.Error() != ref {
		t.Errorf("Result is %d (%d*256+%d)", d, d/256, d%256)
		t.Errorf("Parsing %s error is \n%s should be \n%s", value, err.Error(), ref)
	}
}

func TestListRW(t *testing.T) {
	var buf bytes.Buffer
	values := []Dec64{Dec64(1 * 256), Dec64(-1 * 256), Dec64(100 * 256)}

	err := ListToWriter(&buf, values)
	if err != nil {
		t.Error(err)
		return
	}
	read, err := ListFromReader(&buf)
	if err != nil {
		if err != io.EOF {
			t.Error(err)
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

// Dec64 mutiplied to a int huge enough to make it an int64
func testMultInt64(t *testing.T, d Dec64, i, ref int64) {
	d = d.MultInt64(i)
	if ref != Int64(d) {
		t.Errorf("%d Result is %d should be %d", i, d, ref)
	}
}

func TestMultInt64(t *testing.T) {
	testMultInt64(t, Dec64(-4*256), 3600, -4*3600)
	testMultInt64(t, Dec64(2*256), 3600, 2*3600)
	// Kabul like +2.5 hours
	testMultInt64(t, Dec64(25*256+256-1), 3600, 9000)
	testMultInt64(t, Dec64(-25*256+256-1), 3600, -9000)
}

func TestNormalize(t *testing.T) {
	d := Dec64(10 * 256)
	ref := int64(1*256 + 1)
	if ref != int64(Normalize(d)) {
		t.Errorf("%s normalized is %d should be %d", d,
			int64(Normalize(d)), ref)
	}
	d = Dec64(-10 * 256)
	ref = -1*256 + 1
	if ref != int64(Normalize(d)) {
		t.Errorf("%s normalized is %d should be %d", d,
			int64(Normalize(d)), ref)
	}
	d = Dec64(-10*256 + (256 - 3))
	ref = -1*256 + (256 - 2)
	if ref != int64(Normalize(d)) {
		t.Errorf("%s normalized is %d should be %d", d,
			int64(Normalize(d)), ref)
	}
}

func TestNeg(t *testing.T) {
	one := Dec64(1 * 256)
	mone := Dec64(-1 * 256)
	// first ensure Neg isnot identity
	if !mone.Equal(one.Neg()) {
		t.Errorf("%s should be %s", one.Neg(), mone)
	}
	// check neg of neg!
	for _, v := range sVBench {
		d, err := Parse(v)
		if err != nil {
			t.Error(err)
			return
		}
		if !d.Equal(d.Neg().Neg()) {
			t.Errorf("%s should be %s", d.Neg().Neg(), d)
		}
	}
}

func testAdd(t *testing.T, a, b, ref Dec64) {
	if !ref.Equal(a.Add(b)) {
		t.Errorf("%s+%s is %s(%d) should be %s", a, b, a.Add(b), int64(a.Add(b)), ref)
	}
	if !ref.Equal(a.Sub(b.Neg())) {
		t.Errorf("%s-%s is %s should be %s", a, b.Neg(),
			a.Sub(b.Neg()), ref)
	}
	if !ref.Equal(b.Add(a)) {
		t.Errorf("%s+%s is %s should be %s", b, a, b.Add(a), ref)
	}
	if !ref.Equal(b.Sub(a.Neg())) {
		t.Errorf("%s-%s is %s should be %s", b, a.Neg(),
			b.Sub(a.Neg()), ref)
	}
}

func TestAdd(t *testing.T) {
	a := Dec64(1 * 256)
	b := Dec64(-1 * 256)
	ref := Dec64(0)
	testAdd(t, a, b, ref)
	// Huge and small
	var err error
	a, err = Parse("0.000000000000001")
	if err != nil {
		t.Error(err)
		return
	}
	b, err = Parse("100000000000000")
	if err != nil {
		t.Error(err)
		return
	}
	ref = b
	testAdd(t, a, b, ref)
	a = Dec64(9999999999999999*256 + 127)
	b = Dec64(1*256 + (256 - 128))
	ref = a
	testAdd(t, a, b, ref)

	a = Dec64(1*256 + 1)
	b = Dec64(-10 * 256)
	ref = Dec64(0)
	testAdd(t, a, b, ref)

	a = Dec64(-10*256 + (256 - 3))
	b = Dec64(1*256 + (256 - 2))
	testAdd(t, a, b, ref)

	a = Dec64(1*256 + (256 - 3))
	b = Dec64(1*256 + 2)
	ref = Dec64(100001*256 + (256 - 3))
	testAdd(t, a, b, ref)

	// a max precisiion
	a = Dec64((2<<(6*8+6) - 1) * 256)
	b = Dec64(1*256 + (256 - 128))
	ref = a
	testAdd(t, a, b, ref)
	// wrong 6629.509425000001 > 64842
	// -6629.509425000001+64842 is 58212.49057499999
	a = Dec64(-6629509425000001*256 + 244)
	b = Dec64(64842 * 256)
	ref = Dec64(1490239758719999989)
	testAdd(t, a, b, ref)

	// wrong
	a, _ = Parse("0.026000000000000002")
	b, _ = Parse("0.034999999999999996")
	ref, _ = Parse("0.06099999999999999")
	testAdd(t, a, b, ref)
}

func testBug(t *testing.T) {
	var a, b, ref Dec64
	// overflow
	a = Dec64((2 << (6*8 + 5)) * 256)
	b = Dec64((2 << (6*8 + 5)) * 256)
	ref = a
	testAdd(t, a, b, ref)
}

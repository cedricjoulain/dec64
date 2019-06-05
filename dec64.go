// DEC64 int as described in http://dec64.com/
// As exchanges always publish decimal values
// it's more accurate to store some sort of decimals
package dec64

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Dec64 int64
type ParseError error

const Empty = Dec64(0x0000000000000001)
const NotAvailable = Dec64(0x00000000000000ff)

// return a dec64 form string
func Parse(s string) (res Dec64, err error) {
	res = Empty
	if len(s) == 0 {
		return
	}
	start := 0
	neg := false
	if s[0] == '+' {
		// just forget
		start = 1
	} else {
		if s[0] == '-' {
			neg = true
			start = 1
		}
	}
	dot := false
	var exp, ncoef, coef, addExp, factor int64
	factor = 1
	expMode := false
	i := start
	for ; i < len(s); i++ {
		if s[i] == 'E' || s[i] == 'e' {
			expMode = true
			break
		}
		if s[i] == '0' {
			if coef == 0 && !dot {
				continue
			}
			factor *= 10
			if dot {
				exp--
			}
			addExp += 1
			continue
		}
		coef *= factor
		factor = 1
		addExp = 0
		if s[i] >= '1' && s[i] <= '9' {
			if neg {
				ncoef = 10*coef - int64(s[i]-'0')
				// check for overload
				if ncoef < (-1 * 0x7fffffffffffff) {
					// Can round up ?
					if s[i] >= '5' && coef > (-1*0x7fffffffffffff) {
						coef--
					}
					if dot {
						break
					}
					exp++
				} else {
					coef = ncoef
					if dot {
						exp--
					}
				}
			} else {
				ncoef = 10*coef + int64(s[i]-'0')
				// check for overload
				if ncoef > 0x7fffffffffffff {
					// Can round up ?
					if s[i] >= '5' && coef != 0x7fffffffffffff {
						coef++
					}
					if dot {
						break
					}
					exp++
				} else {
					coef = ncoef
					if dot {
						exp--
					}
				}
			}
			continue
		}
		if s[i] == '.' {
			if dot {
				err = ParseError(errors.New("Only one dot allowed"))
				return
			}
			dot = true
			continue
		}
		err = ParseError(fmt.Errorf("Unable to parse dec64 from %s", s))
		return
	}
	// was written like 1.5E-7
	if expMode {
		df := int64(1)
		toAdd := int64(0)
		for i++; i < len(s); i++ {
			switch s[i] {
			case '+':
				// Nothing to do...
			case '-':
				df = -1
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				toAdd = 10*toAdd + int64(s[i]-'0')
			default:
				err = ParseError(fmt.Errorf("Unable to handle %c in exponent", s[i]))
				return
			}
		}
		addExp += df * toAdd
	}
	exp += addExp
	// -128 is kept for special values
	if exp < -127 {
		err = ParseError(fmt.Errorf("%s is too small for dec64", s))
		return
	}
	if exp > 127 {
		err = ParseError(fmt.Errorf("%s is too big for dec64", s))
		return
	}

	res = Dec64((coef << 8) | (exp & 0xff))
	return
}

// TODO benchmark and use strings.Builder
func (d Dec64) String() string {
	// normalize to avoid 1.0000 for example
	d = Normalize(d)
	exp := int8(d)
	coef := int64(d) >> 8
	if coef == 0 {
		return "0"
	}
	minus := false
	if coef < 0 {
		minus = true
		coef *= -1
	}
	var bb bytes.Buffer
	dotLast := false
	for ; coef != 0; coef /= 10 {
		bb.WriteByte(uint8(coef%10) + '0')
		dotLast = false
		if exp < 0 {
			exp += 1
			if exp == 0 {
				bb.WriteByte('.')
				dotLast = true
			}
		}
	}
	if dotLast {
		bb.WriteByte('0')
	}
	// Smaller
	for ; exp < 0; exp++ {
		if exp == -1 {
			bb.Write([]byte{'0', '.', '0'})
		} else {
			bb.WriteByte('0')
		}
	}
	if minus {
		bb.WriteByte('-')
	}
	chr := bb.Bytes()
	// reverse buffer
	for i := 0; i < len(chr)/2; i++ {
		k := chr[i]
		chr[i] = chr[len(chr)-1-i]
		chr[len(chr)-1-i] = k
	}
	if exp > 0 {
		// Bigger
		var z strings.Builder
		z.Grow(int(exp) + len(chr))
		z.Write(chr)
		for i := int8(0); i < exp; i++ {
			z.WriteByte('0')
		}
		return z.String()
	} else {
		return string(chr)
	}
}

var (
	expf []float64
	expi []int64
)

func init() {
	expf = make([]float64, 256)
	expi = make([]int64, 256)
	f := 1.0
	for i := 0; i < 128; i++ {
		expf[i] = f
		expi[i] = int64(f)
		f *= 10
	}
	f = .1
	in := int64(10)
	for i := 255; i > 128; i-- {
		expf[i] = 1.0 / float64(in)
		expi[i] = in
		f /= 10
		in *= 10
	}
}

// Convert Dec64 to float64 using precomputed exponent
func Float64(d Dec64) (f float64) {
	if d&0xff > 127 {
		return float64(int64(d)>>8) / expf[256-d&0xff]
	} else {
		return float64(int64(d)>>8) * expf[d&0xff]
	}
}

func FromFloat64(f float64) (Dec64, error) {
	// TODO optimize !
	return Parse(strconv.FormatFloat(f, 'g', -1, 64))
}

// Convert int64 to Dec64
// i must be <= 0x00FFFFFF
// which makes around 281 474 976 711 000
func FromInt64(i int64) (Dec64, error) {
	// TODO > 0x00FFFFFF FFFFFFFF
	return Dec64(i * 256), nil
}

// Convert Dec64 to "normal" int64 keeping sign
func Int64(d Dec64) int64 {
	mant := int64(d) >> 8
	exp := int64(d) & 0xff
	if exp > 127 {
		return mant / expi[exp]
	} else {
		return mant * expi[exp]
	}
}

// Normalize Dec64 -> mantisse % 10 != 0
func Normalize(d Dec64) Dec64 {
	mant := int64(d) >> 8
	if mant == 0 {
		return 0
	}
	exp := int64(d) & 0xff
	for mant%10 == 0 {
		mant /= 10
		exp++
	}
	return Dec64(mant<<8 | (exp & 0xff))
}

// Compare 2 dec64, empty and not available are every thing
func (a *Dec64) Equal(b Dec64) bool {
	if *a == Empty || *a == NotAvailable {
		return true
	}
	if b == Empty || b == NotAvailable {
		return true
	}
	if *a != b {
		// Try to normalize to ensure they are different
		if Normalize(*a) != Normalize(b) {
			return false
		}
	}
	return true
}

// Ensure all exponents will be the same or closer as possible
// Dec64 will probably no longer be normalized
func Homogenize(values []Dec64) {
	// First find smaller exponent
	exp := int64(127)
	for _, d := range values {
		if d == 0 {
			// forget 0
			continue
		}
		e := int64(d) & 0xff
		if e > 127 {
			// negative
			e -= 256
		}
		if exp > e {
			exp = e
		}
	}
	// modify !
	for i, d := range values {
		if d == 0 {
			// forget 0
			continue
		}
		e := int64(d) & 0xff
		if e > 127 {
			// negative
			e -= 256
		}
		// we need to multiply mantisse by 10^exp
		e = e - exp
		if e == 0 {
			// nothing to do
			continue
		}
		old := int64(d) >> 8
		mant := (old * expi[e]) & 0xffffffffffffff
		// simple overflow check TODO optimize?
		if mant/expi[e] != old {
			// do nothing
			continue
		}
		// do not set to zero when undefined
		if mant != 0 || (exp&0xff) != 0 {
			values[i] = Dec64(mant<<8 | (exp & 0xff))
		}
	}
}

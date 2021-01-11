// Package dec64 int as described in http://dec64.com/
// As exchanges always publish decimal values
// it's more accurate to store some sort of decimals
package dec64

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Dec64 decimal number representation.
type Dec64 int64

// ParseError input has a wron format.
type ParseError error

// Empty no value encoding.
const Empty = Dec64(0x0000000000000001)

// NotAvailable missing encoding.
const NotAvailable = Dec64(0x00000000000000ff)

// Epsilon tolerance for comparaison with Float64.
const Epsilon = 3e-13

// Parse returns a dec64 form string.
func Parse(s string) (res Dec64, err error) {
	res = Empty
	if len(s) == 0 {
		return
	}
	start := 0
	neg := false
	for s[start] == ' ' {
		// trim starting space
		start++
	}
	if s[start] == '+' {
		// just forget
		start++
	} else {
		if s[start] == '-' {
			neg = true
			start++
		}
	}
	dot := false
	var (
		exp, ncoef, coef, addExp, factor int64
		overflow                         bool
	)
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
			if neg {
				overflow = (coef * factor * 10) < (-1 * 0x7fffffffffffff)
			} else {
				overflow = (coef * factor * 10) > 0x7fffffffffffff
			}
			if !overflow {
				factor *= 10
				if dot {
					exp--
				}
				addExp++
				continue
			}
		}
		coef *= factor
		factor = 1
		addExp = 0
		if s[i] >= '0' && s[i] <= '9' {
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
	// if early stop look if some exponent
	if !expMode {
		for ; i < len(s); i++ {
			if s[i] == 'E' || s[i] == 'e' {
				expMode = true
				break
			}
		}
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
			exp++
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
	}
	return string(chr)
}

// MarshalJSON Dec64 as a decimal
func (d *Dec64) MarshalJSON() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalJSON Dec64 as a decimal
func (d *Dec64) UnmarshalJSON(data []byte) (err error) {
	*d, err = Parse(string(data))
	return
}

var (
	// Expf list of 10^i as float64
	Expf []float64
	// Expi list of 10^i as int64
	Expi []int64
)

func init() {
	Expf = make([]float64, 256)
	Expi = make([]int64, 256)
	for i := 0; i < 128; i++ {
		f := math.Pow10(i)
		Expf[i] = f
		Expi[i] = int64(f)
	}
	in := int64(10)
	for i := 255; i > 128; i-- {
		Expf[i] = 1.0 / float64(in)
		Expi[i] = in
		in *= 10
	}
}

// Float64 converts Dec64 to float64 using precomputed exponent.
func Float64(d Dec64) (f float64) {
	if d&0xff > 127 {
		return float64(int64(d)>>8) / Expf[256-d&0xff]
	}
	return float64(int64(d)>>8) * Expf[d&0xff]
}

// FromFloat64 converts float64 to Dec64.
func FromFloat64(f float64) (Dec64, error) {
	// TODO optimize !
	return Parse(strconv.FormatFloat(f, 'g', -1, 64))
}

// FromInt64 converts int64 to Dec64
// i must be <= 0x00FFFFFF
// which makes around 281 474 976 711 000
func FromInt64(i int64) (Dec64, error) {
	// TODO > 0x00FFFFFF FFFFFFFF
	return Dec64(i * 256), nil
}

// Int64 converts Dec64 to "normal" int64 keeping sign
func Int64(d Dec64) int64 {
	mant := int64(d) >> 8
	exp := int64(d) & 0xff
	if exp > 127 {
		return mant / Expi[exp]
	}
	return mant * Expi[exp]
}

// Normalize Dec64 -> mantisse % 10 != 0
func Normalize(d Dec64) Dec64 {
	mant := int64(d) >> 8
	if mant == 0 {
		return 0
	}
	exp := int64(d) & 0xff
	if exp > 127 {
		exp -= 256
	}
	for mant%10 == 0 {
		mant /= 10
		exp++
	}
	return Dec64(mant<<8 | (exp & 0xff))
}

// Equal compares 2 dec64, empty and not available are every thing.
func (d *Dec64) Equal(b Dec64) bool {
	if *d == Empty || *d == NotAvailable {
		return true
	}
	if b == Empty || b == NotAvailable {
		return true
	}
	if *d != b {
		// Try to normalize to ensure they are different
		if Normalize(*d) != Normalize(b) {
			return false
		}
	}
	return true
}

// Homogenize ensures all exponents will be the same or closer as possible.
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
		if e <= 0 || e > 18 {
			// same or to huge difference, nothing to do
			continue
		}
		old := int64(d) >> 8
		mant := (old * Expi[e]) & 0xffffffffffffff
		// simple overflow check TODO optimize?
		if mant/Expi[e] != old {
			// do nothing
			continue
		}
		// do not set to zero when undefined
		if mant != 0 || (exp&0xff) != 0 {
			values[i] = Dec64(mant<<8 | (exp & 0xff))
		}
	}
}

// IsInt checks that d is an integer with no decimal parts.
func (d Dec64) IsInt() bool {
	// Normalize to ensure exponant is fully significativ
	e := int64(Normalize(d)) & 0xff
	if e > 127 {
		// negative, we have a decimal part
		return false
	}
	return true
}

// DEC64 int as described in http://dec64.com/
// As exchanges always publish decimal values
// it's more accurate to store some sort of decimals
package dec64

import (
	"errors"
	"fmt"
)

type Dec64 int64

const Empty = Dec64(0x0000000000000001)
const NotAvailable = Dec64(0x00000000000000ff)

// return a dec64 form string
func Parse(s string) (res Dec64, err error) {
	if len(s) == 0 {
		return
	}
	start := 0
	neg := false
	if s[0] == '-' {
		neg = true
		start = 1
	}
	dot := false
	var exp, coef, addExp, factor int64
	factor = 1
	for i := start; i < len(s); i++ {
		if s[i] == '0' {
			if coef == 0 {
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
				coef = 10*coef - int64(s[i]-'0')
			} else {
				coef = 10*coef + int64(s[i]-'0')
			}
			if dot {
				exp--
			}
			continue
		}
		if s[i] == '.' {
			if dot {
				err = errors.New("Only one dot allowed")
				return
			}
			dot = true
			continue
		}
		err = fmt.Errorf("Unable to parse dec64 from %s", s)
		return
	}
	exp += addExp
	// -128 is kept for special values
	if exp < -127 {
		err = fmt.Errorf("%s is to small for dec64")
		return
	}

	res = Dec64((coef << 8) | (exp & 0xff))
	return
}

func (d Dec64) String() (s string) {
	exp := int8(d)
	coef := int64(d) >> 8
	sign := ""
	if coef < 0 {
		sign = "-"
		coef *= -1
	}
	for ; coef != 0; coef /= 10 {
		s = string((coef%10)+'0') + s
		if exp < 0 {
			exp += 1
			if exp == 0 {
				s = "." + s
			}
		}
	}
	for ; exp > 0; exp-- {
		s += "0"
	}
	return sign + s
}

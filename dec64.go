// DEC64 int as described in http://dec64.com/
// As exchanges always publish decimal values
// it's more accurate to store some sort of decimals
package dec64

import (
	"errors"
	"fmt"
)

type Dec64 int64

// return a dec64 form string
func Parse(s string) (res Dec64, err error) {
	var d int64
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
	factor := int64(0)
	for i := start; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			if neg {
				d = d*10 - int64(s[i]-'0')
			} else {
				d = d*10 + int64(s[i]-'0')
			}
			if dot {
				factor--
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

	if factor < -127 {
		err = fmt.Errorf("%s is to small for dec64")
		return
	}

	d <<= 8
	d |= (factor & 0xff)
	res = Dec64(d)
	return
}

func (d Dec64) String() (s string) {
	return
}

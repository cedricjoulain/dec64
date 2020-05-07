// DEC64 int as described in http://dec64.com/
// As exchanges always publish decimal values
// it's more accurate to store some sort of decimals
package dec64

// returns 1 if a > 0, -1 if a < 0, 0 if a == 0
func Signum(d Dec64) int {
	if (uint64(d) & MMask) == 0 {
		return 0
	}
	if int64(d) < 0 {
		return -1
	}
	return 1
}

// Round to nearest, presicions is 10^n
func Round(d Dec64, n int64) Dec64 {
	mant := int64(d) >> 8
	if mant == 0 {
		return 0
	}
	e := int64(d) & 0xff
	if e > 127 {
		// negative
		e -= 256
	}
	// Normalize
	for mant%10 == 0 {
		mant /= 10
		e++
	}
	for e < n {
		if (n-e) > 1 || (mant%10 < 5 && mant%10 > -5) {
			mant /= 10
		} else {
			mant /= 10
			if mant > 0 {
				mant++
			} else {
				mant--
			}
		}
		e++
	}
	return Dec64(mant<<8 | (e & 0xff))
}

// Keep on mantisse
const (
	MMask     = 0xffffffffffffff00
	MOverflow = 0x0080000000000000
)

// Multiply Dec64 by an int64
func (d *Dec64) MultInt64(i int64) Dec64 {
	mant := (uint64(*d) & MMask) * uint64(i)
	return Dec64(int64(mant) | (int64(*d) & 0xff))
}

// Neg -> *-1
func (a Dec64) Neg() Dec64 {
	mant := uint64(a) & MMask
	return Dec64((-int64(mant)) | (int64(a) & 0xff))
}

// Add
func (a Dec64) Add(b Dec64) Dec64 {
	ea := int64(a) & 0xff
	eb := int64(b) & 0xff
	if ea == eb {
		// same exp, take care of overflow
		coef := int64(a)>>8 + int64(b)>>8
		// overflow ?
		if coef >= MOverflow || coef <= -MOverflow {
			coef /= 10
			ea++
		}
		return Dec64(coef<<8 | ea&0xff)
	} else {
		// different exp
		// first normalize
		na := Normalize(a)
		nb := Normalize(b)
		ea = int64(na) & 0xff
		if ea > 127 {
			// negative
			ea -= 256
		}
		eb = int64(nb) & 0xff
		if eb > 127 {
			// negative
			eb -= 256
		}
		if ea == eb {
			// same exp, take care of overflow
			coef := int64(na)>>8 + int64(nb)>>8
			// overflow ?
			if coef >= MOverflow || coef <= -MOverflow {
				coef /= 10
				ea++
			}
			return Dec64(coef<<8 | ea&0xff)
		}
		if ea > eb {
			// Switch to get ea > eb
			na, nb = nb, na
			ea, eb = eb, ea
		}
		var ncoef int64
		coefb := int64(nb) >> 8
		for (ea - eb) != 0 {
			ncoef = coefb * 10
			if (uint64(ncoef)^uint64(coefb))&0xff00000000000000 != 0 {
				// overflow on b!
				break
			}
			coefb = ncoef
			eb--
		}
		coefa := int64(na) >> 8
		// overflow loose precision on a
		if (eb - ea) != 0 {
			if (eb - ea) >= 128 {
				//a too small compared to b return b
				return nb
			}
			coefa /= expi[eb-ea]
			ea = eb
		}
		ncoef = coefa + coefb
		// overflow ?
		if ncoef >= MOverflow || ncoef <= -MOverflow {
			ncoef /= 10
			ea++
		}
		return Dec64(ncoef<<8 | ea&0xff)
	}
}

// Sub
func (a Dec64) Sub(b Dec64) Dec64 {
	return a.Add(b.Neg())
}

// Multiply
func (a Dec64) Mult(b Dec64) Dec64 {
	mant := (uint64(a) & MMask) * (uint64(b) & MMask)
	e := (int64(a) & 0xff) + (int64(b) & 0xff)
	return Dec64(int64(mant) | e&0xff)
}

// TODO
func (a Dec64) Div(b Dec64) (res Dec64) {
	// Trick...
	res, _ = FromFloat64(Float64(a) / Float64(b))
	return
}

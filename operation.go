// DEC64 int as described in http://dec64.com/
// As exchanges always publish decimal values
// it's more accurate to store some sort of decimals
package dec64

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
const mMask = 0xffffffffffffff00
// Multiply Dec64 by an int64
func (d *Dec64) MultInt64(i int64) Dec64 {
	mant := (uint64(*d) & mMask) * uint64(i)
	return Dec64(int64(mant) | (int64(*d) & 0xff))
}

// Add
func (a *Dec64) Add(b Dec64) Dec64 {
	mant := (uint64(*d) & mMask) * uint64(i)
	return Dec64(int64(mant) | (int64(*d) & 0xff))
}
// Sub
func (a *Dec64) Sub(b Dec64) Dec64 {
	mant := (uint64(*d) & mMask) * uint64(i)
	return Dec64(int64(mant) | (int64(*d) & 0xff))
}
// Multiply
func (a *Dec64) Mult(b Dec64) Dec64 {
	mant := (uint64(*a) & mMask) * (uint64(b) & mMask)
	e := (int64(*a) & 0xff) + (int64(b) & 0xff)
	return Dec64(mant | e  & 0xff)
}
// Div
func (a *Dec64) Div(b Dec64) Dec64 {
	mant := (uint64(*a) & mMask) * (uint64(*b) & mMask)

	return Dec64(int64(mant) | (int64(*d) & 0xff))
}

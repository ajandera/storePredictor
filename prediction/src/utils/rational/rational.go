package rational

import (
	"math"
	"math/big"
	"strconv"
	"strings"
)

// Rational stores a rational value.
type Rational struct {
	numerator   int64
	denominator int64
}

// New returns new rational number representation.
func New(n, d int64) Rational {
	return Rational{
		numerator:   n,
		denominator: d,
	}
}

// NewFromFloat returns new rational number representation retrieved from float64.
func NewFromFloat(f float64) (ev Rational, err error) {
	d, _ := math.Modf(f)

	fSl := strings.Split(strconv.FormatFloat(f, 'f', -1, 64), ".")
	fStr := "0"
	if len(fSl) == 2 {
		fStr = fSl[1]
	}

	var numerator int64
	numerator, err = strconv.ParseInt(fStr, 10, 64)
	if err != nil {
		return
	}
	denominator := int64(math.Pow(10, float64(len(fStr))))

	var negativeCoef int64 = 1
	if d < 0 {
		negativeCoef = -1
		d *= -1
	}

	ev = New(negativeCoef*(denominator*int64(d)+numerator), denominator)
	ev.Simplify()

	return
}

// Divide divides a rational value by the provided one.
func (ev Rational) Divide(e Rational) (nv Rational) {
	newNumerator := ev.numerator * e.denominator
	newDenominator := ev.denominator * e.numerator
	if newNumerator == newDenominator {
		nv = New(1, 1)
	} else {
		solveNegatives(&newNumerator, &newDenominator)
		nv = New(newNumerator, newDenominator)
		nv.Simplify()
	}
	return
}

// DivideByNum divides a rational value by the provided integer.
func (ev Rational) DivideByNum(i int64) Rational {
	return ev.Divide(New(i, 1))
}

// Multiply multiplies a rational value by provided one.
func (ev Rational) Multiply(e Rational) (nv Rational) {
	newNumerator := ev.numerator * e.numerator
	if newNumerator != 0 {
		newDenominator := ev.denominator * e.denominator
		if newNumerator == newDenominator {
			nv = New(1, 1)
		} else {
			solveNegatives(&newNumerator, &newDenominator)
			nv = New(newNumerator, newDenominator)
			nv.Simplify()
		}
	} else {
		nv = New(0, 1)
	}
	return
}

// MultiplyByNum multiplies a rational value by the provided integer.
func (ev Rational) MultiplyByNum(i int64) Rational {
	return ev.Multiply(New(i, 1))
}

// Add adds the provided rational value to an existing one.
func (ev Rational) Add(e Rational) (nv Rational) {
	newNumerator := ev.numerator*e.denominator + e.numerator*ev.denominator
	if newNumerator != 0 {
		newDenominator := ev.denominator * e.denominator
		if newNumerator == newDenominator {
			nv = New(1, 1)
		} else {
			solveNegatives(&newNumerator, &newDenominator)
			nv = New(newNumerator, newDenominator)
			nv.Simplify()
		}
	} else {
		nv = New(0, 1)
	}
	return
}

// AddNum adds the provided integer to an existing rational number.
func (ev Rational) AddNum(i int64) Rational {
	return ev.Add(New(i, 1))
}

// Subtract subtracts the provided rational value from an existing one.
func (ev Rational) Subtract(e Rational) (nv Rational) {
	newNumerator := ev.numerator*e.denominator - e.numerator*ev.denominator
	if newNumerator != 0 {
		newDenominator := ev.denominator * e.denominator
		if newNumerator == newDenominator {
			nv = New(1, 1)
		} else {
			solveNegatives(&newNumerator, &newDenominator)
			nv = New(newNumerator, newDenominator)
			nv.Simplify()
		}
	} else {
		nv = New(0, 1)
	}
	return
}

// SubtractNum subtracts the provided integer from an existing rational number.
func (ev Rational) SubtractNum(i int64) Rational {
	return ev.Subtract(New(i, 1))
}

// Simplify simplifies the rational number by dividing it's numerator and
// denominator by the GCD.
func (ev *Rational) Simplify() {
	if ev.denominator < 0 {
		ev.numerator *= -1
		ev.denominator *= -1
	}

	currentNumerator := ev.numerator
	currentDenominator := ev.denominator

	if currentNumerator < 0 && currentDenominator > 0 {
		currentNumerator *= -1
	} else if currentDenominator < 0 && currentNumerator >= 0 {
		currentDenominator *= -1
	}

	n := big.NewInt(currentNumerator)
	d := big.NewInt(currentDenominator)

	gcd := new(big.Int).GCD(nil, nil, n, d).Int64()

	if gcd > 1 {
		ev.numerator /= gcd
		ev.denominator /= gcd
	}
}

// IsNatural determines whether the rational number is also natural.
func (ev Rational) IsNatural() bool {
	if ev.numerator%ev.denominator == 0 {
		return true
	}
	return false
}

// Float64 returns the float64 representation of a rational number.
func (ev Rational) Float64() float64 {
	return float64(ev.numerator) / float64(ev.denominator)
}

// Get returns a value.
func (ev Rational) Get() (numerator, denominator int64) {
	return ev.numerator, ev.denominator
}

// GetNumerator returns a numerator.
func (ev Rational) GetNumerator() int64 {
	return ev.numerator
}

// GetDenominator returns a denominator.
func (ev Rational) GetDenominator() int64 {
	return ev.denominator
}

// GetModule returns rational number's module.
func (ev Rational) GetModule() Rational {
	solveNegatives(&ev.numerator, &ev.denominator)
	if ev.LessThanNum(0) {
		ev = ev.MultiplyByNum(-1)
	}
	return ev
}

// IsNull determines whether the value is zero.
func (ev Rational) IsNull() (n bool) {
	if ev.numerator == 0 {
		n = true
	}
	return
}

// RationalsAreNull determines whether the slice of Rationals contains only zero values.
func RationalsAreNull(l []Rational) (isNull bool) {
	isNull = true
	for _, v := range l {
		if !v.IsNull() {
			isNull = false
		}
	}
	return
}

func solveNegatives(n, d *int64) {
	if *d < 0 {
		*n *= -1
		*d *= -1
	}
}

package utils

import (
	"github.com/gonum/matrix/mat64"
	"math"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const eps = 1e-12

func Matrix(str string) *mat64.Dense {

	// remove [ and ]
	str = strings.Replace(str, "[", "", -1)
	str = strings.Replace(str, "]", "", -1)

	// calculate total number of rows
	parts := strings.SplitN(str, ";", -1)
	rows := len(parts)

	// calculate total number of columns
	colSlice := strings.Fields(parts[0])
	columns := len(colSlice)

	// replace all ; with space
	str = strings.Replace(str, ";", " ", -1)

	// convert str to slice
	elements := strings.Fields(str)

	// populate data for the new matrix(Dense type)
	data := make([]float64, rows*columns)
	for i := range data {
		floatValue, _ := strconv.ParseFloat(elements[i], 64)
		data[i] = floatValue
	}

	M := mat64.NewDense(rows, columns, data)
	return M
}

func MatrixByRow(a mat64.Matrix, rowNumber int) []float64 {
	dst := []float64{0.0, 0.0, 0.0}
	mat64.Row(dst, rowNumber, a)
	return dst
}

func SampleSet(cdf []float64) int {
	// https://stackoverflow.com/questions/50507513/golang-choice-number-from-slice-array-with-given-probability

	rand.Seed(time.Now().UnixNano())
	r := rand.Float64()

	bucket := 0

	for r > cdf[bucket] {
		bucket++
		if bucket == len(cdf) {
			bucket--
			break // need this to prevent runtime error: index out of range
		}
	}
	return bucket
}

func Heaviside(number float64) float64 {
	if number > 0 {
		return 1
	} else {
		return 0
	}
}

func LevinsonDurbin(d []float64, order int) ([]float64, float64) {
	coeff := AutoCorr(d, order)
	err := coeff[0]
	if math.Abs(err) < eps {
		err = 1.0 / eps
	}

	coeffTmp := make([]float64, order+1)

	i := 1
	for i <= order {
		k := coeff[i]
		for j := 1; j < i; j++ {
			fub := coeff[j] * coeff[i-j]
			k -= fub
		}
		k /= err
		//if math.Abs(k) > 1.0 {
		//	k = 1.0 / k
		//}
		coeffTmp[i] = k
		for j := 1; j < i; j++ {
			coeffTmp[j] -= k * coeff[i-j]
		}
		copy(coeff, coeffTmp[:i+1])
		err *= 1.0 - k*k
		i++
	}

	coeff = coeff[:i]
	for item := range coeff {
		if math.IsNaN(coeff[item]) {
			coeff[item] = 0
		}
		coeff[item] /= 20
	}
	return coeff, err
}

func AutoCorr(d []float64, order int) []float64 {
	//t := dsputils.NextPowerOf2(2*len(d)-1)
	rs := make([]float64, len(d))
	N := len(d) - order
	u := 0.0
	for i := 0; i < len(d); i++ {
		if i >= len(d) {
			u = 0.0
		} else {
			u = d[i]
		}
		for j := 0; j < len(rs); j++ {
			if i+j < len(d) {
				v := d[i+j]
				rs[j] += u * v
			}
		}
	}
	for i := range rs {
		rs[i] /= float64(N)
	}
	return rs
}

func Argmax(A []float64) int {
	var x int
	v := -1.0
	for i := 0; i < len(A); i++ {
		if A[i] > v {
			v = A[i]
			x = i
		}
	}
	return x
}

func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func Sum(array []int) int {
	result := 0
	for _, v := range array {
		result += v
	}
	return result
}

func GeneratePriceOfProduct(min float64, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return float64(rand.Intn(int(max)-int(min)+1) + int(min))
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func CalculateShift(vector []float64, order int) (int, int) {
	meanPercent := (sum(vector) / len(vector)) / 20
	forSort := make([]float64, len(vector))
	copy(forSort, vector)
	sort.Float64s(forSort)
	one := int(forSort[len(forSort)-1])
	two := int(forSort[len(forSort)-2])
	three := int(forSort[len(forSort)-3])
	four := int(forSort[len(forSort)-4])
	five := int(forSort[len(forSort)-5])

	var result []int
	var result2 []int
	var result3 []int
	var result4 []int
	var result5 []int

	for i, v := range vector {
		if int(v) > (one-meanPercent) && int(v) < (one+meanPercent) {
			result = append(result, i)
		}

		if int(v) > (two-meanPercent) && int(v) < (two+meanPercent) {
			result2 = append(result2, i)
		}

		if int(v) > (three-meanPercent) && int(v) < (three+meanPercent) {
			result3 = append(result3, i)
		}

		if int(v) > (four-meanPercent) && int(v) < (four+meanPercent) {
			result4 = append(result4, i)
		}

		if int(v) > (five-meanPercent) && int(v) < (five+meanPercent) {
			result5 = append(result5, i)
		}
	}

	sort.Ints(result)
	sort.Ints(result2)
	sort.Ints(result3)
	sort.Ints(result4)
	sort.Ints(result5)
	sum := 0
	divide := 0
	shift := 0
	longShift := 0
	if len(result) > 1 {
		sum += result[len(result)-1] + result[len(result)-2]
		divide += 2
	}
	if len(result2) > 1 {
		sum += result2[len(result2)-1] + result2[len(result2)-2]
		divide += 2
	}
	if len(result3) > 1 {
		sum += result3[len(result3)-1] + result3[len(result3)-2]
		divide += 2
	}
	if len(result4) > 1 {
		sum += result4[len(result4)-1] + result4[len(result4)-2]
		divide += 2
	}
	if len(result5) > 1 {
		sum += result5[len(result5)-1] + result5[len(result5)-2]
		divide += 2
	}

	if divide > 0 {
		shift = int(math.Ceil(float64(sum / divide)))
		if shift > 5 {
			shift = 4
		}
		longShift = shift * order
	}
	return shift, longShift
}

func sum(array []float64) int {
	result := 0.0
	for _, v := range array {
		result += v
	}
	return int(result)
}

func DFT_naive(input []float64) ([]float64, []float64) {
	realV := make([]float64, len(input))
	imagV := make([]float64, len(input))
	arg := -2.0 * math.Pi / float64(len(input))
	for k := 0; k < len(input); k++ {
		r, i := 0.0, 0.0
		for n := 0; n < len(input); n++ {
			r += input[n] * math.Cos(arg*float64(n)*float64(k))
			i += input[n] * math.Sin(arg*float64(n)*float64(k))
		}
		realV[k], imagV[k] = r*r, i
	}
	return realV, imagV
}

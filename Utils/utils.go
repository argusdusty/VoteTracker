package Utils

import (
	"math"
	"math/rand"
	"strings"
	"time"
)

const (
	LOG4          float64 = 2.0 * math.Ln2
	Ln3           float64 = 1.09861228866810969139524523692252570464749055782274945173469433363749429321860896687361575481373208878797
	SG_MAGICCONST float64 = 1 + Ln3 + Ln3 - math.Ln2
)

func init() {
	InitRand()
}

func InitRand() {
	rand.Seed(time.Now().UnixNano())
}

func RandGammaVariate(k float64) float64 {
	// k >= 0, mean is k*theta, variance is k*theta**2
	if k < 0.0 {
		panic("gammavariate: k must be >= 0.0")
	}
	if k == 0.0 {
		return 0.0
	}
	if k > 1.0 {
		// Uses R.C.H. Cheng, "The generation of Gamma
		// variables with non-integral shape parameters",
		// Applied Statistics, (1977), 26, No. 1, p71-74
		var ainv, bbb, ccc, u1, u2, v, x, z, r float64
		ainv = math.Sqrt(2.0*k - 1.0)
		bbb = k - LOG4
		ccc = k + ainv
		for true {
			u1 = rand.Float64()
			if !(1e-7 < u1 && u1 < 0.9999999) {
				continue
			}
			u2 = 1.0 - rand.Float64()
			v = math.Log(u1/(1.0-u1)) / ainv
			x = k * math.Exp(v)
			z = u1 * u1 * u2
			r = bbb + ccc*v - x
			if r+SG_MAGICCONST-4.5*z >= 0.0 || r >= math.Log(z) {
				return x
			}
		}
		return x
	} else if k < 1.0 {
		// k is between 0 and 1 (exclusive)
		// Uses ALGORITHM GS of Statistical Computing - Kennedy & Gentle
		var u, b, p, x, u1 float64
		for true {
			u = rand.Float64()
			b = (math.E + k) / math.E
			p = b * u
			if p <= 1.0 {
				x = math.Pow(p, 1.0/k)
			} else {
				x = -math.Log((b - p) / k)
			}
			u1 = rand.Float64()
			if p > 1.0 {
				if u1 <= math.Pow(x, k-1.0) {
					break
				}
			} else if u1 <= math.Exp(-x) {
				break
			}
		}
		return x
	} else {
		return rand.ExpFloat64()
	}
}

func FixName(c string) string {
	vals := strings.Split(c, " ")
	first := vals[0]
	lastidx := len(vals) - 1
	for (strings.HasPrefix(vals[lastidx], "(") && strings.HasSuffix(vals[lastidx], ")")) || vals[lastidx] == "Jr" || vals[lastidx] == "Jr." || vals[lastidx] == "Sr" || vals[lastidx] == "Sr." || vals[lastidx] == "I" || vals[lastidx] == "II" || vals[lastidx] == "III" || vals[lastidx] == "" {
		lastidx--
	}
	last := vals[lastidx]
	last = strings.TrimSuffix(last, ",")
	return strings.Title(strings.ToLower(first + " " + last))
}

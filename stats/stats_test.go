package stats

import (
	"math"
	"testing"
)

// Used to test automatically Lfit
func lfit_innertest(ndata int, coefs2 []float64, t *testing.T) {
	// data
	x := make([]float64, ndata)
	y := make([]float64, ndata)
	devs := make([]float64, ndata)
	for i := range x {
		x[i] = float64(i)
		y[i] = Fpoly(x[i], coefs2)
		devs[i] = 1.0
	}
	// Interpolation
	l := len(coefs2)
	ia := make([]bool, l)
	cofs := make([]float64, l)
	for i := 0; i < l; i++ {
		ia[i] = true
		cofs[i] = 1.
	}
	if _, _, err := Lfit(x, y, devs, cofs, ia, l, Fcoefs); err != nil {
		panic(err)
	}
	// Verification
	eps := 0.1
	for i := range coefs2 {
		if math.Abs(cofs[i]-coefs2[i]) > eps {
			t.Errorf("Bad coef: wanted: %f found: %f", coefs2[i], cofs[i])
		}
	}
}

// Test the interpolation of a trinome
func TestLfitDegree2PositiveCoefs(t *testing.T) {
	lfit_innertest(10, []float64{1., 2., 1.}, t)
}

func TestLfitDegree2NegativeCoefs(t *testing.T) {
	lfit_innertest(10, []float64{-5., -2., -3.}, t)
}

func TestLfitDegree2With100Points(t *testing.T) {
	lfit_innertest(100, []float64{1., 2., 1.}, t)
}

// Test the interpolation of a polynom of degree 3
func TestLfitDegree3(t *testing.T) {
	lfit_innertest(100, []float64{1., 3., 3, 1.}, t)
}

// Test LSFitLinear
// Coefs2 reads y = coefs2[0] + coefs2[1] * x
func lsFitLinear_innertest(ndata int, coefs2 []float64, t *testing.T) {
	// data
	x := make([]float64, ndata)
	y := make([]float64, ndata)
	for i := range x {
		x[i] = float64(i)
		y[i] = Fpoly(x[i], coefs2)
	}
	// Intepolation
	a, b, siga, sigb, chi2, sigdat := LSFitLinear(x, y)
	// Verification
	eps := 0.1
	if math.Abs(a-coefs2[0]) > eps {
		t.Errorf("Bad slope a : wanted: %f found: %f", coefs2[0], a)
	}
	if math.Abs(b-coefs2[1]) > eps {
		t.Errorf("Bad origin b : wanted: %f found: %f", coefs2[1], b)
	}
	if math.Abs(siga) > eps {
		t.Errorf("Too big siga : wanted: %f found: %f", eps, siga)
	}
	if math.Abs(sigb) > eps {
		t.Errorf("Too big sigb : wanted: %f found: %f", eps, sigb)
	}
	if math.Abs(chi2) > eps {
		t.Errorf("Too big chi2 : wanted: %f found: %f", eps, chi2)
	}
	if math.Abs(sigdat) > eps {
		t.Errorf("Too big sigdat : wanted: %f found: %f", eps, sigdat)
	}
}

func TestLsFitLinear_1(t *testing.T) {
	lsFitLinear_innertest(100, []float64{1., 3.}, t)
}

func TestLsFitLinear_Hor(t *testing.T) {
	lsFitLinear_innertest(100, []float64{1., 0.}, t)
}

func TestLsFitLinear_Vert(t *testing.T) {
	lsFitLinear_innertest(100, []float64{0., 1.}, t)
}

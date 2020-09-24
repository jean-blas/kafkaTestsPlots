package stats

import (
	"errors"
	"math"
	"plots/sliceutil"
)

// Fitting routine for polynomial of degree len(a)-1
// return the coefficients
// To be used with Lfit when fitting a polynom
func Fcoefs(x float64, a []float64) []float64 {
	res := make([]float64, len(a))
	res[0] = 1.
	for i := 1; i < len(a); i++ {
		res[i] = res[i-1] * x
	}
	return res
}

// The corresponding polynome of degree len(coefs)-1
// To be used with Lfit when fitting a polynom
func Fpoly(x float64, coefs []float64) float64 {
	res := 0.
	xx := 1.
	for _, c := range coefs {
		res += c * xx
		xx *= x
	}
	return res
}

// Given a set of data points x[0..ndat-1], y [0..ndat-1] with individual standard deviations sig[0..ndat-1],
// use chi2 minimization to fit for some or all of the coefficients a[0..ma-1] of a function that depends linearly on a,
// y = sum{i} a_i afunc_i(x)
// The input array ia[0..ma-1] indicates those components of a that should be fitted for (others are held fixed)
// the program returns values for a[0..ma-1], chi2 and the covariance matrix covar[0..ma-1][0..ma-1]
// The user supplies a routine funcs(x, afunc) that returns the ma basis functions evaluated at x in the array afunc[0..ma-1]
func Lfit(x, y, sig, a []float64, ia []bool, ma int, funcs func(float64, []float64) []float64) (float64, [][]float64, error) {
	mfit := 0
	for _, iia := range ia {
		if iia {
			mfit++
		}
	}
	if mfit == 0 {
		return 0., nil, errors.New("lfit: no parameters to be fitted")
	}
	beta := make([][]float64, ma)
	for i := range beta {
		beta[i] = make([]float64, 1)
	}
	covar := make([][]float64, ma)
	for i := range covar {
		covar[i] = make([]float64, ma)
	}
	afunc := make([]float64, ma)
	copy(afunc, a)
	ndat := len(x)
	for i := 0; i < ndat; i++ {
		afunc = funcs(x[i], afunc)
		ym := y[i]
		if mfit < ma {
			for j := 0; j < ma; j++ {
				if !ia[j] {
					ym -= a[j] * afunc[j]
				}
			}
		}
		sig2i := 1. / (sig[i] * sig[i])
		j := -1
		for l := 0; l < ma; l++ {
			if ia[l] {
				wt := afunc[l] * sig2i
				j++
				k := -1
				for m := 0; m <= l; m++ {
					if ia[m] {
						k++
						covar[j][k] += wt * afunc[m]
					}
				}
				beta[j][0] += ym * wt
			}
		}
	}
	for j := 1; j < mfit; j++ {
		for k := 0; k < j; k++ {
			covar[k][j] = covar[j][k]
		}
	}
	if err := gaussj(covar, mfit, beta); err != nil {
		return 0., nil, err
	}
	j := -1
	for l := 0; l < ma; l++ {
		if ia[l] {
			j++
			a[l] = beta[j][0]
		}
	}
	chisq := 0.
	for i := 0; i < ndat; i++ {
		afunc = funcs(x[i], afunc)
		sum := 0.
		for j := 0; j < ma; j++ {
			sum += a[j] * afunc[j]
		}
		chisq += (y[i] - sum) / sig[i] * (y[i] - sum) / sig[i]
	}
	covsrt(covar, ma, ia, mfit)
	return chisq, covar, nil
}

// Expand in storage the covariance matrix covar, so as to take into account parameters that are being held fixed.
func covsrt(covar [][]float64, ma int, ia []bool, mfit int) {
	for i := mfit; i < ma; i++ {
		for j := 0; j < i; j++ {
			covar[i][j] = 0.
			covar[j][i] = 0.
		}
	}
	k := mfit - 1
	for j := ma - 1; j >= 0; j-- {
		if ia[j] {
			for i := 0; i < ma; i++ {
				covar[i][k], covar[i][j] = covar[i][j], covar[i][k]
			}
			for i := 0; i < ma; i++ {
				covar[k][i], covar[j][i] = covar[j][i], covar[k][i]
			}
			k--
		}
	}
}

// Linear equation solution by Gauss-Jordan elimination
// a[1..n][1..n] is the input matrix
// b[1..n][1..m] is input containing the m right-hand side vectors
// On output,
//	a is replaced by its matrix inverse,
//	b is replaced by the corresponding set of solution vectors
func gaussj(a [][]float64, n int, b [][]float64) error {
	var icol, irow int
	indxc := make([]int, n)
	indxr := make([]int, n)
	ipiv := make([]int, n)
	m := len(b[0])

	for i := 0; i < n; i++ {
		big := 0.
		for j := 0; j < n; j++ {
			if ipiv[j] != 1 {
				for k := 0; k < n; k++ {
					if ipiv[k] == 0 {
						if math.Abs(a[j][k]) >= big {
							big = math.Abs(a[j][k])
							irow = j
							icol = k
						}
					} else if ipiv[k] > 1 {
						return errors.New("gaussj : singular matrix - 1")
					}
				}
			}
		}
		(ipiv[icol])++
		if irow != icol {
			for l := 0; l < n; l++ {
				a[irow][l], a[icol][l] = a[icol][l], a[irow][l]
			}
			for l := 0; l < m; l++ {
				b[irow][l], b[icol][l] = b[icol][l], b[irow][l]
			}
		}
		indxr[i] = irow
		indxc[i] = icol
		if a[icol][icol] == 0 {
			return errors.New("gausj : singular matrix - 2")
		}
		pivinv := 1. / a[icol][icol]
		a[icol][icol] = 1.
		for l := 0; l < n; l++ {
			a[icol][l] *= pivinv
		}
		for l := 0; l < m; l++ {
			b[icol][l] *= pivinv
		}
		for ll := 0; ll < n; ll++ {
			if ll != icol {
				dum := a[ll][icol]
				a[ll][icol] = 0.
				for l := 0; l < n; l++ {
					a[ll][l] -= a[icol][l] * dum
				}
				for l := 0; l < m; l++ {
					b[ll][l] -= b[icol][l] * dum
				}
			}
		}
	}
	for l := n - 1; l >= 0; l-- {
		if indxr[l] != indxc[l] {
			for k := 0; k < n; k++ {
				a[k][indxr[l]], a[k][indxc[l]] = a[k][indxc[l]], a[k][indxr[l]]
			}
		}
	}
	return nil
}

// Given arrays x[0..n-1] and y[0..n-1] containing a tabulated function yi = f(xi),
// returns an array of coefficients cof[0..n-1], such that yi = sum[0,n-1]{cof_j xi^j}
func PolyCoefs(x, y []float64) []float64 {
	n := len(x)
	s := make([]float64, n)
	cof := make([]float64, n)
	s[n-1] = -x[0]
	for i := 1; i < n; i++ {
		for j := n - 1 - i; j < n-1; j++ {
			s[j] -= x[i] * s[j+1]
		}
		s[n-1] -= x[i]
	}
	for j := 0; j < n; j++ {
		phi := float64(n)
		for k := n - 1; k > 0; k-- {
			phi = float64(k)*s[k] + x[j]*phi
		}
		ff := y[j] / phi
		b := 1.0
		for k := n - 1; k >= 0; k-- {
			cof[k] += b * ff
			b = s[k] + x[j]*b
		}
	}
	return cof
}

// Compute the mean of the data
func Mean(data []int64) float64 {
	if len(data) == 0 {
		return 0
	}
	var m int64 = 0
	var d int64
	for _, d = range data {
		m += d
	}
	return float64(m) / float64(len(data))
}

// Given a set of data points x[0..ndata-1], y[0..ndata-1],
// sets a, b and their uncertainties siga and sigb, and the chi-square chi2
// Note : fit line reads y = b * x + a
func LSFitLinear(x, y []float64) (float64, float64, float64, float64, float64, float64) {
	var i int
	var ss, sxoss float64
	ndata := len(x)
	sx := 0.
	sy := 0.
	st2 := 0.
	chi2 := 0.
	sigdat := 0.
	b := 0.0
	for i = 0; i < ndata; i++ {
		sx += x[i]
		sy += y[i]
	}
	ss = float64(ndata)
	sxoss = sx / ss
	for i = 0; i < ndata; i++ {
		t := x[i] - sxoss
		st2 += t * t
		b += t * y[i]
	}
	b /= st2
	a := (sy - sx*b) / ss
	siga := math.Sqrt((1.0 + sx*sx/(ss*st2)) / ss)
	sigb := math.Sqrt(1.0 / st2)
	for i = 0; i < ndata; i++ {
		t := y[i] - a - b*x[i]
		chi2 += t * t
	}
	if ndata > 2 {
		sigdat = math.Sqrt(chi2 / float64(ndata-2))
	}
	siga *= sigdat
	sigb *= sigdat
	return a, b, siga, sigb, chi2, sigdat
}

// Given an array of data[0..n-1], this routine returns its:
// mean ave
// average deviation adev
// standard deviation sdev (aka sqrt of variance)
// skewness skew
// kurtosis curt
func Moments(data []float64) (float64, float64, float64, float64, float64, error) {
	n := len(data)
	if n <= 1 {
		return 0, 0, 0, 0, 0, errors.New("Moments: n must be at least 2")
	}
	s := 0.0 //First pass to get the mean.
	for j := 0; j < n; j++ {
		s += data[j]
	}
	an := float64(n)
	ave := s / an
	ep := 0.0
	adev := 0.0
	var2 := 0.0
	skew := 0.0
	curt := 0.0
	var p float64
	for j := 0; j < n; j++ {
		s = data[j] - ave
		adev += math.Abs(s)
		ep += s
		p = s * s
		var2 += (p)
		p *= s
		skew += (p)
		p *= s
		curt += (p)
	}
	adev /= an
	var2 = (var2 - ep*ep/an) / float64(n-1) // Corrected two-pass formula.
	sdev := math.Sqrt(var2)
	if var2 != 0.0 {
		skew /= (an * var2 * sdev)
		curt = curt/(an*var2*var2) - 3.0
	} else {
		return 0, 0, 0, 0, 0, errors.New("Moments: no skew/kurtosis when variance = 0")
	}
	return ave, adev, sdev, skew, curt, nil
}

// ToHisto Transform a slice of data into an histogram of size ncol
// and normalize it
// Returns the histo together with its xmin and xmax values
func ToHisto(data []float64, ncol int) (float64, float64, []float64) {
	min, max := sliceutil.MinMax(data)
	acol := float64(ncol)
	itv := (max - min) / acol        // interval (size of a histo column)
	histo := make([]float64, ncol+1) // ncol + 1 to add max into the last column
	for _, d := range data {
		col := int((d - min) / itv)
		histo[col]++
	}
	// Normalize to 1
	sum := 0.
	for _, h := range histo {
		sum += h
	}
	sum *= itv
	for i := range histo {
		histo[i] /= sum
	}
	return min, max, histo
}

// Gauss returns the normale function at x
// mean : average
// sdev : standard deviation
func Gauss(x, mean, sdev float64) float64 {
	return math.Exp(-(x-mean)*(x-mean)/(2.*sdev*sdev)) / (sdev * math.Sqrt(2.*math.Pi))
}

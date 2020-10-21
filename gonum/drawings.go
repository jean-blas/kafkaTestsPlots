package main

import (
	"fmt"
	"image/color"
	"math"
	"path/filepath"
	"plots/parser"
	"plots/plotfunc"
	"plots/sliceutil"
	"plots/stats"

	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var ComparePNGsuffix string // A suffixe to be added to the PNG when comparing configs

// Print the values of x and y to screen
func print(x, y []float64, label string) {
	if !PRINT {
		return
	}
	fmt.Println(label)
	for i := range x {
		fmt.Println("\t", x[i], y[i])
	}
}

// Process the comparison of the specified configs
func compareConfigs(confs []Config) error {
	cfgs := make([]Config, len(confs))
	for i, c := range confs {
		c.prepare()
		cfgs[i] = c
	}
	if err := compareThroughputs(cfgs); err != nil {
		return err
	}
	if err := compareNbMsgPerSec(cfgs); err != nil {
		return err
	}
	if err := compareMeansErr(cfgs); err != nil {
		return err
	}
	if err := compareMeansLine(cfgs); err != nil {
		return err
	}
	return nil
}

// used to pass the func as first citizen
type fdraw func(string, int) error

// Draw the function for one or all files, according to the value of "n"
func drawCFiles(c Config, n int, f fdraw) {
	if n < 0 {
		for _, file := range c.files {
			if err := f(file, c.nbPtsDiscard); err != nil {
				panic(err)
			}
		}
	} else {
		if err := f(c.files[n], c.nbPtsDiscard); err != nil {
			panic(err)
		}
	}
}

// Draw a single config "c" according to the Draws enum "d" value
// "n" is the number of the config sample file (-1 = draw all files of the config)
func drawConfig(c Config, d Draws, n int) {
	c.prepare()

	if d == Dall || d == Dfile {
		drawCFiles(c, n, drawFile)
	}
	if d == Dall || d == DslideFile {
		drawCFiles(c, n, drawSlideFile)
	}
	if d == Dall || d == DhistoFile {
		drawCFiles(c, n, drawHistoFile)
	}
	if d == Dall || d == DmeansFile {
		if err := drawMeansFiles(c); err != nil {
			panic(err)
		}
	}
	if d == Dall || d == DmeansErrFiles {
		if err := drawMeansErrFiles(c); err != nil {
			panic(err)
		}
	}
	if d == Dall || d == Dthroughput {
		if err := drawThroughputsFiles(c); err != nil {
			panic(err)
		}
	}
	if d == Dall || d == DnbMsgPerSec {
		if err := drawNbMsgPerSecFiles(c); err != nil {
			panic(err)
		}
	}
}

// Compute the number of messages per seconds for each file
func computeNbMsgPerSecFiles(files []string, sizes []float64, nbPtsDiscard, ndata int, abscisIsNb bool) ([]float64, error) {
	NB_MSG := float64(ndata - nbPtsDiscard) // default number of messages sent
	trput := make([]float64, len(files))
	for i, f := range files {
		ts1, ts2, err := parser.ParseData(f)
		if err != nil {
			return nil, err
		}
		seconds := float64(ts2[len(ts2)-1]-ts1[nbPtsDiscard]) / 1.e9
		if abscisIsNb {
			NB_MSG = sizes[i]
		}
		trput[i] = NB_MSG / seconds
	}
	return trput, nil
}

// Comparison of number of messages per seconds for different configs
func compareNbMsgPerSec(confs []Config) error {
	// Create the plot
	p, err := plotfunc.NewPlot("Msg / s", confs[0].xlabel, "nb of Msg / s")
	if err != nil {
		return err
	}
	for i, c := range confs {
		trput, err := computeNbMsgPerSecFiles(c.files, c.abscis, c.nbPtsDiscard, c.ndata, c.abscisIsNb)
		if err != nil {
			return err
		}
		print(c.abscis, trput, "NbMsgPerSec")
		if err = plotfunc.AddWithLineXY(c.abscis, trput, c.legend(), i, p); err != nil {
			return err
		}
	}
	// Save the plot to a PNG file.
	return p.Save(10*vg.Centimeter, 10*vg.Centimeter, confs[0].xlabel+"_nbMsgPerSec_"+ComparePNGsuffix+".png")
}

// Compute the number of messages per second for every dataset and draw it
// files : files to parse
// sizes : files corresponding abcissa
func drawNbMsgPerSecFiles(c Config) error {
	trput, err := computeNbMsgPerSecFiles(c.files, c.abscis, c.nbPtsDiscard, c.ndata, c.abscisIsNb)
	if err != nil {
		return err
	}
	base := filepath.Base(c.root)
	return drawPointsXY(c.abscis, trput, c.xlabel, "nb of msg / s", base+c.title, base+"_nbmespersec.png")
}

// Compute the throughput for each file
func computeThoughputFiles(files []string, sizes []float64, mb float64, ndata, nbPtsDiscard int,
	abscisIsNb, abscisIsSz bool) ([]float64, error) {
	SIZE_MSG := mb                          // default message size in Mb
	NB_MSG := float64(ndata - nbPtsDiscard) // default number of messages sent
	trput := make([]float64, len(files))
	for i, f := range files {
		ts1, ts2, err := parser.ParseData(f)
		if err != nil {
			return nil, err
		}
		seconds := float64(ts2[len(ts2)-1]-ts1[nbPtsDiscard]) / 1.e9
		if abscisIsNb {
			NB_MSG = sizes[i]
		}
		if abscisIsSz {
			SIZE_MSG = sizes[i] / 1000. // size en Mb
		}

		trput[i] = NB_MSG * SIZE_MSG / seconds
	}
	return trput, nil
}

// Comparison of throughputs for different configs
func compareThroughputs(confs []Config) error {
	// Create the plot
	p, err := plotfunc.NewPlot("Throughputs", confs[0].xlabel, "nb of Mb / s")
	if err != nil {
		return err
	}
	for i, c := range confs {
		trput, err := computeThoughputFiles(c.files, c.abscis, c.mb, c.ndata, c.nbPtsDiscard, c.abscisIsNb, c.abscisIsSz)
		if err != nil {
			return err
		}
		print(c.abscis, trput, "Throughput")
		if err = plotfunc.AddWithLineXY(c.abscis, trput, c.legend(), i, p); err != nil {
			return err
		}
	}
	// Save the plot to a PNG file.
	return p.Save(10*vg.Centimeter, 10*vg.Centimeter, confs[0].xlabel+"_throughputs_"+ComparePNGsuffix+".png")
}

// Compute the throughput for every dataset and draw it
// files : files to parse
// sizes : files corresponding abcissa
// func drawThroughputsFiles(files []string, sizes []float64, nbPtsDiscard int, xlabel, title, outPng string, abscisIsNb,
// abscisIsSz bool, mb float64, ndata int) error {
func drawThroughputsFiles(c Config) error {
	trput, err := computeThoughputFiles(c.files, c.abscis, c.mb, c.ndata, c.nbPtsDiscard, c.abscisIsNb, c.abscisIsSz)
	if err != nil {
		return err
	}
	print(c.abscis, trput, "Throughput")
	base := filepath.Base(c.root)
	return drawPointsXY(c.abscis, trput, c.xlabel, "nb of Mb / s", base+c.title, base+"_throughputs.png")
}

// call ParseFile (with the given filename)
// call slide (with "nval", the number of samples to slide)
func drawSlideFile(filename string, nbPtsDiscard int) error {
	fvalues, err := parseFile(filename)
	if err != nil {
		return err
	}
	base := filepath.Base(filename)
	outPng := fmt.Sprintf("%s_nval%d_slide.png", base, NVAL)
	title := fmt.Sprintf("%s\n(nval=%d)", base, NVAL)
	return drawSlide(fvalues, NVAL, nbPtsDiscard, "Msg number", "times (ms)", title, outPng)
}

// slide the data with an interval of nval data values
// for each window, compute the mean and dev
// draw and save into a PNG image
func drawSlide(data []float64, nval int, nbPtsDiscard int, xlabel, ylabel, title, outPng string) error {
	var means, devs, x []float64
	temp := make([]float64, nval)
	for i := 0; i < len(data)-nval; i++ {
		for j := 0; j < nval; j++ {
			temp[j] = data[i+j]
		}
		mean, _, sdev, _, _, err := stats.Moments(temp)
		if err != nil {
			return err
		}
		means = append(means, mean)
		devs = append(devs, sdev)
		x = append(x, float64(i))
	}
	return drawLinearFit(x[nbPtsDiscard:], means[nbPtsDiscard:], xlabel, ylabel, title, outPng)
}

// Parse the filename in the root folder
// transform the data into millis
func parseFile(filename string) ([]float64, error) {
	values, err := parser.ParseAndDiff(filename)
	if err != nil {
		return nil, err
	}
	// Transform nano into milli
	fvalues := sliceutil.MapF64(sliceutil.I64ToF64(values), func(x float64) float64 { return x / 1000000. })
	return fvalues, nil
}

// call parseFile and drawHisto
// image name = ${filename}_histo.png
func drawHistoFile(filename string, nbPtsDiscard int) error {
	fvalues, err := parseFile(filename)
	if err != nil {
		return err
	}
	return drawHisto(fvalues, filepath.Base(filename), filepath.Base(filename)+"_histo.png", nbPtsDiscard)
}

// Draw a normalized histogram
// compare with its gaussian function
// save the plot to PNG image file (name is filename_histo.png)
func drawHisto(data []float64, title, outPng string, nbPtsDiscard int) error {
	// Compute the moments
	mean, adev, sdev, skew, curt, err := stats.Moments(data[nbPtsDiscard:])
	if PRINT {
		fmt.Printf("Moments : mean=%.3e adev=%.3e sdev=%.3e skew=%.3e curt=%.3e %s\n", mean, adev, sdev, skew, curt, title)
	}
	// Create the plot
	p, err := plotfunc.NewPlot(title, "Latency (ms)", "Nb of values")
	if err != nil {
		return err
	}
	// Draw an histogram
	v := make(plotter.Values, len(data[nbPtsDiscard:]))
	for i := range v {
		v[i] = data[i+nbPtsDiscard]
	}
	h, err := plotter.NewHist(v, NCOL)
	if err != nil {
		return err
	}
	h.Normalize(1)
	p.Add(h)
	// Add the normal distribution function
	plotfunc.AddGaussian(mean, sdev, p)
	// Save the plot to a PNG file.
	return p.Save(15*vg.Centimeter, 15*vg.Centimeter, outPng)
}

// Comparison of means with deviations for different configs
func compareMeansErr(confs []Config) error {
	// Create the plot
	p, err := plotfunc.NewPlot("Means", confs[0].xlabel, "times (ms)")
	if err != nil {
		return err
	}
	for i, c := range confs {
		means, devs, err := computeMeansErrFiles(c.files, c.abscis, c.nbPtsDiscard)
		if err != nil {
			return err
		}
		print(c.abscis, means, "MeansErr")
		if err = plotfunc.AddWithErrXY(c.abscis, means, devs, c.legend(), i, p); err != nil {
			return err
		}
	}
	// Save the plot to a PNG file.
	return p.Save(10*vg.Centimeter, 10*vg.Centimeter, confs[0].xlabel+"_meansErr_"+ComparePNGsuffix+".png")
}

// Comparison of means for different configs
func compareMeansLine(confs []Config) error {
	// Create the plot
	p, err := plotfunc.NewPlot("Means", confs[0].xlabel, "times (ms)")
	if err != nil {
		return err
	}
	for i, c := range confs {
		means, _, err := computeMeansErrFiles(c.files, c.abscis, c.nbPtsDiscard)
		if err != nil {
			return err
		}
		print(c.abscis, means, "Means")
		if err = plotfunc.AddWithLineXY(c.abscis, means, c.legend(), i, p); err != nil {
			return err
		}
	}
	// Save the plot to a PNG file.
	return p.Save(10*vg.Centimeter, 10*vg.Centimeter, confs[0].xlabel+"_means_"+ComparePNGsuffix+".png")
}

// Compute the means and deviations for each file
func computeMeansErrFiles(files []string, sizes []float64, nbPtsDiscard int) ([]float64, []float64, error) {
	means := make([]float64, len(files))
	devs := make([]float64, len(files))
	for i, f := range files {
		fvalues, err := parseFile(f)
		if err != nil {
			return nil, nil, err
		}
		mean, adev, sdev, skew, curt, err := stats.Moments(fvalues[nbPtsDiscard:])
		if err != nil {
			return nil, nil, err
		}
		if PRINT {
			fmt.Printf("Moments : mean=%.3e adev=%.3e sdev=%.3e skew=%.3e curt=%.3e %s\n", mean, adev, sdev, skew, curt, filepath.Base(f))
		}
		means[i] = mean
		devs[i] = sdev / math.Sqrt(float64(len(fvalues)))
	}
	return means, devs, nil
}

// Parse each file of suffixes
// compute the means and draw it with the error bars
// save the plot to a PNG file
func drawMeansErrFiles(c Config) error {
	base := filepath.Base(c.root)
	means, devs, err := computeMeansErrFiles(c.files, c.abscis, c.nbPtsDiscard)
	if err != nil {
		return err
	}
	return drawErrsXY(c.abscis, means, devs, c.xlabel, "times (ms)", base+c.title, base+"_mean_err.png")
}

// Draw the x,y  for every dataset with deviation as Y error bars
func drawErrsXY(x, y, devs []float64, xlabel, ylabel, title, outPng string) error {
	// Create the plot
	p, err := plotfunc.NewPlot(title, xlabel, ylabel)
	if err != nil {
		return err
	}
	// Add the means with errors
	plotfunc.AddWithErrXY(x, y, devs, "", 0, p)

	// a, b, siga, sigb, chi2, sigdat := plotfunc.AddLinearfit(x[1:], y[1:], p)
	// if PRINT {
	// 	fmt.Println("drawMeansErr moments", a, b, siga, sigb, chi2, sigdat)
	// }

	// Add regression trinome
	// coefs, err := plotfunc.AddPolyfit(x, y, devs, 3, p)
	// if err != nil {
	// 	return err
	// }
	// if PRINT {
	// 	fmt.Println("drawMeansErr poly coefs", coefs)
	// }
	// Save the plot to a PNG file.
	return p.Save(10*vg.Centimeter, 10*vg.Centimeter, outPng)
}

// Compute the means for every dataset and draw it
// Compute the linear regression that fits the means and draw it
func drawMeansFiles(c Config) error {
	var means []float64
	// Parse the files and compute the means
	for _, f := range c.files {
		values, err := parseFile(f)
		if err != nil {
			return err
		}
		ave, adev, sdev, skew, curt, err := stats.Moments(values[c.nbPtsDiscard:])
		if PRINT {
			fmt.Printf("Moments : mean=%.3e adev=%.3e sdev=%.3e skew=%.3e curt=%.3e %s\n", ave, adev, sdev, skew, curt, c.title)
		}
		means = append(means, ave)
	}
	base := filepath.Base(c.root)
	return drawLinearFit(c.abscis, means, c.xlabel, "times (ms)", base+c.title, base+"_mean.png")
}

// Draw the data
func drawPointsXY(x, y []float64, xlabel, ylabel, title, outPng string) error {
	// Create the plot
	p, err := plotfunc.NewPlot(title, xlabel, ylabel)
	if err != nil {
		return err
	}
	// Add the data
	if err = plotfunc.AddWithPointsXY(x, y, "", 0, p); err != nil {
		return err
	}
	// Save the plot to a PNG file.
	return p.Save(10*vg.Centimeter, 10*vg.Centimeter, outPng)
}

// Draw the data and
// compute the linear regression that fits the data
func drawLinearFit(x, means []float64, xlabel, ylabel, title, outPng string) error {
	// Create the plot
	p, err := plotfunc.NewPlot(title, xlabel, ylabel)
	if err != nil {
		return err
	}
	// Add the means
	if err = plotfunc.AddWithPointsXY(x, means, "", 0, p); err != nil {
		return err
	}
	// Add a regression line
	a, b, siga, sigb, chi2, sigdat := plotfunc.AddLinearfit(x[1:], means[1:], p)
	if PRINT {
		fmt.Printf("Linear fit : a=%.3e b=%.3e siga=%.3e sigb=%.3e chi2=%.3e sigdat=%.3e %s\n", a, b, siga, sigb, chi2, sigdat, title)
	}
	// Save the plot to a PNG file.
	return p.Save(10*vg.Centimeter, 10*vg.Centimeter, outPng)
}

// Parse a file and draw the data and some stats
func drawFile(filename string, nbPtsDiscard int) error {
	fvalues, err := parseFile(filename)
	if err != nil {
		return err
	}
	base := filepath.Base(filename)
	// Create the plot
	p, err := plotfunc.NewPlot(base, "Msg number", "times (ms)")
	if err != nil {
		return err
	}
	p.Legend.Top = true
	// Add the data
	if err = plotfunc.AddWithPoints(fvalues, base, 0, p); err != nil {
		return err
	}
	// Compute mean regression
	ave, adev, sdev, skew, curt, err := stats.Moments(fvalues[nbPtsDiscard:])
	if PRINT {
		fmt.Printf("Moments : ave=%.3e adev=%.3e sdev=%.3e skew=%.3e curt=%.3e %s\n", ave, adev, sdev, skew, curt, filepath.Base(filename))
	}
	plotfunc.AddHLine(ave, float64(nbPtsDiscard), float64(len(fvalues)), "", color.Black, p)
	// Save the plot to a PNG file.
	return p.Save(10*vg.Centimeter, 10*vg.Centimeter, filepath.Base(filename)+".png")
}

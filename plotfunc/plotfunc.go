package plotfunc

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"plots/stats"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Create a plot with title and axis labels
func NewPlot(title, xlabel, ylabel string) (*plot.Plot, error) {
	p, err := plot.New()
	if err != nil {
		return nil, err
	}
	p.Title.Text = title
	p.X.Label.Text = xlabel
	p.Y.Label.Text = ylabel
	return p, nil
}

// AddLabel add a text at (x, y)
func AddLabel(x, y float64, text string, p *plot.Plot) error {
	labels, err := plotter.NewLabels(plotter.XYLabels{
		XYs: []plotter.XY{
			{X: x, Y: y},
		},
		Labels: []string{text},
	},
	)
	if err != nil {
		return errors.New("could not creates labels plotter: %+v")
	}
	p.Add(labels)
	return nil
}

// AddHLine Add a horizontal line
func AddHLine(y, xmin, xmax float64, legend string, c color.Color, p *plot.Plot) error {
	return AddStraightLine(xmin, y, xmax, y, legend, c, p)
}

// AddVLine Add a vertical line
func AddVLine(x, ymin, ymax float64, legend string, c color.Color, p *plot.Plot) error {
	return AddStraightLine(x, ymin, x, ymax, legend, c, p)
}

// AddStraightLine Add a straight line
func AddStraightLine(xmin, ymin, xmax, ymax float64, legend string, c color.Color, p *plot.Plot) error {
	pts := make(plotter.XYs, 2)
	pts[0].X = xmin
	pts[0].Y = ymin
	pts[1].X = xmax
	pts[1].Y = ymax
	line, err := plotter.NewLine(pts)
	if err != nil {
		return err
	}
	line.Color = c
	p.Add(line)
	if legend != "" {
		p.Legend.Add(legend, line)
	}
	return nil
}

// SimpleAdd Use the simplest way to create the plot (using plotUtil)
func SimpleAdd(data []float64, legend string, p *plot.Plot) error {
	return plotutil.AddLinePoints(p, legend, CreatePoints(data))
}

// AddWithLine Draw the data with a line
func AddWithLine(data []float64, legend string, p *plot.Plot) error {
	lpLine, err := plotter.NewLine(CreatePoints(data))
	if err != nil {
		return err
	}
	c1, c2, c3 := uint8(rand.Int()%255), uint8(rand.Int()%255), uint8(rand.Int()%255)
	lpLine.Color = color.RGBA{R: c1, G: c2, B: c3, A: 255}
	p.Add(lpLine)
	if legend != "" {
		p.Legend.Add(legend, lpLine)
	}
	return nil
}

// AddWithPoints Draw the data with points
func AddWithPoints(data []float64, legend string, p *plot.Plot) error {
	points, err := plotter.NewScatter(CreatePoints(data))
	if err != nil {
		return err
	}
	points.Radius = 1
	c1, c2, c3 := uint8(rand.Int()%255), uint8(rand.Int()%255), uint8(rand.Int()%255)
	points.Shape = draw.CircleGlyph{}
	points.Color = color.RGBA{R: c1, G: c2, B: c3, A: 255}
	p.Add(points)
	if legend != "" {
		p.Legend.Add(legend, points)
	}
	return nil
}

type commaTicks struct{}

// Ticks computes the default tick marks, and define the label for the major tick marks.
func (commaTicks) Ticks(min, max float64) []plot.Tick {
	tks := plot.DefaultTicks{}.Ticks(min, max)
	for i, t := range tks {
		if t.Label == "" { // Skip minor ticks, they are fine.
			continue
		}
		tks[i].Label = fmt.Sprintf("%.2f", tks[i].Value)
	}
	return tks
}

type errPoints struct {
	plotter.XYs
	plotter.YErrors
}

// AddWithErrXY Draw the data (x, y) with their error bars (devs)
func AddWithErrXY(x, y, devs []float64, legend string, p *plot.Plot) error {
	rand.Seed(time.Now().UnixNano())
	xys := make(plotter.XYs, len(y))
	yer := make(plotter.YErrors, len(y))
	for j := range xys {
		xys[j].X = x[j]
		xys[j].Y = y[j]
		yer[j].High = devs[j]
		yer[j].Low = devs[j]
	}
	data := errPoints{XYs: xys, YErrors: yer}
	yerrs, err := plotter.NewYErrorBars(data)
	if err != nil {
		return err
	}
	scatter, err := plotter.NewScatter(data)
	if err != nil {
		return err
	}
	scatter.Radius = 1
	scatter.Shape = draw.CircleGlyph{}
	c1, c2, c3 := uint8(rand.Int()%255), uint8(rand.Int()%255), uint8(rand.Int()%255)
	scatter.Color = color.RGBA{R: c1, G: c2, B: c3, A: 255}
	yerrs.Color = color.RGBA{R: c1, G: c2, B: c3, A: 255}
	p.Add(scatter, yerrs)
	if legend != "" {
		p.Legend.Add(legend, scatter)
	}
	p.Legend.Top = true
	return nil
}

// AddWithPoints Draw the data with points
func AddWithPointsXY(x, y []float64, legend string, p *plot.Plot) error {
	points, err := plotter.NewScatter(CreatePointsXY(x, y))
	if err != nil {
		return err
	}
	points.Radius = 2
	c1, c2, c3 := uint8(rand.Int()%255), uint8(rand.Int()%255), uint8(rand.Int()%255)
	points.Shape = draw.CircleGlyph{}
	points.Color = color.RGBA{R: c1, G: c2, B: c3, A: 255}
	p.Add(points)
	if legend != "" {
		p.Legend.Add(legend, points)
	}
	p.Legend.Top = true
	p.Y.Tick.Marker = commaTicks{}
	return nil
}

// CreatePoints Transform a []float64 into a plotter
func CreatePoints(data []float64) plotter.XYs {
	pts := make(plotter.XYs, len(data))
	for i := range data {
		pts[i].X = float64(i)
		pts[i].Y = data[i]
	}
	return pts
}

// CreatePointsXY Transform the x, y slices into plotter
func CreatePointsXY(x, y []float64) plotter.XYs {
	pts := make(plotter.XYs, len(x))
	for i := range x {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}
	return pts
}

// Interpolation of x, y (with deviations devs) with a polynom of degree n
// Returns the fitting coefficients
func AddPolyfit(x, y, devs []float64, n int, p *plot.Plot) ([]float64, error) {
	// Add regression trinome
	ia := make([]bool, n)
	cofs := make([]float64, n)
	for i := 0; i < n; i++ {
		ia[i] = true
		cofs[i] = 1.
	}
	if _, _, err := stats.Lfit(x, y, devs, cofs, ia, n, stats.Fcoefs); err != nil {
		return nil, err
	}
	poly := func(x float64) float64 {
		return stats.Fpoly(x, cofs)
	}
	fp := plotter.NewFunction(poly)
	fp.Color = color.RGBA{B: 255, A: 255}
	fp.Dashes = []vg.Length{vg.Points(10), vg.Points(10)}
	p.Add(fp)
	return cofs, nil
}

// Interpolation of x, y with a straight line
// Return a, b, siga, sigb, chi2 and sigdat
func AddLinearfit(x, y []float64, p *plot.Plot) (float64, float64, float64, float64, float64, float64) {
	a, b, siga, sigb, chi2, sigdat := stats.LSFitLinear(x, y)
	// Add regression line
	ymin1 := b*x[0] + a
	ymax1 := b*x[len(x)-1] + a
	legend1 := fmt.Sprintf("%0.2e x + %0.2f", b, a)
	AddStraightLine(x[0], ymin1, x[len(x)-1], ymax1, legend1, color.RGBA{R: 255, A: 255}, p)
	return a, b, siga, sigb, chi2, sigdat
}

// Add the normal distribution function with "mean" average and "sdev" standard deviation
func AddGaussian(mean, sdev float64, p *plot.Plot) {
	sdev2 := 2. * sdev * sdev
	pi2 := sdev * math.Sqrt(2.*math.Pi)
	gauss := func(x float64) float64 {
		return math.Exp(-(x-mean)*(x-mean)/sdev2) / pi2
	}
	gaus := plotter.NewFunction(gauss)
	gaus.Color = color.RGBA{B: 255, A: 255}
	p.Add(gaus)
}

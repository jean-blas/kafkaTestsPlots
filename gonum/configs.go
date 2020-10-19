package main

import (
	"path/filepath"
	"plots/plotfunc"
	"strconv"
	"sync"
)

// Definition of a Config fields
type Config struct {
	nbPtsDiscard int      // [optional] number of points to discard from the beginning when fitting (default 0)
	root         string   // root folder of the data files
	prefix       string   // constant prefix in the name of all data files
	postfix      string   // [optional] constant postfix in the name of all data files (mostly empty)
	sufix        []string // variable part in the name of the data files
	xlabel       string   // xlabel of the graphics
	abscisIsNb   bool     // [optional] true if the nb of messages is represented by the absissa, needed to compute the throughput (default false)
	abscisIsSz   bool     // [optional] true if the size of the messages is represented by the absissa, needed to compute the throughput (default false)
	title        string   // [optional] Add a title line (default is empty)
	mb           float64  // [optional] default size of the messages in Mb (default = 0.1)
	ndata        int      // [optional] number of data values in each file (default = 2000)

	files  []string  // real file names (root + prefix + sufix + postfix), computed automatically
	abscis []float64 // corresponding abscissa of the data files, in the correct unit, computed automatically
}

// Create a string legend from the config fields
func (c Config) legend() string {
	// return fmt.Sprintf("size = %0.2f Mb", c.mb) // return the size
	return filepath.Base(filepath.Dir(c.root)) // return the last folder
}

// Prepare the config object before using it in the draw functions
func (c *Config) prepare() {
	// Replace sufix with the real path (root + prefix + sufix + postfix) for each sufix
	sfx := make([]string, len(c.sufix))
	for i := range sfx {
		sfx[i] = filepath.Join(c.root, c.prefix+c.sufix[i]+c.postfix)
	}
	c.files = sfx
	// Compute the abscissa
	abs := make([]float64, len(c.sufix))
	for i := range c.sufix {
		vi, err := strconv.Atoi(c.sufix[i])
		if err != nil {
			panic(err)
		}
		abs[i] = float64(vi)
	}
	c.abscis = abs
	// default size of messages if not set
	if c.mb == 0 {
		c.mb = 0.1
	}
	if c.ndata == 0 {
		c.ndata = 2000
	}
}

// Add the diagram types here
type Draws int

const (
	Dall           Draws = iota // Draw all diagram types (except Dcompare)
	Dfile                       // Draw the raw points
	DhistoFile                  // Draw the equivalent histogram
	DmeansFile                  // Draw the computed mean
	DmeansErrFiles              // Draw the computed mean with error deviations
	DslideFile                  // Draw a sliding window accross the points
	Dthroughput                 // Draw the throughput
	DnbMsgPerSec                // Draw the number of messages per second
)

func (d Draws) String() string {
	return [...]string{"Draw all", "Comparison", "Draw file raw data", "Draw histograms", "Draw means",
		"Draw means with errors", "Draw a sliding window", "Draw throughput", "Draw the number of messages per seconds"}[d]
}

// If true, print the moments ofCqueueBufMaxKbytes_30k, CqueueBufMaxKbytes_100k, CqueueBufMaxKbytes_300k, CqueueBufMaxKbytes_3000k the distribution for each diagram
const PRINT = false

// Window interval when using drawSlide
const NVAL = 5

// Number of columns of the histograms
const NCOL = 30

// Main entry point.
func main() {
	// draw := Dall                  // Define the diagrams to create (Dall = all diagrams)
	// fileNb := -1                  // file to process as example (-1 = all files)
	// conf := []Confs{CmsgSizeAck1} //Define the configs to process (Call = all configs)
	// process(draw, conf, fileNb, "")

	compareAll()
}

// Run all comparisons in parallel
func compareAll() {
	ComparePNGsuffix = "per_partition" // sufix added to PNG names
	plotfunc.N = 10
	var wg sync.WaitGroup
	wg.Add(1)
	// go func() {
	// 	compareConfigs([]Confs{Cp6_queueBufMaxMs_100k, Cp36_queueBufMaxMs_100k, Cp72_queueBufMaxMs_100k, Cp108_queueBufMaxMs_100k, Cp180_queueBufMaxMs_100k, Cp360_queueBufMaxMs_100k})
	// 	wg.Done()
	// }()
	// go func() {
	// 	compareConfigs([]Confs{Cp6_queuedMinMessages_100k, Cp36_queuedMinMessages_100k, Cp72_queuedMinMessages_100k, Cp108_queuedMinMessages_100k, Cp180_queuedMinMessages_100k, Cp360_queuedMinMessages_100k})
	// 	wg.Done()
	// }()
	// go func() {
	// 	compareConfigs([]Confs{Cp6_queueBufMaxMsg_100k, Cp36_queueBufMaxMsg_100k, Cp72_queueBufMaxMsg_100k, Cp108_queueBufMaxMsg_100k, Cp180_queueBufMaxMsg_100k, Cp360_queueBufMaxMsg_100k})
	// 	wg.Done()
	// }()
	// go func() {
	// 	compareConfigs([]Confs{Cp6_batchNumMsg_100k, Cp36_batchNumMsg_100k, Cp72_batchNumMsg_100k, Cp108_batchNumMsg_100k, Cp180_batchNumMsg_100k, Cp360_batchNumMsg_100k})
	// 	wg.Done()
	// }()
	// go func() {
	// 	compareConfigs([]Confs{Cp6_fetchMinBytes_100k, Cp36_fetchMinBytes_100k, Cp72_fetchMinBytes_100k, Cp108_fetchMinBytes_100k, Cp180_fetchMinBytes_100k, Cp360_fetchMinBytes_100k})
	// 	wg.Done()
	// }()
	// go func() {
	// 	compareConfigs([]Confs{Cp6_fetchWaitMaxMs_100k, Cp36_fetchWaitMaxMs_100k, Cp72_fetchWaitMaxMs_100k, Cp108_fetchWaitMaxMs_100k, Cp180_fetchWaitMaxMs_100k, Cp360_fetchWaitMaxMs_100k})
	// 	wg.Done()
	// }()
	// go func() {
	// 	compareConfigs([]Confs{Cp6_queueBufMaxKbytes_100k, Cp36_queueBufMaxKbytes_100k, Cp72_queueBufMaxKbytes_100k, Cp108_queueBufMaxKbytes_100k, Cp180_queueBufMaxKbytes_100k, Cp360_queueBufMaxKbytes_100k})
	// 	wg.Done()
	// }()
	go func() {
		compareConfigs([]Confs{Cp6_msgSize, Cp36_msgSize, Cp72_msgSize, Cp108_msgSize})
		wg.Done()
	}()

	wg.Wait()
}

// Process the config according to the parameters
func process(draw Draws, conf []Confs, fileNb int, comparePNGsuffix string) {
	ComparePNGsuffix = comparePNGsuffix
	plotfunc.N = len(conf)
	if plotfunc.N > 1 {
		compareConfigs(conf)
	} else {
		switch conf[0] {
		case Call:
			// Process all configs in parallel
			var wg sync.WaitGroup
			for _, cfg := range Configs {
				wg.Add(1)
				go func(c Config) {
					drawConfig(c, draw, fileNb)
					wg.Done()
				}(cfg)
			}
			wg.Wait()
		default:
			// Process the only config defined (only one here)
			for _, c := range conf {
				drawConfig(Configs[c-1], draw, fileNb)
			}
		}
	}
}

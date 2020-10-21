package main

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"plots/plotfunc"
	"strconv"
	"sync"
)

// Definition of a Config fields
type Config struct {
	name         string   // unique name of the config
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

// Return the index of the config that has the same name, or -1 if not found
func findConfigIdx(name string) int {
	for i, conf := range Configs {
		if name == conf.name {
			return i
		}
	}
	return -1
}

// Transform the slice of strings into the slice of corresponding configs
func toConfigs(cs []string) ([]Config, error) {
	cfg := make([]Config, len(cs))
	for i, name := range cs {
		idx := findConfigIdx(name)
		if idx == -1 {
			return nil, errors.New("Config not found with name : " + name)
		}
		cfg[i] = Configs[idx]
	}
	return cfg, nil
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

var draws = []Draws{
	Dall, Dfile, DhistoFile, DmeansFile, DmeansErrFiles, DslideFile, Dthroughput, DnbMsgPerSec,
}

func (d Draws) String() string {
	return [...]string{"Draw all", "Draw file raw data", "Draw histograms", "Draw means",
		"Draw means with errors", "Draw a sliding window", "Draw throughput", "Draw the number of messages per seconds"}[d]
}

// Describe the different draws in the help (-h)
func helpDraw() string {
	var s string
	for i, d := range draws {
		s = fmt.Sprintf("%s\n%d = %s", s, i, d)
	}
	return s
}

// If true, print the moments ofCqueueBufMaxKbytes_30k, CqueueBufMaxKbytes_100k, CqueueBufMaxKbytes_300k, CqueueBufMaxKbytes_3000k the distribution for each diagram
const PRINT = false

// Window interval when using drawSlide
const NVAL = 5

// Number of columns of the histograms
const NCOL = 30

// Main entry point.
func main() {
	d := flag.Int("d", 0, "Drawing type (default 0)"+helpDraw())
	n := flag.Int("n", -1, "File number to process as example or -1 for all")
	c := flag.String("c", "msgSizeAck1", "Name of the config to process")
	compar := flag.Bool("C", false, "Run comparison mode")
	flag.Parse()

	if *compar {
		compareAll()
		return
	}
	if *d < 0 || *d >= len(draws) {
		fmt.Println("Error : bad drawing type. Should be in [ 0, ", len(draws)-1, "]")
		return
	}
	if *n < -1 {
		fmt.Println("Error : bad file number.")
		return
	}
	if *c == "all" {
		processAllConfigs(draws[*d], *n)
		return
	}
	idx := findConfigIdx(*c)
	if idx == -1 {
		fmt.Println("No config found with name : ", *c)
		return
	}
	drawConfig(Configs[idx], draws[*d], *n)
}

// Compare the configs defined by their name in the given slice one each other
func compareConfig(names []string) error {
	confs, err := toConfigs(names)
	if err != nil {
		return err
	}
	return compareConfigs(confs)
}

// Run all comparisons in parallel
func compareAll() {
	ComparePNGsuffix = "per_partition" // sufix added to PNG names
	plotfunc.N = 10
	var wg sync.WaitGroup
	wg.Add(8)
	go func() {
		if err := compareConfig([]string{"p6_queueBufMaxMs_100k", "p36_queueBufMaxMs_100k", "p72_queueBufMaxMs_100k", "p108_queueBufMaxMs_100k", "p180_queueBufMaxMs_100k", "p360_queueBufMaxMs_100k"}); err != nil {
			fmt.Println(err)
		}
		wg.Done()
	}()
	go func() {
		if err := compareConfig([]string{"p6_queuedMinMessages_100k", "p36_queuedMinMessages_100k", "p72_queuedMinMessages_100k", "p108_queuedMinMessages_100k", "p180_queuedMinMessages_100k", "p360_queuedMinMessages_100k"}); err != nil {
			fmt.Println(err)
		}
		wg.Done()
	}()
	go func() {
		if err := compareConfig([]string{"p6_queueBufMaxMsg_100k", "p36_queueBufMaxMsg_100k", "p72_queueBufMaxMsg_100k", "p108_queueBufMaxMsg_100k", "p180_queueBufMaxMsg_100k", "p360_queueBufMaxMsg_100k"}); err != nil {
			fmt.Println(err)
		}
		wg.Done()
	}()
	go func() {
		if err := compareConfig([]string{"p6_batchNumMsg_100k", "p36_batchNumMsg_100k", "p72_batchNumMsg_100k", "p108_batchNumMsg_100k", "p180_batchNumMsg_100k", "p360_batchNumMsg_100k"}); err != nil {
			fmt.Println(err)
		}
		wg.Done()
	}()
	go func() {
		if err := compareConfig([]string{"p6_fetchMinBytes_100k", "p36_fetchMinBytes_100k", "p72_fetchMinBytes_100k", "p108_fetchMinBytes_100k", "p180_fetchMinBytes_100k", "p360_fetchMinBytes_100k"}); err != nil {
			fmt.Println(err)
		}
		wg.Done()
	}()
	go func() {
		if err := compareConfig([]string{"p6_fetchWaitMaxMs_100k", "p36_fetchWaitMaxMs_100k", "p72_fetchWaitMaxMs_100k", "p108_fetchWaitMaxMs_100k", "p180_fetchWaitMaxMs_100k", "p360_fetchWaitMaxMs_100k"}); err != nil {
			fmt.Println(err)
		}
		wg.Done()
	}()
	go func() {
		if err := compareConfig([]string{"p6_queueBufMaxKbytes_100k", "p36_queueBufMaxKbytes_100k", "p72_queueBufMaxKbytes_100k", "p108_queueBufMaxKbytes_100k", "p180_queueBufMaxKbytes_100k", "p360_queueBufMaxKbytes_100k"}); err != nil {
			fmt.Println(err)
		}
		wg.Done()
	}()
	go func() {
		if err := compareConfig([]string{"p6_msgSize", "p36_msgSize", "p72_msgSize", "p108_msgSize", "p180_msgSize", "p360_msgSize"}); err != nil {
			fmt.Println(err)
		}
		wg.Done()
	}()
	wg.Wait()
}

// Process all configs according to the parameters in parallel
func processAllConfigs(draw Draws, fileNb int) {
	plotfunc.N = 1
	var wg sync.WaitGroup
	for _, cfg := range Configs {
		wg.Add(1)
		go func(c Config) {
			drawConfig(c, draw, fileNb)
			wg.Done()
		}(cfg)
	}
	wg.Wait()
}

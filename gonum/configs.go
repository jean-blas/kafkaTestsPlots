package main

// Definition of a Config fields
type Config struct {
	nbPtsDiscard int      // [optional] number of points to discard when fitting (default 0)
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

var Configs = []Config{
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/msgSizeAck1",
		prefix:       "ms_size_ack1_",
		sufix: []string{"100", "200", "300", "400", "500", "600", "700", "800", "900", "1000",
			"1500", "2000", "2500", "3000", "3500", "4000", "4500", "5000"},
		xlabel:     "size (kb)",
		postfix:    "k_n2000",
		abscisIsSz: true,
		ndata:      2000,
		mb:         0.1,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/fetchWaitMaxMs",
		prefix:       "fetch_wait_max_ms_100k_n2000_",
		sufix:        []string{"2", "4", "6", "8", "10", "20", "40", "60", "80", "100", "200", "400", "600", "800", "1000"},
		xlabel:       "fetch.wait.max.ms",
		ndata:        2000,
		mb:           0.1,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queuedMinMessages",
		prefix:       "queued.min.messages_100k_n2000_",
		sufix: []string{"2000", "4000", "6000", "8000", "10000", "20000",
			"40000", "60000", "80000", "100000", "200000", "400000", "600000", "800000", "1000000"},
		xlabel: "queued.min.messages",
		ndata:  2000,
		mb:     0.1,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/fetchMinBytes_100k",
		prefix:       "fetch_min_bytes_100k_n10000_",
		sufix: []string{"200", "400", "600", "800", "1000", "2000", "4000", "6000", "8000", "10000", "20000",
			"40000", "60000", "80000", "100000", "200000", "400000", "600000", "800000", "1000000"},
		xlabel: "fetch.min.bytes",
		ndata:  10000,
		mb:     0.1,
		title:  "\nsize=100k n=10000 linger.ms=100",
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/fetchMinBytes_300k",
		prefix:       "fetch_min_bytes_300k_n10000_",
		sufix: []string{"200", "400", "600", "800", "1000", "2000", "4000", "6000", "8000", "10000", "20000",
			"40000", "60000", "80000", "100000", "200000", "400000", "600000", "800000", "1000000"},
		xlabel: "fetch.min.bytes",
		mb:     0.3,
		ndata:  10000,
		title:  "\nsize=300k n=10000 linger.ms=100",
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/fetchMinBytes_3000k",
		prefix:       "fetch_min_bytes_3000k_n10000_",
		sufix: []string{"400", "600", "800", "1000", "2000", "4000", "6000", "8000", "10000", "20000",
			"40000", "60000", "80000", "100000", "200000", "400000", "600000", "800000", "1000000"},
		xlabel: "fetch.min.bytes",
		mb:     3,
		ndata:  10000,
		title:  "\nsize=3m n=10000 linger.ms=100",
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxKbytes_big",
		prefix:       "queue.buffering.max.msg_100k_n2000_",
		sufix:        []string{"200", "2000", "20000", "200000", "1000000"},
		xlabel:       "queue.buffering.max.kbytes",
		ndata:        2000,
		mb:           0.1,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxMsg_3000k_n10000",
		prefix:       "queue.buffering.max.messages_3000k_n10000_",
		sufix:        []string{"4", "6", "8", "10", "20", "30", "40", "50", "60", "80", "100", "200", "400"},
		xlabel:       "queue.buffering.max.msg",
		mb:           3.0,
		ndata:        10000,
		title:        "\nsize=3m n=10000 linger.ms=100",
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxMsg_ms200_900k",
		prefix:       "queue.buffering.max.msg_900k_n2000_",
		sufix: []string{"200", "400", "600", "800", "1000", "2000", "4000", "6000", "8000", "10000",
			"20000", "40000", "60000", "80000", "100000", "200000", "400000", "600000", "800000", "1000000"},
		xlabel: "queue.buffering.max.msg",
		mb:     0.9,
		ndata:  2000,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxMsg_ms200_700k",
		prefix:       "queue.buffering.max.msg_700k_n2000_",
		sufix: []string{"200", "400", "600", "800", "1000", "2000", "4000",
			"6000", "8000", "10000", "20000", "40000", "60000", "80000", "100000", "200000", "400000", "600000",
			"800000", "1000000"},
		xlabel: "queue.buffering.max.msg",
		mb:     0.7,
		ndata:  2000,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxMsg_ms200_500k",
		prefix:       "queue.buffering.max.msg_500k_n2000_",
		sufix: []string{"200", "400", "600", "800", "1000", "2000", "4000",
			"6000", "8000", "10000", "20000", "40000", "60000", "80000", "100000", "200000", "400000", "600000",
			"800000", "1000000"},
		xlabel: "queue.buffering.max.msg",
		mb:     0.5,
		ndata:  2000,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxMsg_ms200_300k",
		prefix:       "queue.buffering.max.messages_300k_n10000_",
		sufix:        []string{"2", "4", "6", "8", "10", "20", "40", "60", "80", "100", "200", "400"},
		xlabel:       "queue.buffering.max.msg",
		mb:           0.3,
		ndata:        10000,
		title:        "\nsize=300kb n=10000 linger.ms=200",
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxMsg_ms200_100k",
		prefix:       "queue.buffering.max.msg_100k_n2000_",
		sufix: []string{"200", "400", "600", "800",
			"1000", "2000", "4000", "6000", "8000", "10000", "20000", "40000", "60000", "80000", "100000"},
		xlabel: "queue.buffering.max.msg",
		ndata:  2000,
		mb:     0.1,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxMsg_ms100_30k",
		prefix:       "queue.buffering.max.messages_30k_n10000_",
		sufix:        []string{"2", "4", "6", "8", "10", "20", "40", "60", "80", "100", "200", "400", "600", "800", "1000"},
		xlabel:       "queue.buffering.max.msg",
		mb:           0.03,
		ndata:        10000,
		title:        "\nsize=30kb n=10000 linger.ms=100",
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxMsg_ms100_100k",
		prefix:       "queue.buffering.max.msg_100k_n2000_",
		sufix: []string{"4", "6", "8", "10", "20", "40", "60", "80", "100", "200", "400", "600", "800",
			"1000", "2000", "4000", "6000", "8000", "10000"},
		xlabel: "queue.buffering.max.msg",
		ndata:  2000,
		mb:     0.1,
		title:  "\nsize=100kb n=10000 linger.ms=100",
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxMsg_ms100_300k",
		prefix:       "queue.buffering.max.messages_300k_n10000_",
		sufix:        []string{"2", "4", "6", "8", "10", "20", "40", "60", "80", "100", "200", "400"},
		xlabel:       "queue.buffering.max.msg",
		mb:           0.3,
		ndata:        10000,
		title:        "\nsize=300kb n=10000 linger.ms=100",
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxMsg_ms0_300k",
		prefix:       "queue.buffering.max.messages_300k_n10000_",
		sufix:        []string{"2", "4", "6", "8", "10", "20", "40", "60", "80", "100", "200", "400"},
		xlabel:       "queue.buffering.max.msg",
		mb:           0.3,
		ndata:        10000,
		title:        "\nsize=300kb n=10000 linger.ms=0",
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxMsg_ms0_100k",
		prefix:       "queue.buffering.max.msg_100k_n2000_",
		sufix:        []string{"1", "3", "7", "11", "50", "500", "1000", "1500", "2000", "2500", "4000", "6000", "8000", "10000"},
		xlabel:       "queue.buffering.max.msg",
		ndata:        2000,
		mb:           0.1,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/batchNumMsg_100k",
		prefix:       "batch.num.msg_100k_n2000_",
		sufix: []string{"2", "4", "6", "8", "10", "20", "40", "60", "80", "100", "200", "400", "600", "800", "1000",
			"2000", "4000", "6000", "8000", "10000"},
		xlabel: "batch.num.msg",
		ndata:  2000,
		mb:     0.1,
		title:  "\nsize=100kb n=2000 linger.ms=100",
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/batchNumMsg_300k",
		prefix:       "batch.num.msg_300k_n10000_",
		sufix:        []string{"200", "400", "600", "800", "1000", "2000", "4000", "6000", "8000", "10000"},
		xlabel:       "batch.num.msg",
		mb:           0.3,
		ndata:        10000,
		title:        "\nsize=300kb n=10000 linger.ms=100",
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/batchNumMsg_5Mb",
		prefix:       "batch.num.messages_5000k_n2000_",
		sufix: []string{"20", "40", "60", "80", "100", "200", "400", "600", "800", "1000",
			"2000", "4000", "6000", "8000", "10000"},
		xlabel: "batch.num.msg",
		mb:     5.0,
		ndata:  2000,
		title:  "\nsize=5m n=2000 linger.ms=100",
	},
	{
		nbPtsDiscard: 0,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/nbMsg",
		prefix:       "nbMsg_100k_n",
		sufix: []string{"110", "210", "310", "410", "510", "610", "710", "810", "910",
			"1100", "210", "3100", "4100", "5100", "6100", "7100", "8100", "9100", "11000"},
		xlabel:     "nb of messages",
		abscisIsNb: true,
		ndata:      2000,
		mb:         0.1,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxKbytes_30k",
		prefix:       "queue.buffering.max.kbytes_30k_n10000_",
		sufix: []string{"200", "400", "600", "800", "1000", "2000", "4000", "6000", "8000", "10000", "20000", "40000", "60000", "80000", "100000",
			"200000", "400000", "600000", "800000", "1000000"},
		xlabel: "queue.buffering.max.kbytes",
		mb:     0.03,
		ndata:  10000,
		title:  "\nsize=30kb n=10000 linger.ms=100",
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxKbytes_100k",
		prefix:       "queue.buffering.max.kb_100k_n2000_",
		sufix: []string{"200", "400", "600", "800", "1000", "2000", "4000", "6000", "8000", "10000", "20000", "40000", "60000", "80000", "100000",
			"200000", "400000", "600000", "800000", "1000000"},
		xlabel: "queue.buffering.max.kbytes",
		ndata:  2000,
		mb:     0.1,
		title:  "\nsize=100kb n=10000 linger.ms=100",
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxKbytes_300k",
		prefix:       "queue.buffering.max.kbytes_300k_n10000_",
		sufix: []string{"400", "600", "800", "1000", "2000", "4000", "6000", "8000", "10000", "20000", "40000", "60000", "80000", "100000",
			"200000", "400000", "600000", "800000", "1000000"},
		xlabel: "queue.buffering.max.kbytes",
		mb:     0.3,
		ndata:  10000,
		title:  "\nsize=300kb n=10000 linger.ms=100",
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxMs_30k",
		prefix:       "queue.buffering.max.ms_30k_n10000_",
		sufix:        []string{"2", "4", "6", "8", "10", "12", "14", "16", "18", "20"},
		xlabel:       "queue.buffering.max.ms",
		mb:           0.03,
		ndata:        10000,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxMs_100k",
		prefix:       "queue.buffering.max.ms_100k_n2000_",
		sufix:        []string{"0", "3", "6", "9", "12", "15"},
		xlabel:       "queue.buffering.max.ms",
		ndata:        2000,
		mb:           0.1,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxMs_300k",
		prefix:       "queue.buffering.max.ms_300k_n10000_",
		sufix:        []string{"2", "4", "6", "8", "10", "12", "14", "16", "18", "20"},
		xlabel:       "queue.buffering.max.ms",
		mb:           0.3,
		ndata:        10000,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxMs_1Mb",
		prefix:       "queue.buffering.max.ms_1000k_n2000_",
		sufix:        []string{"1", "3", "5", "7", "9", "11", "13", "15", "17", "19"},
		xlabel:       "queue.buffering.max.ms",
		mb:           1.0,
		ndata:        2000,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/queueBufMaxMs_5Mb",
		prefix:       "queue.buffering.max.ms_5000k_n2000_",
		sufix:        []string{"1", "3", "5", "7", "9", "11", "13", "15", "17", "19"},
		xlabel:       "queue.buffering.max.ms",
		mb:           5.0,
		ndata:        2000,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/sizeMsgGzip",
		prefix:       "latency_gzip_",
		postfix:      "k_2000",
		sufix: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "20", "30", "40", "50", "60", "70", "80", "90",
			"100", "200", "300", "400", "500", "600", "700", "800", "900"},
		xlabel:     "size (kb)",
		abscisIsSz: true,
		ndata:      2000,
		mb:         0.1,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/sizeMsgLz4",
		prefix:       "latency_lz4_",
		postfix:      "k_2000",
		sufix: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "20", "30", "40", "50", "60", "70", "80", "90",
			"100", "200", "300", "400", "500", "600", "700", "800", "900"},
		xlabel:     "size (kb)",
		abscisIsSz: true,
		ndata:      2000,
		mb:         0.1,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/sizeMsgSnappy",
		prefix:       "latency_snappy_",
		postfix:      "k_2000",
		sufix: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "20", "30", "40", "50", "60", "70", "80", "90",
			"100", "200", "300", "400", "500", "600", "700", "800", "900"},
		xlabel:     "size (kb)",
		abscisIsSz: true,
		ndata:      2000,
		mb:         0.1,
	},
	{
		nbPtsDiscard: 500,
		root:         "/home/jimbert/Projects/LibRdKafka/messagebrokerclient/notes/benchmarks/sizeMsg",
		prefix:       "latency_none_",
		postfix:      "k_2000",
		sufix: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "20", "30", "40", "50", "60", "70", "80", "90",
			"100", "200", "300", "400", "500", "600", "700", "800", "900", "1000", "1500", "2000", "2500", "3000", "3500",
			"4000", "4500", "5000"},
		xlabel:     "size (kb)",
		abscisIsSz: true,
		ndata:      2000,
		mb:         0.1,
	},
}

// Add the configs below, in the same order they are defined above (order matters)
type Confs int

const (
	Call Confs = iota
	CmsgSizeAck1
	CfetchWaitMaxMs
	CqueuedMinMessages
	CfetchMinBytes_100k
	CfetchMinBytes_300k
	CfetchMinBytes_3000k
	CqueueBufMaxKbytes_big
	CqueueBufMaxMsg_3000k_n10000
	CqueueBufMaxMsg_ms200_900k
	CqueueBufMaxMsg_ms200_700k
	CqueueBufMaxMsg_ms200_500k
	CqueueBufMaxMsg_ms200_300k
	CqueueBufMaxMsg_ms200_100k
	CqueueBufMaxMsg_ms100_30k
	CqueueBufMaxMsg_ms100_100k
	CqueueBufMaxMsg_ms100_300k
	CqueueBufMaxMsg_ms0_300k
	CqueueBufMaxMsg_ms0_100k
	CbatchNumMsg_100k
	CbatchNumMsg_300k
	CbatchNumMsg_5Mb
	CnbMsg
	CqueueBufMaxKbytes_30k
	CqueueBufMaxKbytes_100k
	CqueueBufMaxKbytes_300k
	CqueueBufMaxMs_30k
	CqueueBufMaxMs_100k
	CqueueBufMaxMs_300k
	CqueueBufMaxMs_1Mb
	CqueueBufMaxMs_5Mb
	CsizeMsgGzip
	CsizeMsgLz4
	CsizeMsgSnappy
	CsizeMsg
)

func (c Confs) String() string {
	return [...]string{"Call", "msg size with Ack=1", "fetch.wait.max.ms", "queued.min.messages", "fetch.min.bytes (100k)",
		"fetch.min.bytes (300k)", "fetch.min.bytes (3m)", "queue.buffering.max.kbytes big", "queue.buffering.max.messages with linger.ms=100 (3m)",
		"queue.buffering.max.messages with linger.ms=200 (900k)", "queue.buffering.max.messages with linger.ms=200 (700k)",
		"queue.buffering.max.messages with linger.ms=200 (500k)", "queue.buffering.max.messages with linger.ms=200 (300k)",
		"queue.buffering.max.messages with linger.ms=200 (100k)", "queue.buffering.max.messages with linger.ms=100 (30k)",
		"queue.buffering.max.messages with linger.ms=100 (100k)", "queue.buffering.max.messages with linger.ms=100 (300k)",
		"queue.buffering.max.messages with linger.ms=0 (300k)", "queue.buffering.max.messages with linger.ms=0 (100k)",
		"batch.num.message (100k)", "batch.num.message (300k)", "batch.num.message (5M)", "nb of messages",
		"queue.buffering.max.kbytes (30k)", "queue.buffering.max.kbytes (100k)", "queue.buffering.max.kbytes (300k)",
		"queue.buffering.max.ms (30k)", "queue.buffering.max.ms (100k)", "queue.buffering.max.ms (300k)",
		"queue.buffering.max.ms (1M)", "queue.buffering.max.ms (5M)",
		"size of messages with Gzip", "size of messages with Lz4", "size of messages with GzipSnappy",
		"size of messages with no compression"}[c]
}

// Add the diagram types here
type Draws int

const (
	Dall           Draws = iota // Draw all diagram types (except Dcompare)
	Dcompare                    // Draw the defined configs (C...) in the same diagram for comparison
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

// If true, print the moments of the distribution for each diagram
const PRINT = true

// window interval when using drawSlide
const NVAL = 5

// Main entry point.
// Define the actions to do hereafter.
// Ex1 Dconpare => comparison beween 3 configs (CfetchMinBytes_100k, CfetchMinBytes_300k, CfetchMinBytes_3000k)
// 		var draw Draws = Dcompare
//		var conf []Confs = []Confs{CfetchMinBytes_100k, CfetchMinBytes_300k, CfetchMinBytes_3000k}
// Ex2 fileNb = -1 => create all diagrams for CqueueBufMaxMsg_ms100_30k (with all histo)
//		var draw Draws = Dall
//		var conf []Confs = []Confs{CqueueBufMaxMsg_ms100_30k}
//		fileNb := -1
// Ex3 Dall + fileNb >= 0 => create only interesting diagrams for CqueueBufMaxMsg_ms100_30k (only 1 histo)
//		var draw Draws = Dall
//		var conf []Confs = []Confs{CqueueBufMaxMsg_ms100_30k}
//		fileNb := 4
// Ex4 Call => create the same diagram (here throughputs) for all configs
//		var draw Draws = Dthroughput
//		var conf []Confs = []Confs{Call}
//		fileNb := 4
// Ex5 Call + Dall => create all diagrams for all configs (big number of generated files)
// Ex5 Call + Dall + fileNb=-1 => create all diagrams for all configs with all histo (huge number of generated files)
func main() {
	// User should modify the 3 lines below according to their needs
	var draw Draws = Dthroughput     // Define the diagrams to create (Dall = all diagrams)
	var conf []Confs = []Confs{Call} // Define the configs to process (Call = all configs)
	fileNb := 4                      // Numero of the file to display as histo sample (-1 = all files in the config)

	// Do NOT modify below
	if draw == Dcompare {
		compareConfigs(conf)
	} else {
		switch conf[0] {
		case Call:
			// Process all configs
			for _, c := range Configs {
				drawConfig(c, draw, fileNb)
			}
		default:
			// Process only the configs defined
			for _, c := range conf {
				drawConfig(Configs[c-1], draw, fileNb)
			}
		}
	}
}

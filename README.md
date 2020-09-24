 ## A. Aiming

1. Testing an application against a Kafka cluster, plot the test results in various ways.

* Draw the defined configs (C...) in the same diagram for comparison	
* Draw the raw points
* Draw the equivalent histogram
* Draw the computed mean
* Draw the computed mean with error deviations
* Draw a sliding window accross the points
* Draw the throughput
* Draw the number of messages per second

2. Compute the distribution moments (mean, standard and absolute deviations, skewness, curtosis)

3. Interpolate the curves with gaussian or linear regressions or polynoms of any degree.

4. Save the diagrams in PNG format

## B. Usage

1. Fill the configuration _Configs_ item in configs.go
1. Set the _draw_, _conf_ parameters
1. go run configs.go drawings.go

You may set the variable _PRINT_ to false to NOT display the moments while computing them for each diagram.

You can vary the size of the window in the sliding diagrams with the parameter _NVAL_

## C. Examples
1. ### Dcompare 
comparison beween 3 configs (CfetchMinBytes_100k, CfetchMinBytes_300k, CfetchMinBytes_3000k)

 		var draw Draws = Dcompare
		var conf []Confs = []Confs{CfetchMinBytes_100k, CfetchMinBytes_300k, CfetchMinBytes_3000k}

2. ### fileNb = -1 
create all diagrams for CqueueBufMaxMsg_ms100_30k (with all histo)

		var draw Draws = Dall
		var conf []Confs = []Confs{CqueueBufMaxMsg_ms100_30k}
		fileNb := -1

3. ### Dall + fileNb >= 0
create only interesting diagrams for CqueueBufMaxMsg_ms100_30k (only 1 histo)

		var draw Draws = Dall
		var conf []Confs = []Confs{CqueueBufMaxMsg_ms100_30k}
		fileNb := 4

4. ### Call
create the same diagram (here throughputs) for all configs

		var draw Draws = Dthroughput
		var conf []Confs = []Confs{Call}
		fileNb := 4

5. ### Call + Dall
create all diagrams for all configs (big number of generated files)

6. ### Call + Dall + fileNb=-1
create all diagrams for all configs with all histo (huge number of generated files)

## D. External libraries

* gonum.org/v1/plot/plotter and gonum.org/v1/plot/vg used to draw
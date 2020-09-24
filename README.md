 Define the actions to do hereafter.
 Ex1 Dconpare => comparison beween 3 configs (CfetchMinBytes_100k, CfetchMinBytes_300k, CfetchMinBytes_3000k)
 		var draw Draws = Dcompare
		var conf []Confs = []Confs{CfetchMinBytes_100k, CfetchMinBytes_300k, CfetchMinBytes_3000k}
 Ex2 fileNb = -1 => create all diagrams for CqueueBufMaxMsg_ms100_30k (with all histo)
		var draw Draws = Dall
		var conf []Confs = []Confs{CqueueBufMaxMsg_ms100_30k}
		fileNb := -1
 Ex3 Dall + fileNb >= 0 => create only interesting diagrams for CqueueBufMaxMsg_ms100_30k (only 1 histo)
		var draw Draws = Dall
		var conf []Confs = []Confs{CqueueBufMaxMsg_ms100_30k}
		fileNb := 4
 Ex4 Call => create the same diagram (here throughputs) for all configs
		var draw Draws = Dthroughput
		var conf []Confs = []Confs{Call}
		fileNb := 4
 Ex5 Call + Dall => create all diagrams for all configs (big number of generated files)
 Ex5 Call + Dall + fileNb=-1 => create all diagrams for all configs with all histo (huge number of generated files)
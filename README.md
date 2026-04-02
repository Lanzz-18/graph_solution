# graph_solution
A program which processes csv files contained directed graph data and display properties from it 

## Requirements
Should have Go installed in your system. Download: https://golang.org/dl/

## How to run
Run on Git bash
```bash
chmod +x graph_solution
./graph_solution path/to/graph.csv
```

## Example
There are three graph csv files (graph1.csv, graph2.csv, graph3.csv)
```bash
./graph_solution test_cases/graph1.csv
```

## Output
```
is_dag: true/false
max_in_degree: <integer>
max_out_degree: <integer>
pr_max: <float>
pr_min: <float>
```

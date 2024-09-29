# Byzantine faults

## Goal
This project is to understand the use of GoRountines and what are the different methods of accessing data from concurrent goroutines.
The algorithm used where **locks** and **Compare And Swap** based.

The final goal was to benchmark each algorithm by changing the szie of input data, and representing the values in an excel file.

## Run
To modifiy the amount of data used by the goroutines to benchmark each algorithm, modify the lines :
```go
jobs := []int{500000, 5000000, 7000000, 50000, 50000}
store_sizes := []int{20, 20, 50, 100, 120}
```
*Keep an equal number of jobs and store sizes to obtain a properly organised graph*

To run the project, enter the following command:
`go run *.go`
or `go build` then `go benchmark`.
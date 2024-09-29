package main 

import ( 
	"fmt" 
	"sync" 
	"time" 
	"github.com/xuri/excelize/v2"
) 

type Store interface { 
	Name() string 
	Size() int 
	At(index int) int64 
	Add(index int, amount int64) 
	Substract(index int, amount int64) 
} 

type Runner struct { 
	THREAD_COUNT 	int 
	JOB_COUNT		int 
	STORE_SIZE		int 
	executor		sync.WaitGroup 
	jobs 			[]func()
}  

func (r *Runner) bench(store Store) time.Duration { 
	for i := 0; i < r.JOB_COUNT; i++ { 
		adding := i%2 == 0 
		r.jobs = append(r.jobs, func() { 
			for j := 0; j < store.Size(); j++ { 
				if adding { 
					store.Add(j, 1) 
				} else { 
					store.Substract(j, 1) 
				} 
			} 
			r.executor.Done() 
		}) 
	} 
	start := time.Now() 
	for _, job := range r.jobs { 
		r.executor.Add(1) 
		go job() 
	} 
	r.executor.Wait() 
	stop := time.Now() 
	
	fmt.Println(store.Name(), ":", stop.Sub(start))
	
	for i := 0; i < store.Size(); i++ { 
		fmt.Print(store.At(i), " ") 
	}
		
	fmt.Println("")
	return stop.Sub(start) 
} 

func main() {
	jobs := []int{500000, 5000000, 7000000, 50000, 50000}
	store_sizes := []int{20, 20, 50, 100, 120}
	
	file := excelize.NewFile()
	sheetName := "Benchmark"

	index, _ := file.NewSheet(sheetName)
	file.SetActiveSheet(index)

	headers := []string{"Jobs", "Store_size", "HippieStore Avg(ms)", "PessimisticStore Avg(ms)", "OptimisticStore Avg(ms)"}

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		file.SetCellValue(sheetName, cell, header)
	}

	row := 2
	for i := 0; i < len(jobs); i++ {
		fmt.Println("Testing with", jobs[i], "jobs and a size of", store_sizes[i])
		
		// 1. HippieStore benchmark
		var hippieTotalTime int64 = 0
		for j := 0; j < 10; j++ {
			runner := Runner{
				THREAD_COUNT: 4,
				JOB_COUNT:    jobs[i],
				STORE_SIZE:   store_sizes[i],
			}
			elapsedTime := runner.bench(NewHippieStore(runner.STORE_SIZE))
			hippieTotalTime += elapsedTime.Milliseconds()
		}
		hippieAvgTime := hippieTotalTime / 10

		// 2. PessimisticStore benchmark
		var pessimisticTotalTime int64 = 0
		for k := 0; k < 10; k++ {
			runner := Runner{
				THREAD_COUNT: 4,
				JOB_COUNT:    jobs[i],
				STORE_SIZE:   store_sizes[i],
			}
			elapsedTime := runner.bench(NewPessimisticStore(runner.STORE_SIZE))
			pessimisticTotalTime += elapsedTime.Milliseconds()
		}
		pessimisticAvgTime := pessimisticTotalTime / 10
		
		// 3. OptimisticStore benchmark
		var optimisticTotalTime int64 = 0
		for l := 0; l < 10; l++ {
			runner := Runner{
				THREAD_COUNT: 4,
				JOB_COUNT:    jobs[i],
				STORE_SIZE:   store_sizes[i],
			}
			elapsedTime := runner.bench(NewOptimisticStore(runner.STORE_SIZE))
			optimisticTotalTime += elapsedTime.Milliseconds()
		}
		optimisticAvgTime := optimisticTotalTime / 10
		
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), jobs[i])
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), store_sizes[i])
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), hippieAvgTime)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", row), pessimisticAvgTime)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", row), optimisticAvgTime)
		
		row++
	}
	
	if err := addSimpleLineChart(file, sheetName, row-1); err != nil {
		fmt.Println("Error adding chart:", err)
	}

	if err := file.SaveAs("benchmark.xlsx"); err != nil {
		fmt.Println("Error saving file:", err)
	} else {
		fmt.Println("Benchmark results saved to benchmark.xlsx")
	}
}

func addSimpleLineChart(f *excelize.File, sheet string, lastRow int) error {
	categories := fmt.Sprintf("%s!$A$2:$A$%d", sheet, lastRow)
	hippieValues := fmt.Sprintf("%s!$C$2:$C$%d", sheet, lastRow)
	pessimisticValues := fmt.Sprintf("%s!$D$2:$D$%d", sheet, lastRow)
	optimisticValues := fmt.Sprintf("%s!$E$2:$E$%d", sheet, lastRow)

	lineChart := excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{
			{
				Name:       "HippieStore",
				Categories: categories,
				Values:     hippieValues,
			},
			{
				Name:       "PessimisticStore",
				Categories: categories,
				Values:     pessimisticValues,
			},
			{
				Name:       "OptimisticStore",
				Categories: categories,
				Values:     optimisticValues,
			},
		},
		Format: excelize.GraphicOptions{
			OffsetX: 15,
			OffsetY: 10,
		},
		Legend: excelize.ChartLegend{
			Position: "bottom",
		},
	}

	return f.AddChart(sheet, "G2", &lineChart)
}


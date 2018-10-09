// parse string by line put into array
/*
in rls run
./pa1Test pa1.go to test

*/
package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"container/ring"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type process struct {
	name       string
	input      int64 // order it was input
	wait       int64
	turnaround int64
	arrival    int64 // time incrament when process arrives
	burst      int64 // how many time units it will take to compleet
	index      int   // only used by heap interface
}

// interface for the priotity queue for sjf that implements a heap interface
type PriorityQueueBurst []*process

func (pq PriorityQueueBurst) Len() int { return len(pq) }

// want shortest jobs on top
func (pq PriorityQueueBurst) Less(i, j int) bool {
	if pq[i].burst < pq[j].burst {
		return true
	} else if pq[i].burst > pq[j].burst {
		return false
	} else {
		return pq[i].arrival < pq[j].arrival
	}
}

func (pq PriorityQueueBurst) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueueBurst) Push(x interface{}) {
	n := len(*pq)
	item := x.(process)
	item.index = n
	*pq = append(*pq, &item)
}

func (pq *PriorityQueueBurst) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[0]
	item.index = -1 // for safety
	*pq = old[1:n]
	return item
}

// return top element without removing it
func (pq *PriorityQueueBurst) Peak() interface{} {
	old := *pq
	//n := len(old)
	item := old[0]
	return item
}

// priority is burst. smaller = higher priority
func (pq *PriorityQueueBurst) Update(item *process, priority int64) {
	item.burst = priority
	heap.Fix(pq, item.index)
}

// for sorting process array by input field
type byInput []process

func (p byInput) Len() int {
	return len(p)
}
func (p byInput) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p byInput) Less(i int, j int) bool {
	return p[i].input < p[j].input
}

// for sorting process array by arrival field
type byArrival []process

func (p byArrival) Len() int {
	return len(p)
}
func (p byArrival) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p byArrival) Less(i int, j int) bool {
	return p[i].arrival < p[j].arrival
}

type sim struct {
	schAlg       string    // which scheduling alg the input file says to use
	processCount int64     // number of processes
	runFor       int64     //
	quant        int64     // only used for round robbin, -1 otherwise
	processes    []process //
}

func errorFunc(eCode int, flavor string, line int) {
	/*
		switch eCode {
		case 1:
			{

			}
	*/
	fmt.Printf("Error on line %d: %s\n", line, flavor)
	os.Exit(eCode)
	//}
}

func fcfs(info sim) string {
	var buffer bytes.Buffer

	buffer.WriteString("  ")
	buffer.WriteString(strconv.FormatInt(info.processCount, 10))
	buffer.WriteString(" processes\n")

	buffer.WriteString("Using First-Come First-Served\n")

	var time int64 = 0
	//var selected int = 0
	for time = 0; time < info.runFor; time++ {

	}

	return buffer.String()
}

/*
This function takes in a file and returns a string
Paramiters File
Returns string
*/
func fileToStr(file *os.File) string {
	scanner := bufio.NewScanner(file)
	txt := ""
	for scanner.Scan() {
		txt += scanner.Text() + "\n"

	}
	return txt
}

/**

 */
func isSchAlg(str string) bool {
	return str == "fcfs" || str == "sjf" || str == "rr"
}

/*
this function takes in an array of strings and puts pertenate info into a sim struct
Paramiters: lines: an array of stings representing each line
Retruns: a sim struct with all the info needed to run a simulation of the processes
*/
func parse(lines []string) sim {
	// this may be a mess
	// need to write and plan this out
	//info :=
	schAlg := ""               // which scheduling alg the input file says to use
	var processCount int64 = 0 // number of processes
	var runFor int64 = 0       //
	var quant int64 = -1       // only used for round robbin, -1 otherwise
	var processes []process    //
	var pIndex int64 = 0       // used to index processes
	done := false
	var index = 0

	for index, line := range lines {
		if done {
			break
		}

		//fmt.Printf("%d: ", index)
		state := 0
		// tokenize line based on white space
		tokens := strings.FieldsFunc(line, Split)
		//token, status := getNextToken(line)

		for _, token := range tokens {
			if token[0] == '#' || done {
				break
			}
			switch state {
			case -1:
				errorFunc(3, token, index)

			case 0: // start
				if token == "processcount" {
					state = 1
				} else if token == "runfor" {
					state = 3
				} else if token == "use" {
					state = 5
				} else if token == "quantum" {
					state = 9
				} else if token == "process" {
					state = 11
					if pIndex >= processCount {
						errorFunc(8, token, index)
						done = true
					}
				} else if token == "end" {
					done = true
					break
				} else { // unknown keyword
					errorFunc(4, token, index)
				}
				break

			case 1: // processCount
				num, err := strconv.ParseInt(token, 10, 64)
				if err == nil {
					processCount = num
					processes = make([]process, processCount)
					state = -1
				} else {
					// could not parse int
					errorFunc(1, token, index)
					done = true

				}
				break

			case 3: // runfor
				num, err := strconv.ParseInt(token, 10, 64)
				if err == nil {
					runFor = num
					state = -1
				} else {
					errorFunc(1, token, index)
					done = true
				}
				break

			case 5: // use
				if isSchAlg(token) {
					schAlg = token
					state = -1
				} else {
					errorFunc(2, token, index)
					done = true
				}
				break

			case 9: // quantum

				num, err := strconv.ParseInt(token, 10, 64)
				if err == nil {

					quant = num
					//fmt.Printf("\n%d\n", quant)
					state = -1
				} else {
					errorFunc(1, token, index)
					done = true
				}
				break

			case 11: //process
				if token == "name" {
					state = 12
				} else {
					errorFunc(5, token, index)
					done = true
				}

				break

			case 12: // name
				processes[pIndex].name = token
				processes[pIndex].input = pIndex
				state = 13

				break

			case 13: //
				if token == "arrival" {
					state = 14
				} else {
					errorFunc(6, token, index)
					done = true
				}
				break

			case 14: //arrival
				num, err := strconv.ParseInt(token, 10, 64)
				if err == nil {
					processes[pIndex].arrival = num
					state = 15
				} else {
					errorFunc(1, token, index)
					done = true
				}
				break

			case 15: //
				if token == "burst" {
					state = 16
				} else {
					errorFunc(7, token, index)
					done = true
				}
				break

			case 16: // burst
				num, err := strconv.ParseInt(token, 10, 64)
				if err == nil {
					processes[pIndex].burst = num
					pIndex++
					state = -1
				} else {
					errorFunc(1, token, index)
					done = true
				}
				break

			}
			//fmt.Printf("%s ", token)
		}
		//fmt.Printf("\n")

	}

	if schAlg == "" {
		errorFunc(9, "", index)
	} else if runFor <= 0 {
		errorFunc(10, "", index)
	} else if schAlg == "rr" && quant < 0 {
		errorFunc(11, "", index)
	} else if pIndex != processCount {
		errorFunc(12, "", index)
	}

	return sim{schAlg: schAlg, processCount: processCount, runFor: runFor, quant: quant, processes: processes}

}

/*
This function simulates the cpu running the process
Paramiters: lines: an array? of stings containing the proceeses to simulate and how to sumulate
Returns: a string to be writen to an output file
*/
func run(lines []string) string {

	info := parse(lines)

	// fmt.Println(info)
	sort.Sort(byArrival(info.processes))

	// the meat
	var output string = ""
	output += fmt.Sprintf("%3d processes\n", info.processCount)
	// determine sch alg
	if info.schAlg == "rr" {
		output += "Using Round-Robin\n"
		output += fmt.Sprintf("Quantum %3d\n\n", info.quant)
		output += rr(info, info.quant)
	} else if info.schAlg == "fcfs" {
		output += "Using First-Come First-Served\n"
		output += rr(info, info.quant)
	} else if info.schAlg == "sjf" {
		output += "Using preemptive Shortest Job First\n"
		output += sjf(info)
	}

	return output
}

func rr(info sim, quant int64) string {
	var buffer bytes.Buffer

	var time int64 = 0
	var fin []process = make([]process, info.processCount)
	var finIndex int64 = 0
	sch := ring.New(1) //int(info.processCount))
	sch = nil
	end := sch // where new elements will be added
	var arrivalIndex int64 = 0

	var switchTimer int64 = 0 // used to keep trac of when to switch

	for time = 0; time < info.runFor; time++ {

		str := fmt.Sprintf("Time %3d :", time)
		//buffer.WriteString(str)

		// check for new arrivals
		if arrivalIndex < info.processCount && time == info.processes[arrivalIndex].arrival {
			line := fmt.Sprintf("%s %s arrived\n", str, info.processes[arrivalIndex].name)
			buffer.WriteString(line)
			// if first arival put in sch
			if sch == nil {
				sch = ring.New(1)
				sch.Value = info.processes[arrivalIndex]
				//end = sch

				line := fmt.Sprintf("%s %s selected (burst %3d)\n", str, info.processes[arrivalIndex].name, info.processes[arrivalIndex].burst)
				buffer.WriteString(line)
				arrivalIndex++
				//continue // start computation on next clock cycle
			} else {
				// splice in another ring
				new := ring.New(1)
				new.Value = info.processes[arrivalIndex]
				// but it behing of current
				end = sch.Move(-1)
				end = end.Link(new)
				arrivalIndex++
			}
		}

		// check if current process is compleate. if burst = 0 remove sch. val from list
		if sch != nil && sch.Value.(process).burst == 0 {
			val := sch.Value.(process)
			line := fmt.Sprintf("%s %s finished\n", str, val.name)
			buffer.WriteString(line)

			// save turnaround and wait
			fin[int(finIndex)] = val
			finIndex++

			// set sch to defined rig zero value
			if sch.Len() == 1 {
				sch = nil
			} else {
				// lop off current element
				sch = sch.Unlink(sch.Len() - 1)
				line := fmt.Sprintf("%s %s selected (burst %3d)\n", str, sch.Value.(process).name, sch.Value.(process).burst)
				buffer.WriteString(line)
				switchTimer = 0
			}
		}

		// if items in scheduling queue simulate compute
		if sch != nil {
			// swithch if time to switch
			if quant != -1 {
				if switchTimer >= quant {
					switchTimer = 0
					//temp := sch
					sch = sch.Next()

					//if sch.Value != temp.Value {
					line := fmt.Sprintf("%s %s selected (burst %3d)\n", str, sch.Value.(process).name, sch.Value.(process).burst)
					buffer.WriteString(line)
					//}

				}
			}

			val := sch.Value.(process)
			val.burst -= 1

			switchTimer++

			// incrament all proceeses' turnaround
			val.turnaround++
			temp := sch.Next()
			for temp.Value != sch.Value {
				p := temp.Value.(process)
				p.turnaround++
				temp.Value = p
				temp = temp.Next()

			}

			// incrament all proceeses' wait exep sch
			temp = sch.Next()
			for temp.Value != sch.Value {
				p := temp.Value.(process)

				p.wait++
				temp.Value = p
				temp = temp.Next()

			}

			// save changes made
			sch.Value = val

		} else {
			line := fmt.Sprintf("%s Idle\n", str)
			buffer.WriteString(line)
		}
	}

	str := fmt.Sprintf("Finished at time %3d\n\n", time)
	buffer.WriteString(str)

	sort.Sort(byInput(fin))
	for _, p := range fin {
		if p.name != "" {
			line := fmt.Sprintf("%s wait %3d turnaround %3d\n", p.name, p.wait, p.turnaround)
			buffer.WriteString(line)
		}
	}

	return buffer.String()
}

/**/
func Split(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t' || r == '\r' //|| r == '\n\r'
}

func sjf(info sim) string {
	var buffer bytes.Buffer

	var time int64 = 0
	var fin []process = make([]process, info.processCount)
	var finIndex int64 = 0
	sch := make(PriorityQueueBurst, 0)
	//sch = nil
	heap.Init(&sch)

	var arrivalIndex int64 = 0

	var current interface{} // holds the current

	//	var switchTimer int64 = 0 // used to keep trac of when to switch

	for time = 0; time < info.runFor; time++ {

		str := fmt.Sprintf("Time %3d :", time)
		//buffer.WriteString(str)

		// check for new arrivals
		if arrivalIndex < info.processCount && time == info.processes[arrivalIndex].arrival {
			line := fmt.Sprintf("%s %s arrived\n", str, info.processes[arrivalIndex].name)
			buffer.WriteString(line)
			old := current
			heap.Push(&sch, info.processes[arrivalIndex])
			heap.Init(&sch)
			current = sch.Peak()
			arrivalIndex++

			if old == nil {
				line = fmt.Sprintf("%s %s selected (burst %3d)\n", str, current.(*process).name, current.(*process).burst)
				buffer.WriteString(line)
			}
			// if premtion happens
			if old != nil && old.(*process).name != current.(*process).name { // comp by val
				line := fmt.Sprintf("%s %s selected (burst %3d)\n", str, current.(*process).name, current.(*process).burst)
				buffer.WriteString(line)

				// debug
				//fmt.Println(fmt.Sprintf("O: %s: %d (should be higher)", old.(*process).name, old.(*process).burst)) //
				// debug
				//fmt.Println(fmt.Sprintf("C: %s: %d\n", current.(*process).name, current.(*process).burst))

			}

		}

		// check if current process is compleate. if burst = 0 remove sch. val from list
		if sch.Len() > 0 && current.(*process).burst <= 0 {
			val := sch.Peak().(*process)
			line := fmt.Sprintf("%s %s finished\n", str, val.name)
			buffer.WriteString(line)
			// remove current top item from queue
			sch.Pop()
			heap.Init(&sch)

			// save turnaround and wait
			fin[int(finIndex)] = *val
			finIndex++

			// select next process
			if len(sch) > 0 {
				current = sch.Peak()
				line = fmt.Sprintf("%s %s selected (burst %3d)\n", str, current.(*process).name, current.(*process).burst)
				buffer.WriteString(line)
			}

		}

		// if items in scheduling queue simulate compute
		if len(sch) > 0 {

			//debug
			//

			// for each item in the priority queue
			for _, proc := range sch {

				//fmt.Println(proc)
				// increment wait time and turnaround
				proc.turnaround++
				proc.wait++

				// active process in queue
				if proc.name == current.(*process).name {

					proc.burst--
					proc.wait--
					//heap.Fix(&sch, index)

					// debug
					//fmt.Printf("Processing %s, %d remaining, len: %d\n", proc.name, proc.burst, sch.Len())
				}
			}
			//fmt.Println()

		} else {
			line := fmt.Sprintf("%s Idle\n", str)
			buffer.WriteString(line)
		}
	}

	str := fmt.Sprintf("Finished at time %3d\n\n", time)
	buffer.WriteString(str)

	sort.Sort(byInput(fin))
	for _, p := range fin {
		if p.name != "" {
			line := fmt.Sprintf("%s wait %3d turnaround %3d\n", p.name, p.wait, p.turnaround)
			buffer.WriteString(line)
		}
	}

	return buffer.String()
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: ./pa1 <input file> <output file>\n")
		return
	}

	inputFile := os.Args[1]

	outputFile := os.Args[2]

	file, errI := os.Open(inputFile)

	if errI != nil {
		file.Close()
		log.Fatal(errI)
	}

	// copy contents of file to string
	txt := fileToStr(file)
	file.Close()
	lines := strings.Split(txt, "\n")

	out, errO := os.Create(outputFile)
	if errO != nil {
		out.Close()
		log.Fatal(errO)
	}

	out.WriteString(run(lines))
	out.Close()
	//fmt.Printf("%s", output)

}

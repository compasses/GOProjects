package main

import (
	"algorithems/bubblesort"
	"algorithems/qsort"
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

var infile *string = flag.String("i", "infile", "Files contains value to sorting")
var outfile *string = flag.String("o", "outfile", "Files save sorting results")
var algorithems *string = flag.String("a", "qsort", "Sort algorithem")

func readValues(infile string) (values []int, err error) {
	file, err := os.Open(infile)
	if err != nil {
		fmt.Println("Fail to open file ", infile)
		return
	}
	defer file.Close()

	br := bufio.NewReader(file)
	values = make([]int, 0)

	for {
		line, isPrefix, err1 := br.ReadLine()
		if err1 != nil {
			if err1 != io.EOF {
				err = err1
			}
			break
		}

		if isPrefix {
			fmt.Println("A too long line, seems unexpected")
		}

		str := string(line)
		value, err1 := strconv.Atoi(str)
		if err1 != nil {
			err = err1
			return
		}
		values = append(values, value)
	}
	return
}

func writeValues(values []int, outfile string) error {
	file, err := os.Create(outfile)
	if err != nil {
		fmt.Println("Failed to created output file ", file)
		return err
	}
	defer file.Close()
	for _, value := range values {
		str := strconv.Itoa(value)
		file.WriteString(str + "\n")
	}
	return nil
}

func main() {
	flag.Parse()

	if infile != nil {
		fmt.Println("infile = ", *infile, "outfile = ", *outfile, "algorithem = ", *algorithems)
	}

	values, err := readValues(*infile)
	if err == nil {
		fmt.Println("Read values: ", values)
		t1 := time.Now()
		switch *algorithems {
		case "bubblesort":
			bubblesort.BubbleSort(values)
		case "qsort":
			qsort.QSort(values)
		default:
			fmt.Println("Sorting Algorithem ", *algorithems, "is either unknown or unsupported")
		}
		t2 := time.Now()
		fmt.Println("The Sorting process costs ", t2.Sub(t1), "to complete.")
		writeValues(values, *outfile)
	} else {
		fmt.Println(err)
	}
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// for sorting by key.
type ByKey []string

// for sorting by key.
func (a ByKey) Len() int      { return len(a) }
func (a ByKey) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool {
	x, _ := strconv.Atoi(a[i])
	y, _ := strconv.Atoi(a[j])
	return x < y
}

func ExternalMergeSort(filename string, memorymaxlines int) {
	l := getFileLength(filename)
	if l <= memorymaxlines {
		InternalSort(filename)
	} else {
		file1, file2 := spiltFile(filename, l)
		ExternalMergeSort(file1, memorymaxlines)
		ExternalMergeSort(file2, memorymaxlines)
		f1, err := os.Open(file1)
		if err != nil {
			return
		}
		f2, err := os.Open(file2)
		if err != nil {
			f1.Close()
			return
		}
		scanner1 := bufio.NewScanner(f1)
		scanner2 := bufio.NewScanner(f2)
		out, err := os.CreateTemp("", "out")
		if err != nil {
			f1.Close()
			f2.Close()
		}
		ok1 := scanner1.Scan()
		ok2 := scanner2.Scan()
		for {
			if ok1 && ok2 {
				line1 := scanner1.Text()
				line2 := scanner2.Text()
				num1, _ := strconv.Atoi(line1)
				num2, _ := strconv.Atoi(line2)
				if num1 < num2 {
					fmt.Fprintf(out, "%s\n", line1)
					ok1 = scanner1.Scan()
				} else {
					fmt.Fprintf(out, "%s\n", line2)
					ok2 = scanner2.Scan()
				}
			} else if !ok1 {
				fmt.Fprintf(out, "%s\n", scanner2.Text())
				for scanner2.Scan() {
					fmt.Fprintf(out, "%s\n", scanner2.Text())
				}
				break
			} else {
				fmt.Fprintf(out, "%s\n", scanner1.Text())
				for scanner1.Scan() {
					fmt.Fprintf(out, "%s\n", scanner1.Text())
				}
				break
			}
		}
		out.Close()
		os.Rename(out.Name(), filename)
	}
}

func InternalSort(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return
	}
	numbers := strings.Split(string(content), "\n")
	sort.Sort(ByKey(numbers))

	ofile, err := os.CreateTemp(".", "tmp")
	if err != nil {
		return
	}
	for _, line := range numbers {
		if line != "" {
			fmt.Fprintf(ofile, "%s\n", line)
		}
	}
	ofile.Close()
	os.Rename(ofile.Name(), filename)
}

func getFileLength(filename string) int {
	length := 0
	f, err := os.Open(filename)
	if err != nil {
		return length
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		length++
	}
	return length
}

func spiltFile(filename string, length int) (string, string) {
	f1, err := os.CreateTemp("", "file1")
	if err != nil {
		return "", ""
	}
	defer f1.Close()
	f2, err := os.CreateTemp("", "file2")
	if err != nil {
		return "", ""
	}
	defer f2.Close()
	f, err := os.Open(filename)
	if err != nil {
		return "", ""
	}
	defer os.Remove(filename)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	index := 0
	for scanner.Scan() {
		if index >= length/2 {
			fmt.Fprintf(f2, "%s\n", scanner.Text())
		} else {
			fmt.Fprintf(f1, "%s\n", scanner.Text())
		}
		index++
	}
	return f1.Name(), f2.Name()
}

func main() {
	for _, s := range []struct {
		File           string
		MemoryMaxLines int
	}{
		{"a.txt", 10},
		{"b.txt", 5},
		{"c.txt", 1},
	} {
		start := time.Now()
		ExternalMergeSort(s.File, s.MemoryMaxLines)
		fmt.Printf("MemoryMaxLines:%d Time:%v\n", s.MemoryMaxLines, time.Since(start))
	}
}

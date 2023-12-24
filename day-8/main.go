package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

func main() {
	puzzle := os.Args[1]

	path := os.Args[2]

	if puzzle == "1" {
		puzzle1(path)
	}

	if puzzle == "2" {
		fmt.Println(time.Now())
		puzzle2(path)
		fmt.Println(time.Now())
	}
}

func puzzle1(path string) {
	directions := parseDirections(path)

	_map := parseMap(path)

	count := 0

	pos := "AAA"

	for x := range repeatDirections(directions) {
		switch x {
		case "L":
			pos = _map[pos].Left()
		case "R":
			pos = _map[pos].Right()
		}

		count += 1

		if pos == "ZZZ" {
			break
		}
	}

	fmt.Println(count, pos)
}

type channels struct {
	In  chan<- string
	Out <-chan string
}

func startWorker(pos string, m map[string]LR) channels {
	in := make(chan string)
	out := make(chan string)

	go func() {
		for direction := range in {
			switch direction {
			case "L":
				pos = m[pos].Left()
			case "R":
				pos = m[pos].Right()
			}

			out <- pos
		}
	}()

	return channels{
		In:  in,
		Out: out,
	}
}

func puzzle2(path string) {
	m := parseMap(path)

	workers := map[string]channels{}

	nodes := startNodes(m)

	for i := range nodes {
		workers[nodes[i]] = startWorker(nodes[i], m)
	}

	count := 0

	for direction := range repeatDirections(parseDirections(path)) {
		for i := range nodes {
			workers[nodes[i]].In <- direction
		}

		z := true

		for i := range nodes {
			v := <-workers[nodes[i]].Out
			if v[2] != 'Z' {
				z = false
			}
		}

		count++

		if count%1000000 == 0 {
			fmt.Print(".")
		}

		if z {
			fmt.Println(count)
			return
		}
	}
}

func startNodes(m map[string]LR) []string {
	var nodes []string

	for x := range m {
		if x[2] == 'A' {
			nodes = append(nodes, x)
		}
	}

	return nodes
}

func repeatDirections(directions string) <-chan string {
	ch := make(chan string)
	go func() {
		for {
			for i := 0; i < len(directions); i++ {
				ch <- string(directions[i])
			}
		}
	}()
	return ch
}

func parseDirections(path string) string {
	f, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	scanner.Scan()

	_ = f.Close()

	return scanner.Text()
}

func parseMap(path string) map[string]LR {
	f, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	scanner.Scan()

	scanner.Scan()

	re := regexp.MustCompile("([A-Z0-9]{3}) = \\(([A-Z0-9]{3}), ([A-Z0-9]{3})\\)")

	result := map[string]LR{}

	for scanner.Scan() {
		line := scanner.Text()

		submatch := re.FindStringSubmatch(line)

		result[submatch[1]] = LR{submatch[2], submatch[3]}
	}

	_ = f.Close()

	return result
}

type LR [2]string

func (x LR) Left() string {
	return x[0]
}

func (x LR) Right() string {
	return x[1]
}

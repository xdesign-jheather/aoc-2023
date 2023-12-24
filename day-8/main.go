package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	puzzle := os.Args[1]

	path := os.Args[2]

	if puzzle == "1" {
		puzzle1(path)
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

	re := regexp.MustCompile("([A-Z]{3}) = \\(([A-Z]{3}), ([A-Z]{3})\\)")

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

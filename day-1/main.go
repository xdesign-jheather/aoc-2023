package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var words = map[string]string{
	"one":   "1",
	"two":   "2",
	"three": "3",
	"four":  "4",
	"five":  "5",
	"six":   "6",
	"seven": "7",
	"eight": "8",
	"nine":  "9",
}

func state() (func(string), func() int) {
	var first string
	var last string

	f1 := func(x string) {
		if first == "" {
			first = x
		}

		last = x
	}

	f2 := func() int {
		n, err := strconv.Atoi(fmt.Sprintf("%s%s", first, last))

		if err != nil {
			panic(err)
		}

		return n
	}

	return f1, f2
}

func number(line string) int {
	set, get := state()

	for i := 0; i < len(line); i++ {
		if line[i] < '0' || line[i] > '9' {
			continue
		}

		set(string(line[i]))
	}

	return get()
}

func numberWithWords(line string) int {
	extended := line + "     "

	set, get := state()

	for i := 0; i < len(line); i++ {
		char := string(line[i])

		if line[i] >= '0' && line[i] <= '9' {
			set(char)
			continue
		}

		for word, value := range words {
			sub := extended[i : i+len(word)]

			if sub == word {
				set(value)
				break
			}
		}
	}

	return get()
}

func puzzle1(path string) {
	r, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	total := 0

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		total += number(line)
	}

	fmt.Println(total)
}

func puzzle2(path string) {
	r, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	total := 0

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		total += numberWithWords(line)
	}

	fmt.Println(total)
}

func main() {
	path := os.Args[1]

	puzzle1(path)

	puzzle2(path)
}

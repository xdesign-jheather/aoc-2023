package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	Red   = "red"
	Blue  = "blue"
	Green = "green"
)

func main() {
	path := os.Args[1]

	puzzle1(path)

	puzzle2(path)
}

type Game struct {
	ID    int
	Hands Hands
}

type Hands []Hand

type Hand struct {
	Colour string
	Count  int
}

func parse(line string) Game {
	line = strings.Replace(line, ";", ",", -1)
	line = strings.Replace(line, ":", ",", -1)
	split := strings.Split(line[5:], ",")

	id, err := strconv.Atoi(split[0])

	if err != nil {
		panic(err)
	}

	var hands Hands

	for i := 1; i < len(split); i++ {
		x := strings.Split(split[i], " ")

		iCount, err := strconv.Atoi(x[len(x)-2])

		if err != nil {
			panic(err)
		}

		hands = append(hands, Hand{
			Colour: x[len(x)-1],
			Count:  iCount,
		})
	}

	return Game{
		ID:    id,
		Hands: hands,
	}
}

func puzzle1(path string) {
	total := 0

	for game := range games(path) {
		if possible(game) {
			total += game.ID
		}
	}

	fmt.Println(total)
}

func puzzle2(path string) {
	total := 0

	for game := range games(path) {
		r, g, b := minimum(game)
		total += r * g * b
	}

	fmt.Println(total)
}

func games(path string) chan Game {
	ch := make(chan Game)

	go func() {
		for line := range lines(path) {
			ch <- parse(line)
		}

		close(ch)
	}()

	return ch
}

func lines(path string) chan string {
	r, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	ch := make(chan string)

	go func() {
		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			line := scanner.Text()

			if line == "" {
				continue
			}

			ch <- line
		}

		close(ch)
	}()

	return ch
}

func possible(game Game) bool {
	for _, hand := range game.Hands {
		if hand.Colour == Red && hand.Count > 12 {
			return false
		}
		if hand.Colour == Green && hand.Count > 13 {
			return false
		}
		if hand.Colour == Blue && hand.Count > 14 {
			return false
		}
	}
	return true
}

func minimumState() (func(int), func(int), func(int), func() (int, int, int)) {
	r, g, b := 0, 0, 0

	return func(i int) {
			if i > r {
				r = i
			}
		}, func(i int) {
			if i > g {
				g = i
			}
		}, func(i int) {
			if i > b {
				b = i
			}
		},
		func() (int, int, int) {
			return r, g, b
		}
}

func minimum(game Game) (int, int, int) {
	r, g, b, result := minimumState()

	for _, hand := range game.Hands {
		if hand.Colour == Red {
			r(hand.Count)
		}
		if hand.Colour == Green {
			g(hand.Count)
		}
		if hand.Colour == Blue {
			b(hand.Count)
		}
	}

	return result()
}

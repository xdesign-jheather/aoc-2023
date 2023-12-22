package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var values = map[int32]int{
	'2': 0,
	'3': 1,
	'4': 2,
	'5': 3,
	'6': 4,
	'7': 5,
	'8': 6,
	'9': 7,
	'T': 8,
	'J': 9,
	'Q': 10,
	'K': 11,
	'A': 12,
}

type Type int

const (
	HighCard Type = iota
	OnePair
	TwoPair
	ThreeOfKind
	FullHouse
	FourOfKind
	FiveOfKind
)

func main() {
	puzzle := os.Args[1]

	path := os.Args[2]

	if puzzle == "1" {
		puzzle1(path)
	}
}

func puzzle1(path string) {
	hands, _ := parse(path)

	for _, h := range hands {
		if h.Type() == FiveOfKind {
			fmt.Println(h, h.Value())
		}
	}
}

func parse(path string) (Hands, []int) {
	f, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	var hands Hands

	var bids []int

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")

		hands = append(hands, Hand(parts[0]))

		bid, _ := strconv.Atoi(parts[1])

		bids = append(bids, bid)
	}

	_ = f.Close()

	return hands, bids
}

type Hand string

func (h Hand) Type() Type {
	cards := map[int32]int{}

	for _, char := range h {
		cards[char] += 1
	}

	same := map[int]int{}

	for _, count := range cards {
		same[count] += 1
	}

	switch {
	case same[5] == 1:
		return FiveOfKind
	case same[4] == 1:
		return FourOfKind
	case same[3] == 1 && same[2] == 1:
		return FullHouse
	case same[3] == 1:
		return ThreeOfKind
	case same[2] == 2:
		return TwoPair
	case same[2] == 1:
		return OnePair
	default:
		return HighCard
	}
}

func (h Hand) Value() int {
	return 0
}

type Hands []Hand

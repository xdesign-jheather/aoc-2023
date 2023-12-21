package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

func main() {
	path := os.Args[1]

	puzzle1(path)
}

type Games []Game

type Game struct {
	Winning []string
	Have    []string
}

func (g *Game) Won() []string {
	have := map[string]bool{}

	for _, num := range g.Have {
		have[num] = true
	}

	var won []string

	for _, num := range g.Winning {
		if have[num] {
			won = append(won, num)
		}
	}

	return won
}

func (g *Game) Score() int {
	won := g.Won()

	score := 0

	for range won {
		switch score {
		case 0:
			score = 1
		default:
			score *= 2
		}
	}

	return score
}

func puzzle1(path string) {
	f, _ := os.Open(path)

	total := 0

	for _, game := range parse(f) {
		score := game.Score()

		total += score
	}

	_ = f.Close()

	fmt.Println(total)
}

func parse(f io.Reader) Games {
	re1 := regexp.MustCompile("(: +| \\| +)")
	re2 := regexp.MustCompile(" +")

	var games Games

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		parts := re1.Split(line, -1)

		games = append(games, Game{
			Winning: re2.Split(parts[1], -1),
			Have:    re2.Split(parts[2], -1),
		})
	}

	return games
}

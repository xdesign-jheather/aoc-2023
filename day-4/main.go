package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

func main() {
	puzzle := os.Args[1]

	path := os.Args[2]

	if puzzle == "1" {
		puzzle1(path)
	}

	if puzzle == "2" {
		puzzle2(path)
	}
}

type Games []Game

type Game struct {
	ID      string
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

func counter() (func(string), func() int) {
	count := map[string]int{}

	return func(s string) {
			count[s] += 1
		}, func() int {
			total := 0
			for _, n := range count {
				total += n
			}
			return total
		}
}

func puzzle2(path string) {
	f, _ := os.Open(path)

	inc, result := counter()

	games := parse(f)

	type recurse func(i int, f recurse)

	walk := recurse(func(i int, f recurse) {
		inc(games[i].ID)

		wins := len(games[i].Won())

		for off := 1; off <= wins; off++ {
			f(i+off, f)
		}
	})

	for i := range games {
		walk(i, walk)
	}

	fmt.Println(result())

	_ = f.Close()
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
			ID:      parts[0],
			Winning: re2.Split(parts[1], -1),
			Have:    re2.Split(parts[2], -1),
		})
	}

	return games
}

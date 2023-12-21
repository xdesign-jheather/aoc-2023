package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	path := os.Args[1]

	puzzle1(path)
}

type Games []Game

type Game struct {
	ID      int
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

		fmt.Printf("%+v\n", game)
		fmt.Println(len(game.Have), len(game.Winning))
		fmt.Println("Won", game.Won())
		fmt.Println("Score", score)

		total += score
	}

	_ = f.Close()

	fmt.Println(total)
}

func parse(f io.Reader) Games {
	re := regexp.MustCompile(" +")

	var games Games

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		line = strings.Replace(line, "Card ", "", -1)
		line = strings.Replace(line, ":", " |", -1)
		parts := strings.Split(line, " | ")

		fmt.Println(parts)

		id, _ := strconv.Atoi(parts[0])

		games = append(games, Game{
			ID:      id,
			Winning: re.Split(parts[1], -1),
			Have:    re.Split(parts[2], -1),
		})
	}

	return games
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	puzzle := os.Args[1]

	path := os.Args[2]

	if puzzle == "1" {
		puzzle1(path)
	}
}

func puzzle1(path string) {
	races := parseRaces(path)

	product := 1

	for _, race := range races {
		fmt.Printf("%+v\n", race)

		product *= len(race.Breakers())
	}

	fmt.Println(product)
}

func parseRaces(path string) []Race {
	f, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile("[0-9]+")

	scanner := bufio.NewScanner(f)

	scanner.Scan()

	times := re.FindAllString(scanner.Text(), -1)

	scanner.Scan()

	distances := re.FindAllString(scanner.Text(), -1)

	if len(times) != len(distances) {
		panic("times and distance count off")
	}

	var races []Race

	for i := 0; i < len(times); i++ {
		time, _ := strconv.Atoi(times[i])
		distance, _ := strconv.Atoi(distances[i])

		races = append(races, Race{
			Time:     time,
			Distance: distance,
		})
	}

	_ = f.Close()

	return races
}

type Race struct {
	Time     int
	Distance int
}

func (r *Race) Speed(button int) int {
	return button * (r.Time - button)
}

func (r *Race) Strategies() [][2]int {
	var x [][2]int
	for button := 0; button <= r.Time; button++ {
		x = append(x, [2]int{button, r.Speed(button)})
	}
	return x
}

func (r *Race) Breakers() [][2]int {
	var x [][2]int
	for _, strategy := range r.Strategies() {
		if strategy[1] > r.Distance {
			x = append(x, strategy)
		}
	}
	return x
}

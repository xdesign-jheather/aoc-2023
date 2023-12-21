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
	puzzle := os.Args[1]

	path := os.Args[2]

	if puzzle == "1" {
		puzzle1(path)
	}

	if puzzle == "2" {
		puzzle2(path)
	}
}

func seedRanges(data []int) <-chan int {
	ch := make(chan int)

	go func() {
		for i := 0; i < len(data); i += 2 {
			start, length := data[i], data[i+1]

			for x := 0; x < length; x++ {
				ch <- start + x
			}
		}
		close(ch)
	}()

	return ch
}

func puzzle1(path string) {
	f1, _ := os.Open(path)
	f2, _ := os.Open(path)

	mappers := parseMappers(f1)
	allSeeds := parseSeeds(f2)

	_ = f1.Close()
	_ = f2.Close()

	lowest := -1

	for _, seed := range allSeeds {

		fmt.Println("Start with seed", seed)

		location := mappers.Resolve(seed)

		switch {
		case lowest == -1:
			lowest = location
		case location < lowest:
			lowest = location
		}

		fmt.Println("End with location", location)

		fmt.Println(".")
	}

	fmt.Println(lowest)
}

func puzzle2(path string) {
	f1, _ := os.Open(path)
	f2, _ := os.Open(path)

	mappers := parseMappers(f1)
	allSeeds := seedRanges(parseSeeds(f2))

	_ = f1.Close()
	_ = f2.Close()

	lowest := -1

	for seed := range allSeeds {
		location := mappers.Resolve(seed)

		switch {
		case lowest == -1:
			lowest = location
		case location < lowest:
			lowest = location
		}
	}

	fmt.Println(lowest)
}

func mapperFactory() (func(string), func() Mappers) {
	var mappers Mappers
	var build *Mapper

	return func(line string) {
			switch {
			case strings.HasSuffix(line, "map:") && build == nil:
				build = &Mapper{
					Name: strings.TrimSuffix(line, " map:"),
				}
			case build != nil && line == "":
				mappers = append(mappers, build)
				build = nil
			case line == "":
			default:
				parts := strings.Split(line, " ")

				if len(parts) != 3 {
					panic(parts)
				}

				n1, _ := strconv.Atoi(parts[0])
				n2, _ := strconv.Atoi(parts[1])
				n3, _ := strconv.Atoi(parts[2])

				build.Mappings = append(build.Mappings, &Mapping{
					n1, n2, n3,
				})
			}

		}, func() Mappers {
			if build != nil {
				mappers = append(mappers, build)
				build = nil
			}
			return mappers
		}
}

func parseSeeds(f io.Reader) []int {
	scanner := bufio.NewScanner(f)

	scanner.Scan()

	line := scanner.Text()

	return seeds(line)
}

func parseMappers(f io.Reader) Mappers {
	scanner := bufio.NewScanner(f)

	scanner.Scan()

	build, results := mapperFactory()

	for scanner.Scan() {
		build(scanner.Text())
	}

	return results()
}

func seeds(line string) []int {
	re := regexp.MustCompile(":? +")

	parts := re.Split(line, -1)

	var nums []int

	for i := range parts {
		if i == 0 {
			continue
		}

		num, _ := strconv.Atoi(parts[i])

		nums = append(nums, num)
	}

	return nums
}

type Mappers []*Mapper

func (m Mappers) Resolve(number int) int {
	for _, mapper := range m {
		number = mapper.Resolve(number)
	}
	return number
}

type Mapper struct {
	Name     string
	Mappings Mappings
}

func (m *Mapper) Resolve(number int) int {
	for _, mapping := range m.Mappings {
		if ok, result := mapping.Resolve(number); ok {
			return result
		}
	}
	return number
}

type Mapping [3]int

func (m *Mapping) Resolve(number int) (bool, int) {
	dest, source, length := m[0], m[1], m[2]

	if number < source || number > source+length-1 {
		return false, 0
	}

	return true, number - (source - dest)
}

type Mappings []*Mapping

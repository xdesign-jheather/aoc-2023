package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

var adjacentDelta = [8][2]int{
	{-1, -1},
	{0, -1},
	{1, -1},
	{-1, 0},
	{1, 0},
	{-1, 1},
	{0, 1},
	{1, 1},
}

type String string

func (s String) Symbol() bool {
	for _, char := range s {
		if char == '.' {
			continue
		}
		if char >= '0' && char <= '9' {
			continue
		}
		return true
	}

	return false
}

type Grid struct {
	size int
	data []byte
}

func (g *Grid) Lands(p Point) bool {
	if p.X < 0 || p.Y < 0 {
		return true
	}
	if p.X+1 > g.size || p.Y+1 > g.size {
		return true
	}

	return false
}

func (g *Grid) Byte(p Point) byte {
	return g.data[g.Address(p)]
}

func (g *Grid) String(points Points) String {
	var buff []byte

	for _, point := range points {
		buff = append(buff, g.Byte(point))
	}

	return String(buff)
}

func (g *Grid) Address(p Point) int {
	return p.X + (p.Y * g.size)
}

func (g *Grid) Points() Points {
	var points Points

	for x := 0; x < g.size; x++ {
		for y := 0; y < g.size; y++ {
			points = append(points, Point{
				X: x,
				Y: y,
			})
		}
	}

	return points
}

type Points []Point

func (p Points) Overlap(other Points) bool {
	seen := map[[2]int]bool{}

	for _, pnt := range p {
		seen[[2]int{pnt.X, pnt.Y}] = true
	}

	for _, pnt := range other {
		k := [2]int{pnt.X, pnt.Y}

		if seen[k] {
			return true
		}
	}

	return false
}

func (p Points) Neighbors(grid *Grid) Points {
	self := map[int]bool{}

	for _, point := range p {
		address := grid.Address(point)
		self[address] = true
	}

	result := map[int]Point{}

	for _, p1 := range p {
		for _, p2 := range p1.Neighbors(grid) {
			address := grid.Address(p2)

			if !self[address] {
				result[address] = p2
			}
		}
	}

	var points Points

	for _, point := range result {
		points = append(points, point)
	}

	return points
}

type Point struct {
	X, Y int
}

func (c Point) Neighbors(grid *Grid) Points {
	var points Points

	for _, delta := range adjacentDelta {
		p := Point{
			X: c.X + delta[0],
			Y: c.Y + delta[1],
		}

		if !grid.Lands(p) {
			points = append(points, p)
		}
	}

	return points
}

func main() {
	path := os.Args[1]

	grid := loadGrid(path)

	puzzle1(grid)

	puzzle2(grid)
}

func puzzle1(grid *Grid) {
	total := 0

	for _, points := range numbers(grid) {
		str := grid.String(points)

		num, err := strconv.Atoi(string(str))

		if err != nil {
			panic(err)
		}

		adjacent := grid.String(points.Neighbors(grid))

		if adjacent.Symbol() {
			total += num
		}
	}

	fmt.Println(total)
}

func puzzle2(grid *Grid) {
	total := 0

	allNumbers := numbers(grid)

	for _, gear := range gears(grid) {
		var nums []Points

		for _, num := range allNumbers {
			if gear.Neighbors(grid).Overlap(num) {
				nums = append(nums, num)
			}
		}

		if len(nums) != 2 {
			continue
		}

		total += ratio(nums[0], nums[1], grid)
	}

	fmt.Println(total)
}

func ratio(num1 Points, num2 Points, grid *Grid) int {
	n1, err := strconv.Atoi(string(grid.String(num1)))

	if err != nil {
		panic(err)
	}

	n2, err := strconv.Atoi(string(grid.String(num2)))

	if err != nil {
		panic(err)
	}

	return n1 * n2
}

func numbers(grid *Grid) []Points {
	var nums []Points

	var points Points

	for pos, char := range grid.data {
		switch {
		case char >= '0' && char <= '9':
			y := int(math.Floor(float64(pos / grid.size)))
			x := pos - (y * grid.size)

			points = append(points, Point{
				X: x,
				Y: y,
			})
		default:
			if len(points) > 0 {
				nums = append(nums, points)
			}
			points = nil
		}
	}

	return nums
}

func gears(grid *Grid) []Point {
	var grs []Point

	for pos, char := range grid.data {
		if char != '*' {
			continue
		}

		y := int(math.Floor(float64(pos / grid.size)))
		x := pos - (y * grid.size)

		grs = append(grs, Point{
			X: x,
			Y: y,
		})
	}

	return grs
}

func loadGrid(path string) *Grid {
	f, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	var result Grid

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		b := scanner.Bytes()
		result.size += 1
		result.data = append(result.data, b...)
	}

	if len(result.data) != result.size*result.size {
		panic("bad grid")
	}

	return &result
}

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var source_file = "./source.txt"

func get_source_data() (data []string) {
	f, err := os.Open(source_file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		data = append(data, line)
	}

	return
}

func generate_name(analysis map[rune]map[rune]float64) (name string) {
	// Q: does it matter about not sorting the keys / are these probabilities working as expected?
	// I think so, the order doesn't matter... ??
	char := '_'
	for {
		pos := 0.0
		d := rand.Float64()
		for x := range analysis[char] {
			pos += analysis[char][x]
			if pos >= d {
				if x == '.' {
					return
				}
				name += string(x)
				char = x
				break
			}
		}
	}
}

// Unused
func visualise(analysis map[rune]map[rune]float64) {
	for m := range analysis {
		if len(analysis[m]) > 0 {
			fmt.Println(string(m))
			for x := range analysis[m] {
				fmt.Println("  ", string(x), " => ", analysis[m][x])
			}
		}
	}
}

func main() {
	data := get_source_data()
	analysis := make(map[rune]map[rune]float64)

	// set up some empty maps for convenience:
	alpha := "abcdefghijklmnopqrstuvwxyz_" // _ is a symbol for the start of a chain
	for i := range alpha {
		a := rune(alpha[i])
		analysis[a] = make(map[rune]float64)
	}

	for _, name := range data {
		// I only want to deal with lowercase...
		name = strings.ToLower(name)

		// deal with the start character, transition from "_" to this start character:
		start_char := rune(name[0])
		// uninitialised values are zero by default, so no need to be very elaborate...
		analysis['_'][start_char] = analysis['_'][start_char] + 1

		// now deal with the remaining characters.
		// for each character, we construct a map of what the next characters are.
		for i, c := range name {
			if i == len(name)-1 {
				// add an ending character "." for this character's map.
				analysis[c]['.'] = analysis[c]['.'] + 1
			} else {
				next_char := rune(name[i+1])
				analysis[c][next_char] = analysis[c][next_char] + 1
			}
		}
	}

	// At this stage I have counts for each character in the chain

	// Now I convert the counts into probabilities
	for m := range analysis {
		if len(analysis[m]) > 0 {
			total := 0.0
			for x := range analysis[m] {
				total += analysis[m][x]
			}
			for x := range analysis[m] {
				f := analysis[m][x] / total
				analysis[m][x] = f
			}
		} else {
			delete(analysis, m)
		}
	}

	// And now I generate the glorious names:
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		// For some reason it sometimes comes back with a 1-letter name, which is no good...
		n := " "
		for {
			n = generate_name(analysis)
			if len(n) > 1 {
				break
			}
		}
		fmt.Println(strings.Title(n))
	}

}

/*
* @(#) solve.go - sudoku solver
* @(#) $Id: solve.go,v 1.5 2018/08/26 21:33:34 bduncan Exp bduncan $
*
* Created: bduncan-sudoku@beachnet.org, Sat Aug 25 17:20:12 EDT 2018
*
* Description:
*   - uses simple recursive backtracking algorithm with up-front and
*     ongoing elimination of invalid tries by tracking for each row,
*     column and region..
*   - uses bitmap arrays for fast lookup
*   - conversion of the original C program
*
* Variables:
*   regmap[r][c]      ; pre-compiled, which region sector is r,c
*   master[r][c]      ; master matrix
*   C[col] & elem     ; true if elem in Column
*   R[row] & elem     ; true if elem in Row
*   Q[reg] & elem     ; true if elem in Region (Quadrant)
*
*
*/


package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const SUBORDER = 3
const ORDER = (SUBORDER * SUBORDER)
const MARK = 1
const UNMARK = 0

var count int = 0

var regmap [ORDER + 1][ORDER + 1]uint
var master [ORDER + 1][ORDER + 1]uint

var C [ORDER + 1]uint
var R [ORDER + 1]uint
var Q [ORDER + 1]uint


/* fregmap - returns the region, given row column */
func fregmap(r, c uint) uint {
	return regmap[r][c]
}

/* initregmap - initialize the region mapping */
func initregmap() {
	var i, j uint

	for i = 0; i < ORDER; i++ {
		for j = 0; j < ORDER; j++ {
			regmap[i+1][j+1] = i/SUBORDER*SUBORDER + j/SUBORDER + 1
		}
	}
}

/* inuse - returns whether value is already used in either row, column or region */
func inuse(r, c, try uint) uint {
	var q = fregmap(r, c)
	var bitmap uint = 1 << try

	return ((R[r] | C[c] | Q[q]) & bitmap)
}

/* mark - mark (or unmark) value in row column (and region) */
func mark(r, c, try, flag uint) {
	var bitmap uint = (1 << try)
	var q uint = fregmap(r, c)

	if flag == MARK {
		Q[q] |= bitmap
		R[r] |= bitmap
		C[c] |= bitmap
		master[r][c] = try
	} else {
		Q[q] &= ^bitmap
		R[r] &= ^bitmap
		C[c] &= ^bitmap
		master[r][c] = 0
	}
}

/* search - do a recursive, depth-first search */
func search(r, c uint) uint {
	var try uint

	count++

	for master[r][c] != 0 {
		c++
		if c > ORDER {
			c = 1
			r++
			if r > ORDER {
				return 1  /* return goodness! */
			}
		}
	}

	for try = 1; try <= ORDER; try++ {
		if inuse(r, c, try) == 0 {
			mark(r, c, try, MARK)
			if search(r, c) != 0 {  /* recurse! */
				return 1
			} /* else zero returned -- unwind, unmark */
			mark(r, c, try, UNMARK)
		}
	}

	return 0  /* return NOT found */
}


/* main - main function which reads from stdin, outputs result */
func main() {
	var i int
	var r uint
	var s []string

	scanner := bufio.NewScanner(os.Stdin)
	initregmap()

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		if strings.Contains(scanner.Text(), "#") {
			continue
		}
		s = strings.Fields(scanner.Text())
		if len(s) == 0 {
			continue
		}
		if len(s) == ORDER {
			r++
			for i = 0; i < len(s); i++ {
				val, err := strconv.Atoi(s[i])
				if err != nil {
					fmt.Printf("Error converting %q to int: %q\n", s[i], err)
				}
				mark(r, uint(i+1), uint(val), MARK)
			}
		} else {
			fmt.Printf("Error on line %d fields %d\n", r, len(s))
		}
	}
	for i = 1; i <= ORDER; i++ {
		fmt.Printf("master=%v\n", master[i][1:])
	}

	fmt.Printf("returned %d, count %d\n", search(1, 1), count)

	for i = 1; i <= ORDER; i++ {
		fmt.Printf("master=%v\n", master[i][1:])
	}

}


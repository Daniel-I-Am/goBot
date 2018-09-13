package main

import (
	"strings"
	"fmt"
	"math/rand"
	"unicode"
	"strconv"
	"regexp"

	"github.com/Knetic/govaluate"
)

func stripWhitespace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
		  return -1
		}
		return r
	  }, str)
}

func roll(content string) string {
	input := stripWhitespace(content[len(prefix)+4:])
	regex := regexp.MustCompile("(?P<pre>.*?)(?P<roll>\\d+d\\d+)(?P<suf>.*)")
	for regex.MatchString(input) == true {
		match := regex.FindStringSubmatch(input)
		input = match[1] + parseRoll(match[2]) + match[3]
	}
	// now that we have the dice rolled, we need to calculate the result
	expression, _ := govaluate.NewEvaluableExpression(input);
	result, _ := expression.Evaluate(nil);
	fmt.Println(input, "->", result)
	return input + " = " + strconv.FormatFloat(result.(float64), 'f', -1, 64)
}

func parseRoll(input string) string {
	regex := regexp.MustCompile("(?P<a>\\d+)d(?P<b>\\d+)")
	if !regex.MatchString(input) { return "" }
	match := regex.FindStringSubmatch(input)
	toRet := "("
	diceCount, _ := strconv.Atoi(match[1])
	diceSize, _ := strconv.Atoi(match[2])
	for i := 0; i < diceCount; i++ {
		rn := rand.Intn(diceSize) + 1 //+1.... 0-indexed
		toRet += strconv.Itoa(rn) + "+"
	}
	toRet = toRet[:len(toRet)-1]
	toRet += ")"
	return toRet
}
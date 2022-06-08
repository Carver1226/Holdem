package main

import (
	"fmt"
	"testing"
)

func Test_isStraight(t *testing.T) {
	test := "5s8d6c9cXn7d3h"
	cards := ConvertCards(test)
	fmt.Printf("%+v\n", cards)
	a, b := isStraight(cards)
	fmt.Printf("%v\t%v\n", a, b)
	test = "Xn6c8sTd9sJd8c"
	cards = ConvertCards(test)
	fmt.Printf("%+v\n", cards)
	a, b = isStraight(cards)
	fmt.Printf("%v\t%v\n", a, b)
}

func Test_Compare(t *testing.T) {
	a := "Jd9dTdAdQd3h8d"
	b := "KsKcJd9dTdAdQd"
	fmt.Println(Compare(a, b))
}

func Test_isStraightFlush(t *testing.T) {
	test := "Jd9dTdAdQd3h8d"
	cards := ConvertCards(test)
	fmt.Printf("%+v\n", cards)
	a, b := isStraightFlush(cards)
	fmt.Printf("%v\t%v\n", a, b)
}

func Test_isSame(t *testing.T) {
	test := "2dXn9h3d7hQhKh"
	cards := ConvertCards(test)
	fmt.Printf("%+v\n", cards)
	a, b := isSame(cards)
	fmt.Printf("%v\t%v\n", a, b)
}
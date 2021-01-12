package main

type Card string
type Cards []Card

var scrumCards Cards

func init() {
	scrumCards = []Card{
		"0",
		"1",
		"2",
		"3",
		"5",
		"8",
		"13",
		"21",
		"40",
		"80",
		"?",
		"coffee",
	}
}

func (cards Cards) contains(card Card) bool {
	for _, c := range cards {
		if c == card {
			return true
		}
	}
	return false
}

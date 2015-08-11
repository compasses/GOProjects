package main

import "fmt"
import "rand"

type score struct {
	player, opponent, thisTurn int
}

type action func(current score) (result score, turnIsOver bool)

func roll(s score) (score, bool) {
	outcome := rand.Intn(6) + 1
	if outcome == 1 {
		return score{s.opponent, s.player, 0}, true
	}
	return score{s.player, s.opponent, outcome + s.thisTurn}, false
}

func stay(s score) (score, bool) {
	return score{s.opponent, s.player + s.thisTurn, 0}, true
}

type strategy func(score) action

func stayAtK(k int) strategy {
	return func(s score) action {
		if s.thisTurn >= k {
			return stay
		}
		return roll
	}
}

func play(strategy1, strategy0 strategy) {

}

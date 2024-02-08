package main

import "testing"

func TestNyan(t *testing.T) {
    if !SolveNyan("nyannyacanyan") {
        t.Errorf("This is nyan")
    }

    if !SolveNyan("🐈dayo") {
        t.Errorf("This is nyan")
    }

    if !SolveNyan("にゃにゃんにゃんにゃか") {
        t.Errorf("This is nyan")
    }

    if SolveNyan("わんわんだ") {
        t.Errorf("This is not nyan")
    }
}

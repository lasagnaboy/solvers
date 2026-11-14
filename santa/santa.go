package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// History maps a year (e.g. "2025") to that year's giver => recipient assignments.
type History map[string]map[string]string

// loadHistory reads the JSON datastore. A missing file is treated as empty history.
func loadHistory(path string) (History, error) {
	h := History{}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return h, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return h, nil
}

// saveHistory writes the datastore atomically (temp file + rename).
func saveHistory(path string, h History) error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(data, '\n'), 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// forbiddenPairs collects giver => {recipients} from the `lookback` years before `year`.
// lookback <= 0 means "all recorded years before `year`".
func forbiddenPairs(h History, year, lookback int) map[string]map[string]bool {
	forbidden := map[string]map[string]bool{}
	for yStr, pairs := range h {
		y, err := strconv.Atoi(yStr)
		if err != nil || y >= year || (lookback > 0 && y < year-lookback) {
			continue
		}
		for giver, recipient := range pairs {
			if forbidden[giver] == nil {
				forbidden[giver] = map[string]bool{}
			}
			forbidden[giver][recipient] = true
		}
	}
	return forbidden
}

// arrange returns the santas ordered as a single gift-giving loop
// (order[i] gives to order[i+1], the last gives to the first), such that
// no giver is paired with a forbidden recipient. It uses randomized
// backtracking, so it always finds a valid loop if one exists.
func arrange(santas []string, forbidden map[string]map[string]bool, rng *rand.Rand) ([]string, error) {
	n := len(santas)
	if n < 2 {
		return nil, fmt.Errorf("need at least 2 santas, got %d", n)
	}
	ok := func(giver, recipient string) bool { return !forbidden[giver][recipient] }

	pool := append([]string(nil), santas...)
	rng.Shuffle(n, func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })

	order := []string{pool[0]}
	used := map[string]bool{pool[0]: true}

	var extend func() bool
	extend = func() bool {
		last := order[len(order)-1]
		if len(order) == n {
			return ok(last, order[0]) // close the loop
		}
		for _, next := range rng.Perm(n) {
			cand := pool[next]
			if used[cand] || !ok(last, cand) {
				continue
			}
			used[cand] = true
			order = append(order, cand)
			if extend() {
				return true
			}
			order = order[:len(order)-1]
			used[cand] = false
		}
		return false
	}

	if !extend() {
		return nil, errors.New("no arrangement avoids the previous recipients; try a smaller -lookback")
	}
	return order, nil
}

func main() {
	var strSantas []string
	s := flag.Bool("stockings", true, "a bool indicating whether to run for stocking stuffers")
	historyPath := flag.String("history", "santa_history.json", "path to the JSON file of past assignments")
	year := flag.Int("year", time.Now().Year(), "the year being arranged")
	lookback := flag.Int("lookback", 1, "how many previous years' recipients to avoid (0 = all recorded years)")
	dryRun := flag.Bool("dry-run", false, "print the assignments without saving them to the history file")
	force := flag.Bool("force", false, "overwrite assignments already saved for -year")
	names := flag.String("names", "", "comma-separated list of santas, e.g. \"Salim,Nico,John\" (overrides the built-in lists)")
	flag.Parse()
	fmt.Println("stockings:", *s)

	if *names != "" {
		for _, n := range strings.Split(*names, ",") {
			if n = strings.TrimSpace(n); n != "" {
				strSantas = append(strSantas, n)
			}
		}
	} else if *s {
		strSantas = []string{"Name1", "Name2"}
	} else {
		strSantas = []string{"Name1", "Name2"}
	}
	seen := map[string]bool{}
	for _, n := range strSantas {
		if seen[n] {
			fmt.Fprintf(os.Stderr, "%s is listed more than once\n", n)
			os.Exit(1)
		}
		seen[n] = true
	}

	history, err := loadHistory(*historyPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error loading history:", err)
		os.Exit(1)
	}
	yearKey := strconv.Itoa(*year)
	if !*force && !*dryRun {
		for _, name := range strSantas {
			if _, exists := history[yearKey][name]; exists {
				fmt.Fprintf(os.Stderr, "%s already has a %s assignment in %s; use -force to replace it or -dry-run to preview\n", name, yearKey, *historyPath)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("Arranging %d Santas ...\n\n", len(strSantas))
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	order, err := arrange(strSantas, forbiddenPairs(history, *year, *lookback), rng)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if *s {
		fmt.Print("Dear stocking stuffers:\n\n")
	} else {
		fmt.Print("Dear grown-up gift-exchangers:\n\n")
	}
	assignments := map[string]string{}
	for j := range order {
		giver, recipient := order[j], order[(j+1)%len(order)]
		assignments[giver] = recipient
		if *s {
			fmt.Printf("Ho! Ho! Ho! %s, you will be making %s especially merry this year by putting a special gift in their stocking. \n", giver, recipient)
		} else {
			fmt.Printf("'Tis the season! %s, you will be making %s especially merry this year with a Christmas gift!\n", giver, recipient)
		}
	}

	if *dryRun {
		return
	}
	// Merge into this year's entry so the stockings and grown-up runs can share a year.
	if history[yearKey] == nil {
		history[yearKey] = map[string]string{}
	}
	for giver, recipient := range assignments {
		history[yearKey][giver] = recipient
	}
	if err := saveHistory(*historyPath, history); err != nil {
		fmt.Fprintln(os.Stderr, "error saving history:", err)
		os.Exit(1)
	}
	fmt.Printf("\nSaved %s assignments to %s\n", yearKey, *historyPath)
} //main
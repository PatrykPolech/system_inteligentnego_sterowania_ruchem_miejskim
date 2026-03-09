package main

import "strings"

func aktualizujSterownikTestu() {
	dt := 1.0 / 60.0
	timerFazy += dt
	switch obecnyStan {
	case StanTestStaly:
		uruchomSterownikStaly()
	case StanTestZachlanny:
		uruchomSterownikZachlanny()
	case StanTestInteligentny:
		uruchomSterownikInteligentny()
	}
}

func uruchomSterownikStaly() {
	fazy := []string{"NS_STRAIGHT", "EW_STRAIGHT", "N_LEFT", "S_LEFT", "E_LEFT", "W_LEFT"}
	czasy := []float64{8.0, 8.0, 4.0, 4.0, 4.0, 4.0}
	dlugoscZielonego := czasy[stanSterownika]
	jestZolte := strings.Contains(obecnaFaza, "YELLOW")
	jestCzerwone := (obecnaFaza == "ALL_RED")
	if jestCzerwone {
		if timerFazy > 1.0 {
			obecnaFaza = fazy[stanSterownika]
			timerFazy = 0
		}
	} else if jestZolte {
		if timerFazy > 2.0 {
			obecnaFaza = "ALL_RED"
			timerFazy = 0
		}
	} else {
		if timerFazy > dlugoscZielonego {
			obecnaFaza += "_YELLOW"
			timerFazy = 0
			stanSterownika = (stanSterownika + 1) % len(fazy)
		}
	}
}

func uruchomSterownikZachlanny() {
	jestZolte := strings.Contains(obecnaFaza, "YELLOW")
	jestCzerwone := (obecnaFaza == "ALL_RED")
	if jestCzerwone {
		if timerFazy > 0.5 {
			najFaza := "NS_STRAIGHT"
			maxQ := -1
			kolejnosc := []string{"NS_STRAIGHT", "EW_STRAIGHT", "N_LEFT", "S_LEFT", "E_LEFT", "W_LEFT"}
			for _, p := range kolejnosc {
				q := 0
				switch p {
				case "NS_STRAIGHT":
					q = pobierzRozmiarKolejki("N", "STRAIGHT") + pobierzRozmiarKolejki("N", "RIGHT") + pobierzRozmiarKolejki("S", "STRAIGHT") + pobierzRozmiarKolejki("S", "RIGHT")
				case "EW_STRAIGHT":
					q = pobierzRozmiarKolejki("E", "STRAIGHT") + pobierzRozmiarKolejki("E", "RIGHT") + pobierzRozmiarKolejki("W", "STRAIGHT") + pobierzRozmiarKolejki("W", "RIGHT")
				case "N_LEFT":
					q = pobierzRozmiarKolejki("N", "LEFT")
				case "S_LEFT":
					q = pobierzRozmiarKolejki("S", "LEFT")
				case "E_LEFT":
					q = pobierzRozmiarKolejki("E", "LEFT")
				case "W_LEFT":
					q = pobierzRozmiarKolejki("W", "LEFT")
				}
				if q > maxQ {
					maxQ = q
					najFaza = p
				}
			}
			if maxQ == 0 {
				timerFazy = 0
				return
			}
			obecnaFaza = najFaza
			timerFazy = 0
		}
	} else if jestZolte {
		if timerFazy > 2.0 {
			obecnaFaza = "ALL_RED"
			timerFazy = 0
		}
	} else {
		if timerFazy > 5.0 {
			obecnaFaza += "_YELLOW"
			timerFazy = 0
		}
	}
}

func uruchomSterownikInteligentny() {
	fazy := []string{"NS_STRAIGHT", "N_LEFT", "S_LEFT", "EW_STRAIGHT", "E_LEFT", "W_LEFT"}
	kolejkiFazy := map[string][][]string{
		"NS_STRAIGHT": {{"N", "STRAIGHT"}, {"N", "RIGHT"}, {"S", "STRAIGHT"}, {"S", "RIGHT"}},
		"EW_STRAIGHT": {{"E", "STRAIGHT"}, {"E", "RIGHT"}, {"W", "STRAIGHT"}, {"W", "RIGHT"}},
		"N_LEFT":      {{"N", "LEFT"}}, "S_LEFT": {{"S", "LEFT"}}, "E_LEFT": {{"E", "LEFT"}}, "W_LEFT": {{"W", "LEFT"}},
	}
	jestZolte := strings.Contains(obecnaFaza, "YELLOW")
	jestCzerwone := (obecnaFaza == "ALL_RED")
	nazwaFazy := fazy[stanSterownika]
	if jestCzerwone {
		if timerFazy > 1.0 {
			rozm := 0
			for _, def := range kolejkiFazy[nazwaFazy] {
				rozm += pobierzRozmiarKolejki(def[0], def[1])
			}
			if rozm == 0 {
				stanSterownika = (stanSterownika + 1) % len(fazy)
				timerFazy = 1.0
				return
			}
			obecnaFaza = nazwaFazy
			timerFazy = 0
		}
	} else if jestZolte {
		if timerFazy > 2.0 {
			obecnaFaza = "ALL_RED"
			timerFazy = 0
			stanSterownika = (stanSterownika + 1) % len(fazy)
		}
	} else {
		if timerFazy < 2.0 {
			return
		}
		if timerFazy > 8.0 {
			obecnaFaza += "_YELLOW"
			timerFazy = 0
			return
		}
		rozm := 0
		for _, def := range kolejkiFazy[nazwaFazy] {
			rozm += pobierzRozmiarKolejki(def[0], def[1])
		}
		if rozm == 0 {
			obecnaFaza += "_YELLOW"
			timerFazy = 0
		}
	}
}

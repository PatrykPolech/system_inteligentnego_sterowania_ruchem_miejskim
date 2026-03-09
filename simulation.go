package main

import (
	"image/color"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func startJazdaWolna() {
	obecnyStan = StanJazdaWolna
	ebiten.SetTPS(int(60 * konfigPredkosc))
	mutexSamochodow.Lock()
	samochody = make([]*Samochod, 0)
	mutexSamochodow.Unlock()
	wszystkieAuta = 0
	czasTrwaniaSym = 0
	timerFazy = 0
	stanSterownika = 0
	obecnaFaza = "ALL_RED"
	timerSpawnuWolnego = 0
	nastepnySpawnWolny = 0.5
}

func aktualizujSpawnerWolny(dt float64) {
	timerSpawnuWolnego += dt
	if timerSpawnuWolnego >= nastepnySpawnWolny {
		stworzLosoweAuto()
		timerSpawnuWolnego = 0
		nastepnySpawnWolny = 0.5 + rand.Float64()
	}
}

func startTest() {
	ebiten.SetTPS(int(60 * konfigPredkosc))
	generujScenariusz()
	wyniki = make([]WynikSymulacji, 0)
	indeksStronyWykresu = 0
	startPrzebiegTestu(StanTestStaly)
}

func generujScenariusz() {
	scenariusz = make([]ZdarzenieSpawnu, 0)
	r := rand.New(rand.NewSource(42))
	czasMs := 0
	for i := 0; i < konfigLiczbaAut; i++ {
		gap := r.Intn(2500) + 500
		czasMs += gap
		kierunki := []string{"N", "S", "E", "W"}
		kier := kierunki[r.Intn(4)]
		akcja := "STRAIGHT"
		pas := "OUTER"
		if r.Float64() < 0.3 {
			akcja = "LEFT"
			pas = "INNER"
		} else if r.Float64() > 0.5 {
			akcja = "RIGHT"
		}
		jestCiezarowka := r.Float64() < 0.2
		jestUprzyw := r.Float64() < 0.05
		scenariusz = append(scenariusz, ZdarzenieSpawnu{OpoznienieMs: czasMs, Kierunek: kier, Akcja: akcja, Pas: pas, JestUprzyw: jestUprzyw, JestCiezarowka: jestCiezarowka})
	}
}

func startPrzebiegTestu(stan StanGry) {
	obecnyStan = stan
	mutexSamochodow.Lock()
	samochody = make([]*Samochod, 0)
	mutexSamochodow.Unlock()
	indeksScenariusza = 0
	czasStartuTestu = time.Now()
	obecnaFaza = "ALL_RED"
	stanSterownika = 0
	timerFazy = 0
	testLacznyCzasOczekiwania = 0
	testLacznyCzasPodrozy = 0
	testMaxCzasOczekiwania = 0
	testLicznikAut = 0
}

func aktualizujSpawnerTestu() {
	minelo := time.Since(czasStartuTestu).Milliseconds()
	minelo = int64(float64(minelo) * konfigPredkosc)
	for indeksScenariusza < len(scenariusz) {
		zd := scenariusz[indeksScenariusza]
		if int64(zd.OpoznienieMs) <= minelo {
			stworzAutoTestowe(zd)
			indeksScenariusza++
		} else {
			break
		}
	}
}

func stworzAutoTestowe(zd ZdarzenieSpawnu) {
	mutexSamochodow.Lock()
	defer mutexSamochodow.Unlock()
	startX, startY := 0.0, 0.0
	switch zd.Kierunek {
	case "N":
		startY = -50
		if zd.Pas == "INNER" {
			startX = X_Kol2_Wew
		} else {
			startX = X_Kol1_Zew
		}
	case "S":
		startY = WysokoscEkranu + 50
		if zd.Pas == "INNER" {
			startX = X_Kol3_Wew
		} else {
			startX = X_Kol4_Zew
		}
	case "E":
		startX = SzerokoscEkranu + 50
		if zd.Pas == "INNER" {
			startY = Y_Wiersz2_Wew
		} else {
			startY = Y_Wiersz1_Zew
		}
	case "W":
		startX = -50
		if zd.Pas == "INNER" {
			startY = Y_Wiersz3_Wew
		} else {
			startY = Y_Wiersz4_Zew
		}
	}
	dl := 22.0
	acc := 0.15
	maxV := 5.0
	paleta := []color.RGBA{{255, 0, 0, 255}, {0, 255, 0, 255}, {0, 0, 255, 255}, {255, 255, 0, 255}, {255, 0, 255, 255}, {0, 255, 255, 255}, {255, 128, 0, 255}, {128, 0, 128, 255}}
	kol := paleta[(indeksScenariusza)%len(paleta)]
	if zd.JestCiezarowka {
		dl = 46.0
		acc = 0.10
		kol = color.RGBA{139, 69, 19, 255}
	}
	if zd.JestUprzyw {
		dl = 22.0
		acc = 0.25
		maxV = 5.5
		kol = color.RGBA{0, 0, 200, 255}
	}
	noweAuto := &Samochod{X: startX, Y: startY, Kierunek: zd.Kierunek, Akcja: zd.Akcja, TypPasa: zd.Pas, Predkosc: maxV, Dlugosc: dl, Szerokosc: 22.0, Kolor: kol, MaxPredkosc: maxV, Przyspieszenie: acc, JestUprzyw: zd.JestUprzyw, CzasPojawienia: time.Now()}
	samochody = append(samochody, noweAuto)
}

func stworzLosoweAuto() {
	mutexSamochodow.Lock()
	defer mutexSamochodow.Unlock()
	kierunki := []string{"N", "S", "E", "W"}
	kier := kierunki[rand.Intn(4)]
	akcja := "STRAIGHT"
	pas := "OUTER"
	if rand.Float64() < 0.3 {
		akcja = "LEFT"
		pas = "INNER"
	} else if rand.Float64() > 0.5 {
		akcja = "RIGHT"
	}
	startX, startY := 0.0, 0.0
	switch kier {
	case "N":
		startY = -50
		if pas == "INNER" {
			startX = X_Kol2_Wew
		} else {
			startX = X_Kol1_Zew
		}
	case "S":
		startY = WysokoscEkranu + 50
		if pas == "INNER" {
			startX = X_Kol3_Wew
		} else {
			startX = X_Kol4_Zew
		}
	case "E":
		startX = SzerokoscEkranu + 50
		if pas == "INNER" {
			startY = Y_Wiersz2_Wew
		} else {
			startY = Y_Wiersz1_Zew
		}
	case "W":
		startX = -50
		if pas == "INNER" {
			startY = Y_Wiersz3_Wew
		} else {
			startY = Y_Wiersz4_Zew
		}
	}
	zablokowane := false
	for _, c := range samochody {
		if math.Abs(c.X-startX) < 70 && math.Abs(c.Y-startY) < 70 {
			zablokowane = true
			break
		}
	}
	if !zablokowane {
		dl := 22.0
		acc := 0.15
		maxV := 5.0
		paleta := []color.RGBA{{255, 0, 0, 255}, {0, 255, 0, 255}, {0, 0, 255, 255}, {255, 255, 0, 255}, {255, 0, 255, 255}, {0, 255, 255, 255}, {255, 128, 0, 255}, {128, 0, 128, 255}}
		kol := paleta[rand.Intn(len(paleta))]
		jestUprzyw := false
		if rand.Float64() < 0.2 {
			dl = 46.0
			acc = 0.10
			kol = color.RGBA{139, 69, 19, 255}
		}
		if rand.Float64() < 0.05 {
			jestUprzyw = true
			kol = color.RGBA{0, 0, 200, 255}
			maxV = 5.5
			dl = 22.0
		}
		noweAuto := &Samochod{X: startX, Y: startY, Kierunek: kier, Akcja: akcja, TypPasa: pas, Predkosc: maxV, Dlugosc: dl, Szerokosc: 22.0, Kolor: kol, MaxPredkosc: maxV, Przyspieszenie: acc, JestUprzyw: jestUprzyw, CzasPojawienia: time.Now()}
		samochody = append(samochody, noweAuto)
	}
}

func sprawdzKoniecTestu() {
	mutexSamochodow.Lock()
	l := len(samochody)
	mutexSamochodow.Unlock()
	if indeksScenariusza >= len(scenariusz) && l == 0 {
		zapiszWynik()
		if obecnyStan == StanTestStaly {
			startPrzebiegTestu(StanTestZachlanny)
		} else if obecnyStan == StanTestZachlanny {
			startPrzebiegTestu(StanTestInteligentny)
		} else if obecnyStan == StanTestInteligentny {
			ebiten.SetTPS(60)
			obecnyStan = StanWyniki
		}
	}
}

func zapiszWynik() {
	nazwa := ""
	switch obecnyStan {
	case StanTestStaly:
		nazwa = "1. Czasowe"
	case StanTestZachlanny:
		nazwa = "2. Czujnik ilosci pojazdow"
	case StanTestInteligentny:
		nazwa = "3. Inteligentne"
	}
	czasReal := time.Since(czasStartuTestu)
	czasSim := time.Duration(float64(czasReal) * konfigPredkosc)
	avgWait := time.Duration(0)
	avgTravel := time.Duration(0)
	maxWait := testMaxCzasOczekiwania
	if testLicznikAut > 0 {
		avgWait = testLacznyCzasOczekiwania / time.Duration(testLicznikAut)
		avgTravel = time.Duration(float64(testLacznyCzasPodrozy)*konfigPredkosc) / time.Duration(testLicznikAut)
	}
	wyniki = append(wyniki, WynikSymulacji{Nazwa: nazwa, CzasCalkowity: czasSim, SredniCzasOczekiwania: avgWait, MaxCzasOczekiwania: maxWait, LacznieAut: testLicznikAut, SredniaPodroz: avgTravel})
}

func aktualizujFizyke() {
	mutexSamochodow.Lock()
	defer mutexSamochodow.Unlock()
	karetkaAktywna = false
	karetkaZrodlo = ""
	for _, c := range samochody {
		if c.JestUprzyw {
			karetkaAktywna = true
			karetkaZrodlo = c.Kierunek
			break
		}
	}
	aktywne := []*Samochod{}
	for _, auto := range samochody {
		if auto.Predkosc < 0.5 {
			auto.CzasOczekiwania += 16 * time.Millisecond
		}
		celPredkosc := auto.MaxPredkosc
		dystansPrzeszkody := 9999.0
		for _, inne := range samochody {
			if auto == inne {
				continue
			}
			tenSamPas := false
			if !auto.Skrecil && !inne.Skrecil {
				if auto.Kierunek == inne.Kierunek && auto.TypPasa == inne.TypPasa {
					tenSamPas = true
				}
			} else if auto.Skrecil && inne.Skrecil {
				if math.Abs(auto.X-inne.X) < 30 && math.Abs(auto.Y-inne.Y) < 30 {
					tenSamPas = true
				}
			}
			if tenSamPas {
				dist := 0.0
				jestPrzed := false
				switch auto.Kierunek {
				case "N":
					if !auto.Skrecil {
						if inne.Y > auto.Y {
							dist = inne.Y - (auto.Y + auto.Dlugosc)
							jestPrzed = true
						}
					} else {
						dx := inne.X - auto.X
						if auto.Akcja == "LEFT" && dx > 0 {
							dist = dx - auto.Dlugosc
							jestPrzed = true
						} else if auto.Akcja == "RIGHT" && dx < 0 {
							dist = -(dx + inne.Dlugosc)
							jestPrzed = true
						}
					}
				case "S":
					if !auto.Skrecil {
						if inne.Y < auto.Y {
							dist = auto.Y - (inne.Y + inne.Dlugosc)
							jestPrzed = true
						}
					} else {
						dx := inne.X - auto.X
						if auto.Akcja == "LEFT" && dx < 0 {
							dist = -(dx + inne.Dlugosc)
							jestPrzed = true
						} else if auto.Akcja == "RIGHT" && dx > 0 {
							dist = dx - auto.Dlugosc
							jestPrzed = true
						}
					}
				case "E":
					if !auto.Skrecil {
						if inne.X < auto.X {
							dist = auto.X - (inne.X + inne.Dlugosc)
							jestPrzed = true
						}
					} else {
						dy := inne.Y - auto.Y
						if auto.Akcja == "LEFT" && dy > 0 {
							dist = dy - auto.Dlugosc
							jestPrzed = true
						} else if auto.Akcja == "RIGHT" && dy < 0 {
							dist = -(dy + inne.Dlugosc)
							jestPrzed = true
						}
					}
				case "W":
					if !auto.Skrecil {
						if inne.X > auto.X {
							dist = inne.X - (auto.X + auto.Dlugosc)
							jestPrzed = true
						}
					} else {
						dy := inne.Y - auto.Y
						if auto.Akcja == "LEFT" && dy < 0 {
							dist = -(dy + inne.Dlugosc)
							jestPrzed = true
						} else if auto.Akcja == "RIGHT" && dy > 0 {
							dist = dy - auto.Dlugosc
							jestPrzed = true
						}
					}
				}
				if auto.Skrecil && inne.Skrecil {
					dist = 50
					jestPrzed = true
				}
				if jestPrzed && dist < 1000 && dist >= -5 {
					if dist < dystansPrzeszkody {
						dystansPrzeszkody = dist
					}
				}
			}
		}
		if dystansPrzeszkody < BezpiecznyDystans {
			celPredkosc = 0.0
		} else if dystansPrzeszkody < BezpiecznyDystans*4.0 {
			celPredkosc = auto.MaxPredkosc * 0.3
		}
		if !auto.MinalStop {
			moznaJechac := sprawdzSwiatlo(auto.Kierunek, auto.Akcja)
			if auto.JestUprzyw {
				moznaJechac = true
			}
			dystDoLinii := 9999.0
			switch auto.Kierunek {
			case "N":
				dystDoLinii = LiniaStop_N - (auto.Y + auto.Dlugosc)
			case "S":
				dystDoLinii = auto.Y - LiniaStop_S
			case "E":
				dystDoLinii = auto.X - LiniaStop_E
			case "W":
				dystDoLinii = LiniaStop_W - (auto.X + auto.Dlugosc)
			}
			if !moznaJechac {
				if dystDoLinii < 10 {
					celPredkosc = 0
				} else if dystDoLinii < 180 {
					celPredkosc = math.Min(celPredkosc, auto.MaxPredkosc*(dystDoLinii/180.0))
					if dystDoLinii < 12 {
						celPredkosc = 0
					}
				}
			} else {
				if dystDoLinii <= -5 {
					auto.MinalStop = true
				}
			}
		}
		if !auto.Skrecil && auto.Akcja != "STRAIGHT" && auto.MinalStop {
			celPredkosc = math.Min(celPredkosc, 2.0)
		}
		if auto.Predkosc < celPredkosc {
			auto.Predkosc += auto.Przyspieszenie
		} else if auto.Predkosc > celPredkosc {
			auto.Predkosc -= Hamowanie
		}
		ruszAuto(auto)
		wSrodku := true
		if auto.X < -150 || auto.X > SzerokoscEkranu+150 || auto.Y < -150 || auto.Y > WysokoscEkranu+150 {
			wSrodku = false
		}
		if wSrodku {
			aktywne = append(aktywne, auto)
		} else {
			if obecnyStan != StanJazdaWolna {
				testLicznikAut++
				testLacznyCzasOczekiwania += auto.CzasOczekiwania
				testLacznyCzasPodrozy += time.Since(auto.CzasPojawienia)
				if auto.CzasOczekiwania > testMaxCzasOczekiwania {
					testMaxCzasOczekiwania = auto.CzasOczekiwania
				}
			}
			wszystkieAuta++
		}
	}
	samochody = aktywne
}

func ruszAuto(c *Samochod) {
	s := c.Predkosc
	switch c.Kierunek {
	case "N":
		if !c.Skrecil {
			c.Y += s
		} else {
			if c.Akcja == "LEFT" {
				c.X += s
			} else if c.Akcja == "RIGHT" {
				c.X -= s
			} else {
				c.Y += s
			}
		}
		if !c.Skrecil && ((c.Akcja == "RIGHT" && c.Y >= Y_Wiersz1_Zew) || (c.Akcja == "LEFT" && c.Y >= Y_Wiersz3_Wew)) {
			c.Skrecil = true
			if c.Akcja == "RIGHT" {
				c.Y = Y_Wiersz1_Zew
			} else {
				c.Y = Y_Wiersz3_Wew
			}
		}
	case "S":
		if !c.Skrecil {
			c.Y -= s
		} else {
			if c.Akcja == "LEFT" {
				c.X -= s
			} else if c.Akcja == "RIGHT" {
				c.X += s
			} else {
				c.Y -= s
			}
		}
		if !c.Skrecil && ((c.Akcja == "RIGHT" && c.Y <= Y_Wiersz4_Zew) || (c.Akcja == "LEFT" && c.Y <= Y_Wiersz2_Wew)) {
			c.Skrecil = true
			if c.Akcja == "RIGHT" {
				c.Y = Y_Wiersz4_Zew
			} else {
				c.Y = Y_Wiersz2_Wew
			}
		}
	case "E":
		if !c.Skrecil {
			c.X -= s
		} else {
			if c.Akcja == "LEFT" {
				c.Y += s
			} else if c.Akcja == "RIGHT" {
				c.Y -= s
			} else {
				c.X -= s
			}
		}
		if !c.Skrecil && ((c.Akcja == "RIGHT" && c.X <= X_Kol4_Zew) || (c.Akcja == "LEFT" && c.X <= X_Kol2_Wew)) {
			c.Skrecil = true
			if c.Akcja == "RIGHT" {
				c.X = X_Kol4_Zew
			} else {
				c.X = X_Kol2_Wew
			}
		}
	case "W":
		if !c.Skrecil {
			c.X += s
		} else {
			if c.Akcja == "LEFT" {
				c.Y -= s
			} else if c.Akcja == "RIGHT" {
				c.Y += s
			} else {
				c.X += s
			}
		}
		if !c.Skrecil && ((c.Akcja == "RIGHT" && c.X >= X_Kol1_Zew) || (c.Akcja == "LEFT" && c.X >= X_Kol3_Wew)) {
			c.Skrecil = true
			if c.Akcja == "RIGHT" {
				c.X = X_Kol1_Zew
			} else {
				c.X = X_Kol3_Wew
			}
		}
	}
}

func sprawdzSwiatlo(kier, akcja string) bool {
	if karetkaAktywna {
		if kier == karetkaZrodlo {
			return true
		} else {
			return false
		}
	}
	zieloneProsto := false
	zieloneLewo := false
	if strings.HasPrefix(obecnaFaza, "NS_STRAIGHT") && !strings.Contains(obecnaFaza, "YELLOW") {
		if kier == "N" || kier == "S" {
			zieloneProsto = true
		}
	}
	if strings.HasPrefix(obecnaFaza, "EW_STRAIGHT") && !strings.Contains(obecnaFaza, "YELLOW") {
		if kier == "E" || kier == "W" {
			zieloneProsto = true
		}
	}
	if obecnaFaza == "N_LEFT" && kier == "N" {
		zieloneLewo = true
	}
	if obecnaFaza == "S_LEFT" && kier == "S" {
		zieloneLewo = true
	}
	if obecnaFaza == "E_LEFT" && kier == "E" {
		zieloneLewo = true
	}
	if obecnaFaza == "W_LEFT" && kier == "W" {
		zieloneLewo = true
	}
	if akcja == "LEFT" {
		return zieloneLewo
	}
	return zieloneProsto
}

func pobierzRozmiarKolejki(kier, akcja string) int {
	liczba := 0
	mutexSamochodow.Lock()
	defer mutexSamochodow.Unlock()
	for _, c := range samochody {
		if c.Kierunek == kier && !c.MinalStop && c.Predkosc < 4.0 {
			if akcja == "" || c.Akcja == akcja {
				liczba++
			}
		}
	}
	return liczba
}

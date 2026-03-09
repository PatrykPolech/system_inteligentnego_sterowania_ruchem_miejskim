package main

import (
	"image/color"
	"sync"
	"time"
)

const (
	LiczbaAutTest     = 50
	SzerokoscEkranu   = 1024
	WysokoscEkranu    = 768
	SrodekX           = 512.0
	SrodekY           = 384.0
	SzerokoscAuta     = 22.0
	RozmiarSwiatla    = 20.0
	Hamowanie         = 0.35
	BezpiecznyDystans = 12.0
	Krawedz_N         = SrodekY - 80.0
	Krawedz_S         = SrodekY + 80.0
	Krawedz_W         = SrodekX - 80.0
	Krawedz_E         = SrodekX + 80.0
	OdstepStop        = 40.0
	LiniaStop_N       = Krawedz_N - OdstepStop
	LiniaStop_S       = Krawedz_S + OdstepStop
	LiniaStop_W       = Krawedz_W - OdstepStop
	LiniaStop_E       = Krawedz_E + OdstepStop
	X_Kol1_Zew        = SrodekX - 71.0
	X_Kol2_Wew        = SrodekX - 31.0
	X_Kol3_Wew        = SrodekX + 9.0
	X_Kol4_Zew        = SrodekX + 49.0
	Y_Wiersz1_Zew     = SrodekY - 71.0
	Y_Wiersz2_Wew     = SrodekY - 31.0
	Y_Wiersz3_Wew     = SrodekY + 9.0
	Y_Wiersz4_Zew     = SrodekY + 49.0
	SwiatloY_N        = LiniaStop_N + 5.0
	SwiatloY_S        = LiniaStop_S - 25.0
	SwiatloX_W        = LiniaStop_W + 5.0
	SwiatloX_E        = LiniaStop_E - 25.0
)

type Samochod struct {
	ID              int
	X, Y            float64
	Kierunek        string
	Akcja           string
	TypPasa         string
	Predkosc        float64
	Skrecil         bool
	MinalStop       bool
	CzasPojawienia  time.Time
	CzasOczekiwania time.Duration
	Dlugosc         float64
	Szerokosc       float64
	Kolor           color.Color
	MaxPredkosc     float64
	Przyspieszenie  float64
	JestUprzyw      bool
}

type ZdarzenieSpawnu struct {
	OpoznienieMs   int
	Kierunek       string
	Akcja          string
	Pas            string
	JestUprzyw     bool
	JestCiezarowka bool
}

type WynikSymulacji struct {
	Nazwa                 string
	CzasCalkowity         time.Duration
	SredniaPodroz         time.Duration
	SredniCzasOczekiwania time.Duration
	MaxCzasOczekiwania    time.Duration
	LacznieAut            int
}

type StanGry int

const (
	StanMenuGlowne StanGry = iota
	StanKonfigAutaTest
	StanKonfigPredkoscTest
	StanKonfigPredkoscWolna
	StanJazdaWolna
	StanTestStaly
	StanTestZachlanny
	StanTestInteligentny
	StanWyniki
)

var (
	konfigLiczbaAut                int     = 50
	konfigPredkosc                 float64 = 1.0
	obecnyStan                     StanGry = StanMenuGlowne
	samochody                      []*Samochod
	mutexSamochodow                sync.Mutex
	obecnaFaza                     string    = "ALL_RED"
	wszystkieAuta                  int       = 0
	czasStartu                     time.Time = time.Now()
	karetkaAktywna                 bool      = false
	karetkaZrodlo                  string    = ""
	czasTrwaniaSym                 float64   = 0
	timerSpawnuWolnego             float64   = 0
	nastepnySpawnWolny             float64   = 0
	scenariusz                     []ZdarzenieSpawnu
	indeksScenariusza              int
	czasStartuTestu                time.Time
	wyniki                         []WynikSymulacji
	indeksStronyWykresu            int = 0
	ostatniaAktualizacjaSterownika time.Time
	stanSterownika                 int
	timerFazy                      float64
	testLacznyCzasOczekiwania      time.Duration
	testLacznyCzasPodrozy          time.Duration
	testMaxCzasOczekiwania         time.Duration
	testLicznikAut                 int
)

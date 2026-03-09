package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Gra struct{}

func (g *Gra) Update() error {
	dt := 1.0 / 60.0
	switch obecnyStan {
	case StanMenuGlowne:
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			obecnyStan = StanKonfigPredkoscWolna
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			obecnyStan = StanKonfigAutaTest
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			return ebiten.Termination
		}
	case StanKonfigPredkoscWolna:
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			konfigPredkosc = 0.5
			startJazdaWolna()
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			konfigPredkosc = 1.0
			startJazdaWolna()
		}
		if inpututil.IsKeyJustPressed(ebiten.Key3) {
			konfigPredkosc = 2.0
			startJazdaWolna()
		}
		if inpututil.IsKeyJustPressed(ebiten.Key4) {
			konfigPredkosc = 5.0
			startJazdaWolna()
		}
	case StanKonfigAutaTest:
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			konfigLiczbaAut = 10
			obecnyStan = StanKonfigPredkoscTest
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			konfigLiczbaAut = 20
			obecnyStan = StanKonfigPredkoscTest
		}
		if inpututil.IsKeyJustPressed(ebiten.Key3) {
			konfigLiczbaAut = 50
			obecnyStan = StanKonfigPredkoscTest
		}
		if inpututil.IsKeyJustPressed(ebiten.Key4) {
			konfigLiczbaAut = 100
			obecnyStan = StanKonfigPredkoscTest
		}
	case StanKonfigPredkoscTest:
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			konfigPredkosc = 0.5
			startTest()
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			konfigPredkosc = 1.0
			startTest()
		}
		if inpututil.IsKeyJustPressed(ebiten.Key3) {
			konfigPredkosc = 2.0
			startTest()
		}
		if inpututil.IsKeyJustPressed(ebiten.Key4) {
			konfigPredkosc = 5.0
			startTest()
		}
	case StanJazdaWolna:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			wrocDoMenu()
		} else {
			czasTrwaniaSym += dt
			aktualizujSpawnerWolny(dt)
			timerFazy += dt
			uruchomSterownikInteligentny()
			aktualizujFizyke()
		}
	case StanTestStaly, StanTestZachlanny, StanTestInteligentny:
		aktualizujSpawnerTestu()
		aktualizujSterownikTestu()
		aktualizujFizyke()
		sprawdzKoniecTestu()
	case StanWyniki:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			indeksStronyWykresu = (indeksStronyWykresu + 1) % 4
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			wrocDoMenu()
		}
	}
	return nil
}

func wrocDoMenu() {
	obecnyStan = StanMenuGlowne
	ebiten.SetTPS(60)
	mutexSamochodow.Lock()
	samochody = make([]*Samochod, 0)
	mutexSamochodow.Unlock()
}

func (g *Gra) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{34, 139, 34, 255})
	switch obecnyStan {
	case StanMenuGlowne:
		rysujWycentrowanyTekst(screen, "Inteligentne sterowanie ruchem ulicznym", 200, color.RGBA{100, 255, 100, 255})
		rysujWycentrowanyTekst(screen, "Patryk Polechonski s197861", 230, color.RGBA{200, 200, 200, 255})
		rysujWycentrowanyTekst(screen, "MENU GLOWNE", 350, color.White)
		rysujWycentrowanyTekst(screen, "[1] Swobodna Jazda", 400, color.RGBA{100, 255, 255, 255})
		rysujWycentrowanyTekst(screen, "[2] Porownanie Algorytmow", 430, color.RGBA{255, 200, 100, 255})
		rysujWycentrowanyTekst(screen, "[ESC] Wyjscie", 500, color.RGBA{150, 150, 150, 255})
	case StanKonfigPredkoscWolna:
		rysujWycentrowanyTekst(screen, "PREDKOSC SYMULACJI:", 300, color.White)
		rysujWycentrowanyTekst(screen, "[1] 50%", 350, color.White)
		rysujWycentrowanyTekst(screen, "[2] 100%", 380, color.White)
		rysujWycentrowanyTekst(screen, "[3] 200%", 410, color.RGBA{255, 200, 100, 255})
		rysujWycentrowanyTekst(screen, "[4] 500%", 440, color.RGBA{255, 100, 100, 255})
	case StanKonfigAutaTest:
		rysujWycentrowanyTekst(screen, "LICZBA POJAZDOW:", 300, color.White)
		rysujWycentrowanyTekst(screen, "[1] 10 Aut", 350, color.White)
		rysujWycentrowanyTekst(screen, "[2] 20 Aut", 380, color.White)
		rysujWycentrowanyTekst(screen, "[3] 50 Aut", 410, color.White)
		rysujWycentrowanyTekst(screen, "[4] 100 Aut", 440, color.RGBA{255, 100, 100, 255})
	case StanKonfigPredkoscTest:
		rysujWycentrowanyTekst(screen, "PREDKOSC SYMULACJI:", 300, color.White)
		rysujWycentrowanyTekst(screen, "[1] 50%", 350, color.White)
		rysujWycentrowanyTekst(screen, "[2] 100%", 380, color.White)
		rysujWycentrowanyTekst(screen, "[3] 200%", 410, color.RGBA{255, 200, 100, 255})
		rysujWycentrowanyTekst(screen, "[4] 500%", 440, color.RGBA{255, 100, 100, 255})
	case StanJazdaWolna:
		rysujSymulacje(screen)
		rysujPanelStatystyk(screen, "TRYB SWOBODNY")
		rysujWycentrowanyTekst(screen, "[ENTER] Powrot do Menu", WysokoscEkranu-20, color.White)
	case StanTestStaly:
		rysujSymulacje(screen)
		rysujPanelStatystyk(screen, fmt.Sprintf("Czasowe (%d Aut)", konfigLiczbaAut))
	case StanTestZachlanny:
		rysujSymulacje(screen)
		rysujPanelStatystyk(screen, fmt.Sprintf("Czujnik ilosci pojazdow (%d Aut)", konfigLiczbaAut))
	case StanTestInteligentny:
		rysujSymulacje(screen)
		rysujPanelStatystyk(screen, fmt.Sprintf("Inteligentne (%d Aut)", konfigLiczbaAut))
	case StanWyniki:
		rysujWykresWynikow(screen)
	}
}

func (g *Gra) Layout(w, h int) (int, int) { return SzerokoscEkranu, WysokoscEkranu }

func main() {
	ebiten.SetWindowSize(SzerokoscEkranu, WysokoscEkranu)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Symulacja inteligentnego systemu sterowania ruchem miejskim")
	if err := ebiten.RunGame(&Gra{}); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"image/color"
	"sort"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

func rysujWycentrowanyTekst(screen *ebiten.Image, msg string, y int, clr color.Color) {
	rect := text.BoundString(basicfont.Face7x13, msg)
	x := (SzerokoscEkranu - rect.Dx()) / 2
	text.Draw(screen, msg, basicfont.Face7x13, x, y, clr)
}

func rysujWykresWynikow(screen *ebiten.Image) {
	tytuly := []string{"CALKOWITY CZAS (s)", "SREDNI CZAS OCZEKIWANIA (s)", "MAX CZAS OCZEKIWANIA (s)", "SREDNIA PODROZ (s)"}
	tytul := tytuly[indeksStronyWykresu]
	text.Draw(screen, "WYNIKI SYMULACJI ("+tytul+")", basicfont.Face7x13, 330, 50, color.White)
	if len(wyniki) == 0 {
		return
	}
	maxVal := 0.0
	for _, r := range wyniki {
		val := 0.0
		switch indeksStronyWykresu {
		case 0:
			val = r.CzasCalkowity.Seconds()
		case 1:
			val = r.SredniCzasOczekiwania.Seconds()
		case 2:
			val = r.MaxCzasOczekiwania.Seconds()
		case 3:
			val = r.SredniaPodroz.Seconds()
		}
		if val > maxVal {
			maxVal = val
		}
	}
	if maxVal == 0 {
		maxVal = 1.0
	}
	wysWykresu := 400.0
	bazaY := 600.0
	szerSlupka := 100.0
	odstep := 200.0
	startX := 200.0
	for i, r := range wyniki {
		val := 0.0
		switch indeksStronyWykresu {
		case 0:
			val = r.CzasCalkowity.Seconds()
		case 1:
			val = r.SredniCzasOczekiwania.Seconds()
		case 2:
			val = r.MaxCzasOczekiwania.Seconds()
		case 3:
			val = r.SredniaPodroz.Seconds()
		}
		wysSlupka := (val / maxVal) * wysWykresu
		x := startX + float64(i)*odstep
		y := bazaY - wysSlupka
		var kolSlupka color.Color
		if i == 0 {
			kolSlupka = color.RGBA{255, 100, 100, 255}
		}
		if i == 1 {
			kolSlupka = color.RGBA{255, 255, 100, 255}
		}
		if i == 2 {
			kolSlupka = color.RGBA{100, 255, 100, 255}
		}
		ebitenutil.DrawRect(screen, x, y, szerSlupka, wysSlupka, kolSlupka)
		text.Draw(screen, r.Nazwa, basicfont.Face7x13, int(x), int(bazaY)+20, color.White)
		text.Draw(screen, fmt.Sprintf("%.2fs", val), basicfont.Face7x13, int(x)+20, int(y)-10, color.White)
	}
	text.Draw(screen, "Nacisnij [SPACJA] - kolejny wykres | [ENTER] - powrot do menu", basicfont.Face7x13, 280, 700, color.RGBA{150, 150, 150, 255})
}

func rysujSymulacje(screen *ebiten.Image) {
	rysujSkrzyzowanie(screen)
	mutexSamochodow.Lock()
	sort.Slice(samochody, func(i, j int) bool { return samochody[i].Y < samochody[j].Y })
	for _, c := range samochody {
		rysujAuto(screen, c)
	}
	mutexSamochodow.Unlock()
	rysujSygnalizatory(screen)
}

func rysujPanelStatystyk(screen *ebiten.Image, tytul string) {
	szerP := 180.0
	wysP := 200.0
	x := float64(SzerokoscEkranu) - szerP
	y := float64(WysokoscEkranu) - wysP
	ebitenutil.DrawRect(screen, x, y, szerP, wysP, color.RGBA{0, 0, 0, 200})
	offX := int(x) + 10
	currentY := int(y) + 20
	text.Draw(screen, tytul, basicfont.Face7x13, offX, currentY, color.RGBA{100, 255, 100, 255})
	currentY += 20
	wyswietlanyCzas := 0.0
	if obecnyStan == StanJazdaWolna {
		wyswietlanyCzas = czasTrwaniaSym
	} else {
		wyswietlanyCzas = time.Since(czasStartuTestu).Seconds() * konfigPredkosc
	}
	text.Draw(screen, fmt.Sprintf("Czas: %.1fs", wyswietlanyCzas), basicfont.Face7x13, offX, currentY, color.White)
	currentY += 20
	pokazaneAuta := wszystkieAuta
	if obecnyStan != StanJazdaWolna {
		pokazaneAuta = testLicznikAut
	}
	text.Draw(screen, fmt.Sprintf("Ilosc aut: %d", pokazaneAuta), basicfont.Face7x13, offX, currentY, color.White)
	currentY += 20
	qN := pobierzRozmiarKolejki("N", "")
	qS := pobierzRozmiarKolejki("S", "")
	qE := pobierzRozmiarKolejki("E", "")
	qW := pobierzRozmiarKolejki("W", "")
	text.Draw(screen, fmt.Sprintf("Kolejki N:%d S:%d", qN, qS), basicfont.Face7x13, offX, currentY, color.White)
	currentY += 15
	text.Draw(screen, fmt.Sprintf("Kolejki E:%d W:%d", qE, qW), basicfont.Face7x13, offX, currentY, color.White)
	currentY += 30
}

func rysujSkrzyzowanie(screen *ebiten.Image) {
	szary := color.RGBA{80, 80, 80, 255}
	bialy := color.RGBA{255, 255, 255, 220}
	ebitenutil.DrawRect(screen, SrodekX-80, 0, 160, WysokoscEkranu, szary)
	ebitenutil.DrawRect(screen, 0, SrodekY-80, SzerokoscEkranu, 160, szary)
	ebitenutil.DrawRect(screen, SrodekX-2, 0, 4, Krawedz_N, bialy)
	ebitenutil.DrawRect(screen, SrodekX-2, Krawedz_S, 4, WysokoscEkranu-Krawedz_S, bialy)
	ebitenutil.DrawRect(screen, 0, SrodekY-2, Krawedz_W, 4, bialy)
	ebitenutil.DrawRect(screen, Krawedz_E, SrodekY-2, SzerokoscEkranu-Krawedz_E, 4, bialy)
	ebitenutil.DrawRect(screen, SrodekX-40, 0, 2, LiniaStop_N, bialy)
	ebitenutil.DrawRect(screen, SrodekX+40, LiniaStop_S, 2, WysokoscEkranu-LiniaStop_S, bialy)
	ebitenutil.DrawRect(screen, 0, SrodekY+40, LiniaStop_W, 2, bialy)
	ebitenutil.DrawRect(screen, LiniaStop_E, SrodekY-40, SzerokoscEkranu-LiniaStop_E, 2, bialy)
	ebitenutil.DrawRect(screen, SrodekX+40, 0, 2, Krawedz_N, bialy)
	ebitenutil.DrawRect(screen, SrodekX-40, Krawedz_S, 2, WysokoscEkranu-Krawedz_S, bialy)
	ebitenutil.DrawRect(screen, 0, SrodekY-40, Krawedz_W, 2, bialy)
	ebitenutil.DrawRect(screen, Krawedz_E, SrodekY+40, SzerokoscEkranu-Krawedz_E, 2, bialy)
	ebitenutil.DrawRect(screen, SrodekX-80, LiniaStop_N, 80, 2, bialy)
	ebitenutil.DrawRect(screen, SrodekX, LiniaStop_S, 80, 2, bialy)
	ebitenutil.DrawRect(screen, LiniaStop_W, SrodekY, 2, 80, bialy)
	ebitenutil.DrawRect(screen, LiniaStop_E, SrodekY-80, 2, 80, bialy)
}

func rysujAuto(screen *ebiten.Image, c *Samochod) {
	w, h := SzerokoscAuta, c.Dlugosc
	zwrot := ""
	if !c.Skrecil {
		switch c.Kierunek {
		case "N":
			zwrot = "DOL"
		case "S":
			zwrot = "GORA"
		case "E":
			zwrot = "LEWO"
			w, h = c.Dlugosc, SzerokoscAuta
		case "W":
			zwrot = "PRAWO"
			w, h = c.Dlugosc, SzerokoscAuta
		}
	} else {
		switch c.Kierunek {
		case "N":
			if c.Akcja == "LEFT" {
				zwrot = "PRAWO"
				w, h = c.Dlugosc, SzerokoscAuta
			} else {
				zwrot = "LEWO"
				w, h = c.Dlugosc, SzerokoscAuta
			}
		case "S":
			if c.Akcja == "LEFT" {
				zwrot = "LEWO"
				w, h = c.Dlugosc, SzerokoscAuta
			} else {
				zwrot = "PRAWO"
				w, h = c.Dlugosc, SzerokoscAuta
			}
		case "E":
			if c.Akcja == "LEFT" {
				zwrot = "DOL"
			} else {
				zwrot = "GORA"
			}
		case "W":
			if c.Akcja == "LEFT" {
				zwrot = "GORA"
			} else {
				zwrot = "DOL"
			}
		}
	}
	ebitenutil.DrawRect(screen, c.X, c.Y, w, h, c.Kolor)
	headColor := color.White
	brakeColor := color.RGBA{255, 0, 0, 255}
	var fx1, fy1, fx2, fy2, bx1, by1, bx2, by2 float64
	switch zwrot {
	case "DOL":
		fx1, fy1 = c.X+2, c.Y+h-4
		fx2, fy2 = c.X+w-6, c.Y+h-4
		bx1, by1 = c.X+2, c.Y
		bx2, by2 = c.X+w-6, c.Y
	case "GORA":
		fx1, fy1 = c.X+2, c.Y
		fx2, fy2 = c.X+w-6, c.Y
		bx1, by1 = c.X+2, c.Y+h-4
		bx2, by2 = c.X+w-6, c.Y+h-4
	case "PRAWO":
		fx1, fy1 = c.X+w-4, c.Y+2
		fx2, fy2 = c.X+w-4, c.Y+h-6
		bx1, by1 = c.X, c.Y+2
		bx2, by2 = c.X, c.Y+h-6
	case "LEWO":
		fx1, fy1 = c.X, c.Y+2
		fx2, fy2 = c.X, c.Y+h-6
		bx1, by1 = c.X+w-4, c.Y+2
		bx2, by2 = c.X+w-4, c.Y+h-6
	}
	ebitenutil.DrawRect(screen, fx1, fy1, 4, 4, headColor)
	ebitenutil.DrawRect(screen, fx2, fy2, 4, 4, headColor)
	if c.Predkosc < 0.5 {
		ebitenutil.DrawRect(screen, bx1, by1, 4, 4, brakeColor)
		ebitenutil.DrawRect(screen, bx2, by2, 4, 4, brakeColor)
	}
	if c.Akcja != "STRAIGHT" && !c.JestUprzyw {
		if (time.Now().UnixMilli()/300)%2 == 0 {
			tx, ty := 0.0, 0.0
			if c.Akcja == "LEFT" {
				switch zwrot {
				case "DOL":
					tx, ty = fx2, fy2
				case "GORA":
					tx, ty = fx1, fy1
				case "PRAWO":
					tx, ty = fx1, fy1
				case "LEWO":
					tx, ty = fx2, fy2
				}
			} else {
				switch zwrot {
				case "DOL":
					tx, ty = fx1, fy1
				case "GORA":
					tx, ty = fx2, fy2
				case "PRAWO":
					tx, ty = fx2, fy2
				case "LEWO":
					tx, ty = fx1, fy1
				}
			}
			ebitenutil.DrawRect(screen, tx, ty, 5, 5, color.RGBA{255, 165, 0, 255})
		}
	}
	if c.JestUprzyw {
		col := color.RGBA{255, 0, 0, 255}
		if (time.Now().UnixMilli()/150)%2 != 0 {
			col = color.RGBA{0, 0, 255, 255}
		}
		ebitenutil.DrawRect(screen, c.X+w/2-4, c.Y+h/2-4, 8, 8, col)
	}
}

func rysujSygnalizatory(screen *ebiten.Image) {
	rysujPojedynczySygnalizator(screen, "N", X_Kol2_Wew-10, SwiatloY_N, X_Kol1_Zew-10, SwiatloY_N)
	rysujPojedynczySygnalizator(screen, "S", X_Kol3_Wew-10, SwiatloY_S, X_Kol4_Zew-10, SwiatloY_S)
	rysujPojedynczySygnalizator(screen, "E", SwiatloX_E, Y_Wiersz2_Wew-10, SwiatloX_E, Y_Wiersz1_Zew-10)
	rysujPojedynczySygnalizator(screen, "W", SwiatloX_W, Y_Wiersz3_Wew-10, SwiatloX_W, Y_Wiersz4_Zew-10)
}

func rysujPojedynczySygnalizator(screen *ebiten.Image, kier string, lx, ly, rx, ry float64) {
	czerwony := color.RGBA{255, 0, 0, 255}
	zielony := color.RGBA{0, 255, 0, 255}
	zolty := color.RGBA{255, 215, 0, 255}
	lKol, pKol := czerwony, czerwony
	if karetkaAktywna {
		if kier == karetkaZrodlo {
			lKol, pKol = zielony, zielony
		} else {
			lKol, pKol = czerwony, czerwony
		}
	} else {
		if strings.Contains(obecnaFaza, "STRAIGHT") {
			if strings.HasPrefix(obecnaFaza, "NS") && (kier == "N" || kier == "S") {
				if strings.Contains(obecnaFaza, "YELLOW") {
					pKol = zolty
				} else {
					pKol = zielony
				}
			}
			if strings.HasPrefix(obecnaFaza, "EW") && (kier == "E" || kier == "W") {
				if strings.Contains(obecnaFaza, "YELLOW") {
					pKol = zolty
				} else {
					pKol = zielony
				}
			}
		}
		if strings.HasPrefix(obecnaFaza, kier+"_LEFT") {
			if strings.Contains(obecnaFaza, "YELLOW") {
				lKol = zolty
			} else {
				lKol = zielony
			}
		}
	}
	ebitenutil.DrawRect(screen, lx, ly, RozmiarSwiatla, RozmiarSwiatla, lKol)
	ebitenutil.DrawRect(screen, rx, ry, RozmiarSwiatla, RozmiarSwiatla, pKol)
}

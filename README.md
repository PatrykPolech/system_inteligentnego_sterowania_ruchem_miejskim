# Symulacja Inteligentnego Systemu Sterowania Ruchem Miejskim

Projekt to interaktywna symulacja skrzyżowania napisana w języku Go z wykorzystaniem biblioteki Ebitengine. Aplikacja umożliwia wizualizację ruchu drogowego oraz testowanie i porównywanie wydajności różnych algorytmów sterowania sygnalizacją świetlną.

<img width="1277" height="992" alt="image" src="https://github.com/user-attachments/assets/6bf13b68-b467-43f9-a078-3e7de7bded2d" />


## Główne funkcjonalności

* **Tryb Swobodnej Jazdy:** Nieskończona symulacja z losowo generowanym ruchem, pozwalająca na obserwację zachowania pojazdów w czasie rzeczywistym.
* **Porównanie Algorytmów:** Przeprowadzenie identycznego scenariusza testowego (dla 10, 20, 50 lub 100 pojazdów) dla trzech różnych sterowników świateł. Wyniki (całkowity czas, średni czas oczekiwania itp.) prezentowane są na czytelnych wykresach słupkowych.
* **Zróżnicowane pojazdy:** W symulacji biorą udział samochody osobowe, ciężarówki (o innej dynamice jazdy) oraz pojazdy uprzywilejowane (karetki), które wymuszają zmianę świateł na zielone.
* **Fizyka ruchu:** Pojazdy dynamicznie przyspieszają, hamują przed przeszkodami oraz zachowują bezpieczny odstęp od siebie i linii zatrzymania.

## Zaimplementowane algorytmy sterowania

1.  **Czasowe (Stałe):** Sygnalizacja zmienia się w sztywnych, z góry określonych odstępach czasu (np. 8 sekund dla jazdy prosto, 4 sekundy dla lewoskrętów).
2.  **Czujnik ilości pojazdów (Zachłanny):** Sterownik analizuje kolejki i zawsze przyznaje zielone światło temu kierunkowi, na którym czeka najwięcej aut.
3.  **Inteligentne:** Sterownik pomija fazy, w których nie ma oczekujących pojazdów, i dynamicznie skraca czas trwania zielonego światła, gdy ulica opustoszeje.

## Technologie

* Go 1.25.5
* Ebitengine v2.9.7

## Jak uruchomić projekt lokalnie

1. Sklonuj repozytorium na swój komputer:
   `git clone https://github.com/PatrykPolech/system_inteligentnego_sterowania_ruchem_miejskim.git`
2. Przejdź do folderu z projektem.
3. Pobierz wymagane zależności:
   `go mod tidy`
4. Uruchom grę:
   `go run .`

<img width="1273" height="988" alt="image" src="https://github.com/user-attachments/assets/0c9c3014-ba60-41e4-9bc6-46b616619b1c" />

<img width="1272" height="991" alt="image" src="https://github.com/user-attachments/assets/4b12a56e-f851-4dc7-9e3b-b04c09ba2782" />

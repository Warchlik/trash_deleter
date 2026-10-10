# Code review — trash_deleter

Data przeglądu: 2026-10-10.

## Ocena ogólna

Projekt jest sensowną bazą na małe narzędzie do portfolio, ale wymaga dopracowania bezpieczeństwa CLI oraz testów potwierdzających, czego program **nie usunie**. Największą przeszkodą nie jest styl kodu, tylko bezpieczeństwo operacji na danych.

Podczas przeglądu nie zmieniano plików projektu. Przeczytano cały kod i oba README. Działanie sprawdzono na danych w katalogu tymczasowym. Ten dokument zapisuje ustalenia — opisane poprawki nie zostały wdrożone.

## 1. Najważniejsze problemy w kodzie

### 1.1. Krytyczne: argument pozycyjny może wyłączyć `-dry-run`

**Miejsce:** `internal/config/config.go:83–98`.

Standardowy `flag.Parse()` kończy parsowanie flag po pierwszym argumencie pozycyjnym. Program nie sprawdza później pozostałych argumentów.

Potwierdzono, że:

```sh
trash_deleter -path ./folder unexpected -dry-run
```

**Usuwa pliki**, zamiast wykonać podgląd. `unexpected` zatrzymuje parsowanie, więc `DryRun` pozostaje `false`.

Także naturalna, choć nieobsługiwana składnia:

```sh
trash_deleter ./folder -dry-run
```

Ignoruje podaną ścieżkę i flagę — program działa na domyślnym `.`.

**Do poprawy:** odrzucać nieobsługiwane argumenty pozycyjne i kończyć działanie błędem **przed jakimkolwiek usuwaniem**. Jeśli program ma obsługiwać ścieżkę pozycyjną, należy zdefiniować jednoznaczny kontrakt składni.

To najważniejszy błąd w tym review.

### 1.2. Wysoki priorytet: niepełne skanowanie kończy się kodem sukcesu

**Miejsce:** `main.go:35–45, 74–76`.

Błędy skanowania i szukania pustych katalogów są tylko logowane. Kod wyjścia zależy wyłącznie od nieudanych operacji usuwania.

Potwierdzono scenariusz z katalogiem bez uprawnień do odczytu:

- Program wypisał `permission denied`.
- Nie przeskanował całego drzewa.
- Zakończył się kodem `0`.

Dla skryptu uruchamiającego narzędzie oznacza to sukces, mimo niekompletnego wyniku.

**Do poprawy:** uwzględnić błędy skanowania w końcowym statusie. Jawnie ustalić, czy częściowy skan pozwala na usuwanie znalezionych plików, czy przerywa operację. Udokumentować tę decyzję.

### 1.3. Wysoki priorytet: symlink jako katalog główny działa niespójnie

**Miejsce:** `main.go:16–22`, `internal/scanner/scanner.go:23, 98–113`.

`os.Stat()` podąża za symlinkiem, więc walidacja akceptuje dowiązanie do katalogu. `filepath.WalkDir()` nie przechodzi jednak przez taki symlink. Z kolei `os.ReadDir()` używane w `EmptyDirs()` odczytuje jego cel.

Potwierdzono:

- Skan przez symlink zgłosił `0 files`, mimo obecności pliku `junk.tmp`.
- Uruchomienie przez symlink z `-empty-dirs` usunęło pusty katalog wewnątrz jego celu.

**Do poprawy:** wybrać jedną politykę: odrzucać symlink jako root albo rozwiązywać go przed walidacją i przekazywać ten sam root do obu operacji. Osobno opisać zachowanie symlinków znalezionych wewnątrz drzewa.

### 1.4. Wysoki priorytet: program usuwa również pliki specjalne

**Miejsce:** `internal/scanner/scanner.go:32–55`.

Skaner wyklucza katalogi, ale nie ogranicza kandydatów do zwykłych plików. Pasująca nazwa wystarczy również dla FIFO czy symlinka.

Potwierdzono, że FIFO o nazwie `service.tmp` zostało usunięte.

**Do poprawy:** określić dozwolone typy obiektów. Dla narzędzia opisanego jako usuwające śmieciowe pliki bezpiecznym domyślnym wyborem są zwykłe pliki. Usuwanie dowiązań i obiektów specjalnych powinno być świadomą, opisaną decyzją.

### 1.5. Średni priorytet: `-max-depth` ma różną semantykę w obu skanerach

**Miejsce:** `internal/scanner/scanner.go:39–40, 105–114`.

Przy `-max-depth 1` skan plików pomija katalogi bezpośrednio pod rootem. Skan pustych katalogów nadal do nich wchodzi.

Potwierdzono, że:

```sh
trash_deleter -path ./folder -max-depth 1 -empty-dirs
```

Usuwa pusty podkatalog, choć skan plików nie wchodzi na ten poziom.

**Do poprawy:** ujednolicić definicję głębokości albo wyraźnie opisać osobne zasady. Obecne README tłumaczy limit przez położenie plików, ale nie rozstrzyga zachowania pustych katalogów.

### 1.6. Średni priorytet: raport „B freed” nie oznacza odzyskanego miejsca

**Miejsce:** `internal/cleaner/cleaner.go:37–39`, `main.go:58`.

Program sumuje `info.Size()` kandydatów. To logiczny rozmiar plików, nie pomiar zwolnionej przestrzeni.

- W `-dry-run` niczego nie zwalnia.
- Przeniesienie do kosza nie jest równoznaczne z odzyskaniem miejsca.
- Hardlinki i pliki sparse dodatkowo rozdzielają rozmiar pliku od zajętości dysku.

**Do poprawy:** raportować „łączny rozmiar plików” lub „rozmiar kandydatów”, zamiast deklarować odzyskane bajty. Projekt nie potrzebuje skomplikowanego pomiaru przestrzeni — wystarczy uczciwa nazwa.

## 2. Co dopracować, żeby projekt dobrze wyglądał w portfolio

### 2.1. Bezpieczny kontrakt działania

Domyślne natychmiastowe usuwanie z `.` jest udokumentowane, więc **nie jest ukrytym błędem**, ale jest ryzykowną decyzją produktową.

Warto rozważyć domyślny podgląd i jawne potwierdzenie wykonania, np. flagę `-apply`. Niezależnie od wyboru:

- Błędne argumenty muszą oznaczać brak zmian.
- Użytkownik powinien widzieć katalog docelowy i tryb operacji.
- README powinno zaznaczać, że `*.tmp` nie gwarantuje, że plik jest zbędny.

### 2.2. Testy zachowania — przede wszystkim ochrony danych

Nie znaleziono plików `*_test.go`. Dla programu usuwającego dane to istotna luka.

Najbardziej wartościowe przypadki:

- `dry-run` pozostawia pliki i katalogi bez zmian.
- Błędne argumenty nie uruchamiają usuwania.
- Ignorowane katalogi pozostają nietknięte.
- Limit głębokości działa na granicach.
- Katalog z zachowanym plikiem nie jest usuwany.
- Błędy skanowania wpływają na wynik operacji.
- Symlinki i pliki specjalne zachowują się zgodnie z ustaloną polityką.

Należy używać `t.TempDir()`. Testy nie powinny dotykać prawdziwych danych użytkownika ani jego systemowego kosza.

### 2.3. Łatwiejsza testowalność CLI

**Miejsce:** `internal/config/config.go:36–106`.

`Parse()` korzysta z globalnego pakietu `flag` i wywołuje `os.Exit()`. To utrudnia izolowane sprawdzanie kolejnych zestawów argumentów.

**Kierunek:** lokalny `flag.FlagSet`, argumenty przekazywane do parsera, wynik jako konfiguracja i błąd. Decyzja o zakończeniu procesu zostaje w `main`.

Nie potrzeba frameworka CLI ani rozbudowanej architektury.

### 2.4. Powtarzalne budowanie i czytelna prezentacja

Nie znaleziono konfiguracji CI ani licencji. Warto dodać:

- CI uruchamiające build, testy i `go vet`.
- Informację o wymaganej wersji Go i obsługiwanych systemach.
- Macierz systemów, jeśli projekt deklaruje wieloplatformowość.
- Opis ograniczeń `-trash`, symlinków, głębokości i kodów wyjścia w **obu** README.
- Krótki przykład „drzewo przed → dry-run → drzewo po”.
- Licencję, jeśli projekt ma być publicznie wykorzystywany.
- Ignorowanie lokalnej binarki — `.gitignore` obecnie zawiera tylko `test_folder`.

Jeśli projekt ma obsługiwać instalację przez `go install`, warto użyć ścieżki modułu odpowiadającej publicznemu repozytorium zamiast samego `trash_deleter`.

## 3. Co już jest dobre

- Czytelny podział na `config`, `scanner` i `cleaner`.
- Mało kodu i brak niepotrzebnych abstrakcji.
- Walidacja niepoprawnych globów i ujemnej głębokości przed usuwaniem.
- Podstawowe czyszczenie, ignorowanie `.git`/`node_modules` i `dry-run` działały w sprawdzonych scenariuszach.
- `os.Remove()` dla katalogów nie usuwa rekursywnie przypadkowo zachowanej zawartości.
- README zawiera konkretne przykłady, flagi i ostrzeżenie.

Nie ma potrzeby rozbudowywania projektu w większą aplikację. Mały CLI jest wystarczający do portfolio, jeśli pokazuje świadome decyzje, ochronę danych i dobre testy.

## 4. Kolejność prac

1. **Naprawić parser argumentów** — obecnie można przypadkiem ominąć `dry-run`.
2. Ustalić politykę błędów skanowania, symlinków i typów plików.
3. Dodać testy chroniące przed niezamierzonym usuwaniem.
4. Ujednolicić głębokość i poprawić raportowanie rozmiaru.
5. Dodać CI i dopracować dokumentację oraz sposób instalacji.

## 5. Weryfikacja i ograniczenia przeglądu

- Build (`go build -o <tymczasowa-binarka> .`) zakończył się poprawnie.
- `go vet ./...` zakończył się poprawnie.
- Uruchomiono CLI na tymczasowych danych i odtworzono opisane problemy parsera, kodu wyjścia przy błędzie odczytu, symlinków, plików specjalnych oraz limitu głębokości.
- Sprawdzono podstawowy `dry-run`, czyszczenie rekursywne, zachowanie zwykłego niedopasowanego pliku, pomijanie `.git` i `node_modules`, limit głębokości dla plików oraz odrzucenie błędnego globu i ujemnej głębokości.
- Nie uruchamiano `-trash`, więc działanie integracji z systemowym koszem pozostaje niezweryfikowane. Uwaga o „B freed” dla kosza wynika z analizy kodu i znaczenia tej operacji, nie z uruchomienia tego trybu.
- Katalog tymczasowy z binarką i danymi scenariuszy został usunięty po przeglądzie.
- Odwołania do numerów linii dotyczą stanu kodu w chwili przeglądu.

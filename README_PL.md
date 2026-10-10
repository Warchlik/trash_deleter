# trash_deleter

[English](README.md)

Małe narzędzie CLI napisane w Go, które znajduje i usuwa śmieciowe pliki (`.DS_Store`, `Thumbs.db`, `*.tmp` i podobne) z drzewa katalogów.

> **Uwaga:** bez `-dry-run` pliki są usuwane od razu. Najpierw uruchom z `-dry-run`, żeby zobaczyć, co zostałoby usunięte.

## Budowanie

```sh
go build -o trash_deleter .
```

## Użycie

```sh
./trash_deleter -path ./jakis/katalog -dry-run
./trash_deleter -path ./jakis/katalog
```

## Flagi

| Flaga | Opis |
|---|---|
| `-path` | Katalog do przeskanowania (domyślnie `.`) |
| `-dry-run` | Tylko wypisz, co zostałoby usunięte, niczego nie zmieniaj |
| `-patterns` | Dodatkowe wzorce nazw plików po przecinku, dodawane do domyślnych (np. `"*.log,backup_*"`) |
| `-ignore` | Dodatkowe nazwy katalogów do pominięcia po przecinku, dodawane do domyślnych (np. `"vendor,dist"`) |
| `-empty-dirs` | Usuń też puste katalogi, w tym te, które opróżniły się po czyszczeniu |
| `-trash` | Przenieś do systemowego kosza zamiast usuwać na stałe |
| `-max-depth` | Głębokość skanowania: `1` = tylko pliki bezpośrednio w `-path`, `0` = bez limitu (domyślnie) |
| `-v`, `-verbose` | Wypisz każdy usunięty plik i katalog |

## Domyślne ustawienia

- Wzorce: `*:Zone.Identifier`, `.DS_Store`, `Thumbs.db`, `*.tmp`, `*~`
- Pomijane katalogi: `.git`, `node_modules`

Wzorce są dopasowywane tylko do nazwy pliku i używają składni glob. Dla prefiksu użyj `backup_*`.

## Uwagi

- Kod wyjścia to `1`, jeśli czegoś nie udało się usunąć.
- Na macOS `-trash` działa przez Findera. Pierwsze uruchomienie może być wolne i poprosić o zgodę Automation.

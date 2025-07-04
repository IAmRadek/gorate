# ZADANIE REKRUTACYJNE

Na wstępie chciałbym podziękować za zainterosowanie się naszą ofertą pracy oraz poświęcony czas aby zaznajomić się z poniszym zadaniem. W przypadku braku nieścisłości lub nie zrozumienia treści zadania, prosimy o kontakt. Tymczasem przystępujemy do treści zadania.

## CEL OGÓLNY

Celem ogólnym zadania jest stworzenie prostej aplikacji w języku GO, przy użyciu frameworka Gin-Gonic. O strukturze umożliwiający jej dalszy rozwój, przy zachowaniu dobrych praktyk programistycznych oraz języka GO.

## CELE ZADANIA

- [X] Utworzenie prostej aplikacji w Go z użyciem frameworka gin-gonic.
- [X] Dodanie jednego (lub dwóch) endpointu HTTP (opis poniżej).
- [X] Dodanie testów do endpoint'u `/rates`.
- [X] Skonteneryzowanie aplikacji przy użyciu Dockera.
- [X] Wynikiem swojej pracy wyślij na utworzone repozytorium na GitHub i podziel się linkiem.

## Opis endpointów

### Endpoint do pobierania kursów wymiany walut - cel wymagany

**Metoda HTTP:** `GET`

**Endpoint**: `/rates`

**Opis:**

Endpoint wymaga jednego parametru `currencies`, w którym ma przyjąć listę walut (wymagane są conajmniej dwie).
Wykorzystując API serwisu `openexchangerates.org` pobierz kursy wymiany walut dla żądanych walut, używając amerykańskiego dolara jako waluty bazowej.
Następnie przelicz brakujące kursy walut i zwróć `status 200` zaś w ciele odpowiedzi, tablicę ze wszystkimi możliwymi parami walut oraz kursami ich wymiany.
W przypadku braku/pustego parametru `currencies`, podaniu tylko jednej waluty, lub błędu zwróconego przez API serwisu `openexchangerates.org` zwróć sam `status 400` z pustym ciałem odpowiedzi.

#### Przykład #1

**Request:** `GET /rates?currencies=USD,GBP,EUR`

**Response:**

> Status: `200`

```
[
  { "from": "USD", "to": "GBP", "rate": 1.0 },
  { "from": "GBP", "to": "USD", "rate": 1.0 },
  { "from": "USD", "to": "EUR", "rate": 1.0 },
  { "from": "EUR", "to": "USD", "rate": 1.0 },
  { "from": "EUR", "to": "GBP", "rate": 1.0 },
  { "from": "GBP", "to": "EUR", "rate": 1.0 },
]
```

#### Przykład #2

**Request:** `GET /rates?currencies=GBP,EUR`

**Response (JSON):**

> Status: `200`

```
[
  { "from": "EUR", "to": "GBP", "rate": 1.0 },
  { "from": "GBP", "to": "EUR", "rate": 1.0 },
]
```

#### Przykład #3

**Request:** `GET /rates?currencies=GBP`

**Response:**

> Status: `400`

---

### Endpoint do przeliczania jednej kryptowaluty na inną - cel opcjonalny
[TASK.md](TASK.md)
**Metoda HTTP:** `GET`

**Endpoint:** `/exchange`

**Opis:**

Endpoint wymaga trzech parametrów:

- `from` - waluta krypto, którą chcemy wymienić
- `to` - waluta krypto, która chcemy otrzymać
- `amount` - kwotę krypto jaką chcemy wymienić

Na bazie powyższych parametrów oraz poniższej tabeli z kursami walut wylicz ile kryptowaluty otrzyma użytkownik, po zamianie jednej kryptowaluty na inną.
W przypadku poprawnego zapytania, zwracamy `status 200` obiekt zawierający informację, które kryptowaluty wymieniany oraz ilość potencjalnie otrzymanej kryptowaluty *(precyzja ilości ma być równa precyzji zawartej tabeli w kolumnie `decimal places` dla danej kryptowaluty)*. 
Ponad to obsługujemy tylko tokeny zawarte w tabeli. W przypadku otrzymania kryptowaluty spoza listy, lub któryś parametr jest pusty zwracamy `status 400` oraz puste ciało odpowiedzi.

#### Kursy walut

| CryptoCurrency | Decimal places | Rate (to USD) |
| ----------- | ----------- | ----------- |
| BEER | 18 | 0.00002461$
| FLOKI | 18 | 0.0001428$
| GATE| 18 | 6.87$
| USDT | 6 | 0.999$
| WBTC | 8 | 57,037.22$

#### Przykład #1

**Request:** `GET /exchange?from=WBTC&to=USDT&amount=1.0`

**Response (JSON):**

> Status: `200`

```
{ "from": "WBTC", "to": "USDT", "amount": 57613.353535 }
```

#### Przykład #2

**Request:** `GET /exchange?from=USDT&to=BEER&amount=1.0`

**Response (JSON):**

> Status: `200`

```
{ "from": "USDT", "to": "USDT", "amount": 7009.810931379558830532 }
```

#### Przykład #3

**Request:** `GET /exchange?from=MATIC&to=GATE&amount=0.999`

**Response:**

> Status: `400`

#### Przykład #4

**Request:** `GET /exchange?from=USDT&to=GATE`

**Response:**

> Status: `400`

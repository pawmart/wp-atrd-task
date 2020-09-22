# Zadanie  GWP

Należy napisać prosty serwer HTTP. Zadanie powinno zostać zrealizowane w dowolnym z następujących języków: C/C++/python/GO. Zachęcamy do próby rozwiązania tego zadania w GO.

## Specyfikacja

- [Features](./features/secrets)
- [Swagger](./api/swagger/swagger.yml)

## Wskazówki

- Pełna dowolność w wyborze bibliotek oraz rodzaju bazy danych lub jej braku.
- Rozwiązanie zadania proszę załączyć jako 'merge request'. 

## Uruchomienie projektu

Projekt składa się z serwera HTTP oraz bazy danych Postgres. Do integracji rozwiązania uzyto Docker-compose.

Uruchomienie projektu:
```
$ cd wp-atrd-task
$ docker-compose up
```

Projekt domyślnie uruchamia się na adresie localhost:8080.

Przykłady poleceń do interakcji z API:

- Stworzenie nowego Secret:
`curl -XPOST localhost:8080/v1/secret/ -d "secret=asdfasdfasdfasdf&expireAfterViews=5&expireAfter=1"`

- Odczytanie istniejącego Secret (odpowiedź w formacie JSON):
`curl localhost:8080/v1/secret/c5772968-24a5-4914-a966-8292cb54585c`

- Odczytanie istniejącego Secret (odpowiedź w formacie XML):
`curl -H "Accept: application/xml" localhost:8080/v1/secret/c5772968-24a5-4914-a966-8292cb54585c`

## Uruchomienie testów

Zakładana jest działająca instalacja Go oraz poprawnie skonfigurowany PATH zawierający $GOPATH/bin.

```
$ go get github.com/cucumber/godog/cmd/godog@v0.10.0
$ cd wp-atrd-task
$ docker-compose up -d
$ godog
```
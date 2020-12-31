# Zadanie  GWP

Należy napisać prosty serwer HTTP. Zadanie powinno zostać zrealizowane w dowolnym z następujących języków: C/C++/python/GO. Zachęcamy do próby rozwiązania tego zadania w GO.

## Specyfikacja

- [Features](./features/secrets)
- [Swagger](./api/swagger/swagger.yml)

## Wskazówki

- Pełna dowolność w wyborze bibliotek oraz rodzaju bazy danych lub jej braku.
- Rozwiązanie zadania proszę załączyć jako 'merge request'.

## Realizacja

Zadanie zostało wykonane przy użuciu języka Go, framework'a go-kit oraz bazy danych MongoDB.

### Uruchomienie

Do poprawnego działania usługi niezbędne jest odpalenie zarówno procesu serwera API jak i usługi zajmującej się czyszczeniem bazy danych.

Całość odpalamy za pomocą komendy `docker-compose up`. API dostępne będzie na porcie `:3000`.


### Uwagi

* Feature testy zostały przeniesione do katalogu `/cmd/secretsd`

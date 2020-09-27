# Zadanie  GWP

Należy napisać prosty serwer HTTP. Zadanie powinno zostać zrealizowane w dowolnym z następujących języków: C/C++/python/GO. Zachęcamy do próby rozwiązania tego zadania w GO.

## Specyfikacja

- [Features](./features/secrets)
- [Swagger](./api/swagger/swagger.yml)

## Wskazówki

- Pełna dowolność w wyborze bibliotek oraz rodzaju bazy danych lub jej braku.
- Rozwiązanie zadania proszę załączyć jako 'merge request'. 

## Uruchamianie

- ```docker-compose build && docker-compose up```
- przykładowy `curl -X POST "http://localhost:8080/v1/secret" -H  "accept: application/json" -H  "Content-Type: application/x-www-form-urlencoded" -d "secret=alamakota&expireAfterViews=0&expireAfter=1"`

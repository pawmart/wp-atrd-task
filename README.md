# Zadanie  GWP

Należy napisać prosty serwer HTTP. Zadanie powinno zostać zrealizowane w dowolnym z następujących języków: C/C++/python/GO. Zachęcamy do próby rozwiązania tego zadania w GO.

## Specyfikacja

- [Features](./features/secrets)
- [Swagger](./api/swagger/swagger.yml)

## Wskazówki

- Pełna dowolność w wyborze bibliotek oraz rodzaju bazy danych lub jej braku.
- Rozwiązanie zadania proszę załączyć jako 'merge request'.

## Dev docs

### Running dev env

```bash
docker-compose up
# REST API http://localhost:3000
# Redis web interface http://admin:admin@localhost:6380
```

### Running tests

```bash
cd cmd/server
godog ../../features
``````

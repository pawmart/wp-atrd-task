# Zadanie  GWP

Należy napisać prosty serwer HTTP. Zadanie powinno zostać zrealizowane w dowolnym z następujących języków: C/C++/python/GO. Zachęcamy do próby rozwiązania tego zadania w GO.

## Specyfikacja

- [Features](./features/secrets)
- [Swagger](./api/swagger/swagger.yml)

## Wskazówki

- Pełna dowolność w wyborze bibliotek oraz rodzaju bazy danych lub jej braku.
- Rozwiązanie zadania proszę załączyć jako 'merge request'. 


## Środowisko deweloperskie
- odpalenie serwera w katalogu projektu, serwer nasłuchuje na porcie `3001`
```
go run cmd/secret_server/secret_server.go
```

## Testy

- Jednostkowe
```
go test ./...
``` 
- Behawioralne
```
cd cmd/secret_server && godog ../../features
```

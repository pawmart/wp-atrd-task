# Zadanie  GWP

Należy napisać prosty serwer HTTP. Zadanie powinno zostać zrealizowane w dowolnym z następujących języków: C/C++/python/GO. Zachęcamy do próby rozwiązania tego zadania w GO.

## Specyfikacja

- [Features](./features/secrets)
- [Swagger](./api/swagger/swagger.yml)

## Wskazówki

- Pełna dowolność w wyborze bibliotek oraz rodzaju bazy danych lub jej braku.
- Rozwiązanie zadania proszę załączyć jako 'merge request'. 

# Rozwiązanie zadania - informacje
Do odpalenia skryptu potrzebne są pakiety: 
- urllib3 - pakiet do (pip install urllib3)
- peewee - ORM do ogarnięcia bazy danych (pip install peewee) 

Wykorzystano bazę danych Sqlite3

Start serwera komendą "python main.py"
Serwer działa na (adres:port): localhost:8000
Zapytanie POST testowane za pomocą skryptu "requesttest.py"
Skrypt pisany i testowany na Windows 10.
Python 3.7.4


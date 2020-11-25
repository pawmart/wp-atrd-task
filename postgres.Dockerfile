FROM library/postgres
COPY scripts/init.sql /docker-entrypoint-initdb.d/
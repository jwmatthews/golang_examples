sudo -u "postgres" createdb gwp
psql -h localhost -U testuser -f setup.sql -d gwp


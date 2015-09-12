sudo -u "postgres" psql -c 'drop database if exists gwp'
sudo -u "postgres" createdb gwp
psql -h localhost -U testuser -f setup.sql -d gwp


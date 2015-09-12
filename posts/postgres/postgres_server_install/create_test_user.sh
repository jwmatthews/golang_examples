sudo -u "postgres" createdb test
sudo -u "postgres" psql -c "CREATE ROLE testuser WITH SUPERUSER LOGIN PASSWORD 'test';"



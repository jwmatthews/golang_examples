sudo -u "postgres" psql -c 'drop database if exists gwp'
sudo -u "postgres" createdb gwp
# Gorm will create the schema and handle migrations


POSTGRESQL_CONF=/var/lib/pgsql/data/postgresql.conf
PG_HBA_CONF=/var/lib/pgsql/data/pg_hba.conf

sudo dnf install -y postgresql postgresql-server postgresql-devel postgresql-libs

sudo systemctl enable postgresql
sudo postgresql-setup --initdb --unit postgresql

sudo sed -i "s/#listen_addresses = 'localhost' */listen_addresses = '*'/" $POSTGRESQL_CONF

sudo sed -i "s/#port = 5432 */port = 5432/" $POSTGRESQL_CONF


sudo cp pg_hba.conf.sample $PG_HBA_CONF
sudo chown postgres:postgres $PG_HBA_CONF
sudo systemctl restart postgresql

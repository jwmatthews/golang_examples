DEFAULT_ZONE=`sudo firewall-cmd --get-default-zone`
sudo firewall-cmd --permanent --zone ${DEFAULT_ZONE} --add-service postgresql
sudo firewall-cmd --reload

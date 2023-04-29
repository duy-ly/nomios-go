FROM arm64v8/mysql:oracle
MAINTAINER duy.1407@gmail.com

COPY ./conf.d /etc/mysql/conf.d
COPY ./inventory.sql /docker-entrypoint-initdb.d/

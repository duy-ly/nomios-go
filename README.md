# nomios-go
A high performance CDC application.

# How to use
- Nomios require binlog_format=row
- Should enable gtid_mode config for supporting failover
- Table must have primary keys
- When change table schema
    - Should not change column type (Column data maybe corrupt)
    - Must add new column to the end (must not add column to beginning or middle)
    - Must not remove existing column

# How to run Nomios at local
## Build and run docker containers
1. Run `maven package` to build jar file
2. Copy jar file to `docker` directory
```
cp target/nomios-1.0-SNAPSHOT-fat.jar docker/
```
2. Run docker-compose
```
cd docker
docker-compose up -d --build
```
3. Add Nomios connector
```
curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" http://localhost:8083/connectors/ -d @nomios-connector.json
```
4. Create table person and add data
```
CREATE TABLE person (
    id int unsigned NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    age INTEGER NULL DEFAULT 10,
    createdAt DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);

INSERT INTO person(name,age) VALUES ('Tony Stark', 43);
```

## How to check kafka state
1. Check kafka records
```
docker-compose exec kafka /opt/kafka/bin/kafka-console-consumer.sh --bootstrap-server kafka:9092 --from-beginning --property print.key=true --topic nomios.catalog_db.person
```
6. Check last offset
```
docker-compose exec kafka /opt/kafka/bin/kafka-console-consumer.sh --bootstrap-server kafka:9092 --from-beginning --property print.key=true --topic connect-offsets
```

## Monitoring
1. Open grafana at: http://localhost:3000/
2. Enter `admin/admin` to login
3. Open Nomios Dashboard at the home page

## Failover demonstration
1. Try to stop mysql1 container
```
docker-compose stop mysql1
```
2. Wait for Nomios connector for reconnecting to mysql source
3. Check server UUID of mysql2 then adding more data to mysql2
```
SHOW GLOBAL VARIABLES LIKE 'server_uuid';
INSERT INTO person(name,age) VALUES ('Tony Stark', 43);
```
4. Wait for records from Kafka topic, verify that `gtid` of new record is set to the `UUID` of mysql2

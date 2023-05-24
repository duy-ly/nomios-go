CREATE TABLE datetime_test (
   id int unsigned NOT NULL AUTO_INCREMENT,
   datetime_col DATETIME NULL,
   datetime_col_2 DATETIME(3) NULL,
   datetime_col_3 DATETIME(6) NULL,
   timestamp_col TIMESTAMP NULL,
   timestamp_col_2 TIMESTAMP(3) NULL,
   timestamp_col_3 TIMESTAMP(6) NULL,
   date_col DATE NULL,
   time_col TIME NULL,
   year_col YEAR NULL,
   PRIMARY KEY (`id`)
);

INSERT INTO `datetime_test`(datetime_col, datetime_col_2, datetime_col_3,
    timestamp_col, timestamp_col_2, timestamp_col_3, date_col, time_col, year_col)
    VALUES ('2022-02-02 12:30:15', '2022-02-02 12:30:15.300', '2022-02-02 12:30:15.222222',
            '2022-02-02 12:30:15', '2022-02-02 12:30:15.300', '2022-02-02 12:30:15.222222', '2022-02-02', '12:30:15', '2022');
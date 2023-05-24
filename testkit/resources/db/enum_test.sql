CREATE TABLE `test_stations` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(500) COLLATE utf8_unicode_ci NOT NULL,
    `type` enum('station', 'post_office') COLLATE utf8_unicode_ci NOT NULL DEFAULT 'station',
    PRIMARY KEY (`id`)
);

INSERT INTO test_stations (`name`, `type`) values ( 'Ha Noi', 'station' );
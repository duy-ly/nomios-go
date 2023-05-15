CREATE TABLE schema_test (
    id int unsigned NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    age INTEGER NULL DEFAULT 10,
    `type` enum('station', 'post_office')  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'station',
    `json_col` JSON NULL,
    salary DECIMAL(9,2) NOT NULL,
    created_at DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);

CREATE TABLE person (
    id int unsigned NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    age INTEGER NULL DEFAULT 10,
    created_at DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);

CREATE TABLE `event_store` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `message_id` varchar(40) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `type` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `topic` varchar(54) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `routing_key` varchar(14) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `source` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `payload` mediumtext COLLATE utf8mb4_unicode_ci,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_message_id` (`message_id`),
    KEY `idx_routing_key` (`routing_key`),
    KEY `idx_topic` (`topic`),
    KEY `idx_type` (`type`)
) ENGINE=InnoDB AUTO_INCREMENT=10994446 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO person(name,age,created_at,updated_at) VALUES ('Georgia',30,'2022-02-25 10:50:54','2022-02-25 10:50:54');

INSERT INTO `event_store` (`id`, `message_id`, `type`, `topic`, `routing_key`, `source`, `payload`, `created_at`) VALUES
    (10994458, 'd868dea7-5f58-4eba-a8a3-0953c6bdaa00', 'Catalog\\Events\\ProductAttributeUpdated', 'catalog.product_event_bus', '2460072', 'Catalog\\Models\\Product', '{\"sourceId\":\"2460072\",\"createdAt\":null,\"payLoad\":null,\"createdBy\":\"erp@tiki.vn\",\"productAttributes\":{\"reason\":null,\"product_seller_id\":null,\"priceComponentParam\":null,\"canBePurchased\":null,\"minSaleQuantity\":null,\"type\":null,\"useConfigMinSaleQuantity\":null,\"sellerId\":1,\"product_master_id\":null,\"useConfigMaxSaleQuantity\":null,\"eavAttributeValues\":{\"support_p2h_delivery\":{\"code\":\"support_p2h_delivery\",\"type\":\"int\",\"value\":\"0\"},\"bulky\":{\"code\":\"bulky\",\"type\":\"int\",\"value\":\"0\"},\"is_fresh\":{\"code\":\"is_fresh\",\"type\":\"int\",\"value\":\"0\"}},\"price\":null,\"newCategories\":[],\"categories\":[],\"id\":2460072,\"sku\":\"9154846413409\",\"inventoryType\":null,\"quantity\":null,\"listReasonIds\":[],\"visibility\":null,\"entityType\":\"master_simple\",\"productSetId\":null,\"maxSaleQuantity\":null,\"version\":null,\"preorderDate\":null,\"productCode\":null,\"name\":null,\"poType\":null,\"newImages\":[],\"minCode\":\"9154846413409\",\"primaryCategoryId\":null,\"sellerProductCode\":null,\"status\":null},\"attributeValues\":{\"support_p2h_delivery\":{\"code\":\"support_p2h_delivery\",\"type\":\"int\",\"value\":\"0\"},\"bulky\":{\"code\":\"bulky\",\"type\":\"int\",\"value\":\"0\"},\"is_fresh\":{\"code\":\"is_fresh\",\"type\":\"int\",\"value\":\"0\"}},\"messageId\":\"d868dea7-5f58-4eba-a8a3-0953c6bdaa00\",\"eventType\":\"Catalog\\\\Events\\\\ProductAttributeUpdated\",\"version\":2}', '2022-06-23 14:02:48');

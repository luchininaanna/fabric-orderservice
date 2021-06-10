CREATE TABLE IF NOT EXISTS `order`
(
    id BINARY(16) NOT NULL,
    status INTEGER NOT NULL,
    cost FLOAT NOT NULL,
    address VARCHAR (255) NOT NULL,
    created_at DATETIME NOT NULL,
    closed_at DATETIME,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE order_item
(
    order_id BINARY(16) NOT NULL,
    fabric_id BINARY(16) NOT NULL,
    quantity INT,
    PRIMARY KEY (order_id, fabric_id),
    FOREIGN KEY (order_id)
        REFERENCES `order`(id)
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
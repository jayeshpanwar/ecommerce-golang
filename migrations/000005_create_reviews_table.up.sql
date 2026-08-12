CREATE TABLE reviews (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    deleted_at DATETIME(3) NULL,

    user_id BIGINT UNSIGNED NOT NULL,

    product_id BIGINT UNSIGNED NOT NULL,

    order_item_id BIGINT UNSIGNED NOT NULL,
    UNIQUE KEY uk_reviews_order_item (order_item_id),

    rating TINYINT UNSIGNED NOT NULL,

    comment TEXT,

    INDEX idx_reviews_deleted_at (deleted_at),
    INDEX idx_reviews_user_id (user_id),
    INDEX idx_reviews_product_id (product_id),

    CONSTRAINT fk_reviews_user
        FOREIGN KEY (user_id)
        REFERENCES users(id),

    CONSTRAINT fk_reviews_product
        FOREIGN KEY (product_id)
        REFERENCES products(id),

    CONSTRAINT fk_reviews_order_item
        FOREIGN KEY (order_item_id)
        REFERENCES order_items(id)
);
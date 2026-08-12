CREATE TABLE review_images (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    deleted_at DATETIME(3) NULL,

    review_id BIGINT UNSIGNED NOT NULL,

    image_url VARCHAR(500) NOT NULL,

    INDEX idx_review_images_deleted_at (deleted_at),
    INDEX idx_review_images_review_id (review_id),

    CONSTRAINT fk_review_images_review
        FOREIGN KEY (review_id)
        REFERENCES reviews(id)
);
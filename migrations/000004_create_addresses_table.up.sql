CREATE TABLE addresses (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    deleted_at DATETIME(3) NULL,

    user_id BIGINT UNSIGNED NOT NULL,

    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,

    address_line_1 VARCHAR(255) NOT NULL,
    address_line_2 VARCHAR(255),

    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    pincode VARCHAR(20) NOT NULL,

    is_default BOOLEAN DEFAULT FALSE,

    INDEX idx_addresses_deleted_at (deleted_at),
    INDEX idx_addresses_user_id (user_id),

    CONSTRAINT fk_addresses_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
);
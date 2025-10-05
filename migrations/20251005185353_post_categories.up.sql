-- Migration UP: post_categories
CREATE TABLE post_categories (
    category_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    `description` VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
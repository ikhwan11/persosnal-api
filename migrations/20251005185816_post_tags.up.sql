-- Migration UP: post_tags
CREATE TABLE post_tags (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tag_id BIGINT NOT NULL,
    post_id BIGINT NOT NULL,
    CONSTRAINT tag_id FOREIGN KEY (tag_id) REFERENCES tags(tag_id),
    CONSTRAINT post_id FOREIGN KEY (post_id) REFERENCES posts(post_id)
);
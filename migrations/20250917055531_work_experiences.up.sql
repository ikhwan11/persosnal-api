-- Migration UP: work_experiences
CREATE TABLE work_experiences (
    we_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    cv_id BIGINT NOT NULL,
    logo VARCHAR(255),
    nama_perusahaan VARCHAR(150) NOT NULL,
    periode VARCHAR(100),
    tugas TEXT,
    pencapaian TEXT,
    CONSTRAINT fk_we_cv FOREIGN KEY (cv_id) REFERENCES cv(cv_id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
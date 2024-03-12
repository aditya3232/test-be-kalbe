CREATE TABLE position
(
    position_id INT NOT NULL AUTO_INCREMENT,
    department_id INT NULL DEFAULT NULL,
    position_name VARCHAR(200) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(200) NULL DEFAULT NULL,
    updated_at DATETIME NULL DEFAULT NULL,
    updated_by VARCHAR(200) NULL DEFAULT NULL,
    deleted_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (position_id),
    FOREIGN KEY (department_id) REFERENCES department(department_id)
) ENGINE = InnoDB;
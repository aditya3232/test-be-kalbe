CREATE TABLE employee
(
    employee_id INT NOT NULL AUTO_INCREMENT,
    employee_code VARCHAR(200) NOT NULL,
    employee_name VARCHAR(200) NOT NULL,
    password VARCHAR(255),
    department_id INT NULL DEFAULT NULL,
    position_id INT NULL DEFAULT NULL,
    superior INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(200) NULL DEFAULT NULL,
    updated_at DATETIME NULL DEFAULT NULL,
    updated_by VARCHAR(200) NULL DEFAULT NULL,
    deleted_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (employee_id),
    FOREIGN KEY (department_id) REFERENCES department(department_id),
    FOREIGN KEY (position_id) REFERENCES `position`(position_id)
) ENGINE = InnoDB;
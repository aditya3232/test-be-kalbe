CREATE TABLE attendance
(
    attendance_id INT NOT NULL AUTO_INCREMENT,
    employee_id INT NULL DEFAULT NULL,
    location_id INT NULL DEFAULT NULL,
    absent_in DATETIME NULL DEFAULT NULL,
    absent_out DATETIME NULL DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(200) NULL DEFAULT NULL,
    updated_at DATETIME NULL DEFAULT NULL,
    updated_by VARCHAR(200) NULL DEFAULT NULL,
    deleted_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (attendance_id),
    FOREIGN KEY (employee_id) REFERENCES employee(employee_id),
    FOREIGN KEY (location_id) REFERENCES location(location_id)
) ENGINE = InnoDB;
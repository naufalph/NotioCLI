CREATE TABLE TODO_LIST (
    ID INT AUTO_INCREMENT PRIMARY KEY,               -- Primary key for tasks
    DESCRIPTION VARCHAR(255) NOT NULL,                      -- Task name
    DUE_DATE DATE DEFAULT NULL,                      -- Optional due date for the task
    CREATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Timestamp when the task is created
    UPDATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- Timestamp for updates
    IS_DELETE BOOL DEFAULT FALSE,
    STATUS VARCHAR(40) NOT NULL DEFAULT 'NOT STARTED'
) ENGINE=InnoDB;
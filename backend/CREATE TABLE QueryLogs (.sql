CREATE TABLE QueryLogs (
    log_id INT IDENTITY(1,1) PRIMARY KEY,
    ioc NVARCHAR(255) NOT NULL,
    last_lookup DATETIME DEFAULT GETDATE(),
    result_count INT,
    user_name NVARCHAR(255) NOT NULL
);

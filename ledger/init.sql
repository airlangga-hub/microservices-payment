CREATE USER 'ledger_user'@'localhost' IDENTIFIED BY 'password';

CREATE DATABASE IF NOT EXISTS ledger;

GRANT ALL PRIVILEGES ON ledger.* TO 'ledger_user'@'localhost';

USE ledger;

CREATE TABLE IF NOT EXISTS ledger (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_id VARCHAR(255),
    user_id VARCHAR(255),
    amount INT,
    operation VARCHAR(255),
    date VARCHAR(255)
);
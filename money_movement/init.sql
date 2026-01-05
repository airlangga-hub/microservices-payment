CREATE USER 'money_movement_user'@'localhost' IDENTIFIED BY 'password';

CREATE DATABASE IF NOT EXISTS money_movement;

GRANT ALL PRIVILEGES ON money_movement.* TO 'money_movement_user'@'localhost';

USE money_movement;

CREATE TABLE IF NOT EXISTS wallets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL UNIQUE,
    wallet_type VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cents INT NOT NULL DEFAULT 0,
    account_type VARCHAR(255) NOT NULL,
    wallet_id INT NOT NULL,
    FOREIGN KEY (wallet_id) REFERENCES wallets (id)
);

CREATE TABLE IF NOT EXISTS transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    pid VARCHAR(255) NOT NULL,
    src_user_id VARCHAR(255) NOT NULL,
    dst_user_id VARCHAR(255) NOT NULL,
    src_wallet_id INT NOT NULL,
    dst_wallet_id INT NOT NULL,
    src_account_id INT NOT NULL,
    dst_account_id INT NOT NULL,
    src_account_type VARCHAR(255) NOT NULL,
    dst_account_type VARCHAR(255) NOT NULL,
    final_dst_merchant_wallet_id INT NOT NULL,
    amount INT NOT NULL,
    INDEX(pid)
);

-- merchant and customer wallets
INSERT INTO wallets (user_id, wallet type) VALUES ('user1@email.com', 'CUSTOMER');
INSERT INTO wallets (user_id, wallet type) VALUES ('merchant_id', 'MERCHANT');

-- customer accounts
INSERT INTO accounts (cents, account_type, wallet_id) VALUES (5000000, 'DEFAULT', 1);
INSERT INTO accounts (cents, account_type, wallet_id) VALUES (0, 'PAYMENT', 1);

-- merchant account
INSERT INTO accounts (cents, account_type, wallet_id) VALUES (0, 'INCOMING', 2);
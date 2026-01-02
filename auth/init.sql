CREATE USER 'auth_user'@'localhost' IDENTIFIED BY 'password';

CREATE DATABASE IF NOT EXISTS users;

GRANT ALL PRIVILEGES ON users.* TO 'auth_user'@'localhost';
MYSQL Queries:

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    mobile_no VARCHAR(9) UNIQUE NOT NULL,
    created_at DATETIME  DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME  DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE cards (
    id INT AUTO_INCREMENT PRIMARY KEY,
    card_no VARCHAR(7) UNIQUE NOT NULL,
    user_id INT,
    created_at DATETIME  DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME  DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE deliveries (
    id INT AUTO_INCREMENT PRIMARY KEY,
    card_id INT,
    status ENUM('card_created','pickedup', 'delivery_attempt_1', 'delivery_attempt_2', 'delivered', 'returned') DEFAULT 'card_created',
    created_at DATETIME  DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME  DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (card_id) REFERENCES cards(id)
);

CREATE TABLE delivery_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    delivery_id INT,
    comment TEXT,
    created_at DATETIME  DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME  DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (delivery_id) REFERENCES deliveries(id)
);

#########################################################################################################


1. Create User API

    curl -X POST http://localhost:8080/users \
    -H "Content-Type: application/json" \
    -d '{
    "mobile_no": "123456789"
    }'

2. Create Card API

    a. Using user_id:


    curl -X POST http://localhost:8080/cards \
    -H "Content-Type: application/json" \
    -d '{
    "user_id": 1
    }'

    b. Using mobile_no:


    curl -X POST http://localhost:8080/cards \
    -H "Content-Type: application/json" \
    -d '{
    "mobile_no": "123456789"
    }'

3. Get Card Status API

    a. Using mobile_no:

    curl -X GET "http://localhost:8080/card-status?mobile_no=123456789"

    b. Using card_id:

    curl -X GET "http://localhost:8080/card-status?card_id=GEB0863"

##########################################################################################

Execution Steps:

A) Non Docker Steps :

go build .
go run main.go


B) Using Docker steps:

1) docker build -t card-service .

2) docker run -p 8080:8080 card-service



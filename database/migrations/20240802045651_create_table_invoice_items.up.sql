CREATE TABLE invoice_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    invoiceId BIGINT NOT NULL,
    itemId BIGINT NOT NULL,
    quantity BIGINT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (invoiceId) REFERENCES invoices(id)
       ON DELETE RESTRICT
       ON UPDATE CASCADE,
    FOREIGN KEY (itemId) REFERENCES items(id)
       ON DELETE RESTRICT
       ON UPDATE CASCADE
);
CREATE TABLE invoices (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    customerId BIGINT NOT NULL,
    invoiceId VARCHAR(225) DEFAULT NULL,
    subject TEXT DEFAULT NULL,
    dueDate DATE DEFAULT NULL,
    status ENUM('Paid', 'Unpaid')  DEFAULT 'Unpaid',
    totalAmount DECIMAL(10, 2) NOT NULL,
    itemCount INT NOT NULL,
    issueDate DATE DEFAULT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE (invoiceId)
);

ALTER TABLE invoices
    ADD CONSTRAINT fk_invoice_customer FOREIGN KEY (customerId) REFERENCES customers(id)
    ON DELETE RESTRICT
    ON UPDATE CASCADE;

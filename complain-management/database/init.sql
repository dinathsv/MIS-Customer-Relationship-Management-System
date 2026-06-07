CREATE TABLE IF NOT EXISTS complaints (
    id SERIAL PRIMARY KEY,
    customer_id VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sample Data for Demo
INSERT INTO complaints (customer_id, title, description, status) VALUES
('john.doe@example.com', 'Delayed Delivery', 'My order #A993 has not arrived yet.', 'Pending'),
('jane.smith@example.com', 'Defective Product', 'The item received has a cracked screen.', 'In Progress'),
('alice.johnson@example.com', 'Billing Error', 'Charged twice for the subscription.', 'Resolved');
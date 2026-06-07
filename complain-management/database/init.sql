CREATE TABLE IF NOT EXISTS complaints (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sample Data for Demo
INSERT INTO complaints (customer_id, title, description, status) VALUES
(101, 'Delayed Delivery', 'My order #A993 has not arrived yet.', 'Pending'),
(102, 'Defective Product', 'The item received has a cracked screen.', 'In Progress'),
(103, 'Billing Error', 'Charged twice for the subscription.', 'Resolved');
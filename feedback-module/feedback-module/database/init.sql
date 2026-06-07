-- ============================================================
-- Customer Feedback Module - PostgreSQL Initialization Script
-- CRM Enterprise System
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Drop table if exists (for clean re-init)
DROP TABLE IF EXISTS customer_feedback;

-- ============================================================
-- MAIN TABLE: customer_feedback
-- ============================================================
CREATE TABLE customer_feedback (
    feedback_id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id   VARCHAR(64)  NOT NULL,
    rating        SMALLINT     NOT NULL CHECK (rating BETWEEN 1 AND 5),
    category      VARCHAR(64)  NOT NULL CHECK (
        category IN (
            'Product Quality',
            'Customer Support',
            'Delivery & Shipping',
            'Pricing & Value',
            'Website & App',
            'Returns & Refunds',
            'Other'
        )
    ),
    comments      TEXT,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Indexes for CRM query patterns
CREATE INDEX idx_cf_customer_id  ON customer_feedback (customer_id);
CREATE INDEX idx_cf_category     ON customer_feedback (category);
CREATE INDEX idx_cf_rating       ON customer_feedback (rating);
CREATE INDEX idx_cf_created_at   ON customer_feedback (created_at DESC);

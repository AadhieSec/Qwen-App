# Database Schema Design

## Overview

This document contains the complete database schema for ShopMonitor, including tables, indexes, constraints, and relationships.

## Database: PostgreSQL 15+

### Extensions Required

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";  -- For text search
```

---

## Tables

### users

Stores user account information.

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255),  -- NULL for OAuth users
    name VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(500),
    timezone VARCHAR(50) DEFAULT 'UTC',
    currency VARCHAR(3) DEFAULT 'USD',
    language VARCHAR(10) DEFAULT 'en',
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

### sessions

Manages user sessions with refresh tokens.

```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(500) NOT NULL UNIQUE,
    user_agent VARCHAR(500),
    ip_address INET,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX idx_sessions_revoked_at ON sessions(revoked_at);

-- Auto-cleanup of expired sessions (run via cron/pg_cron)
-- DELETE FROM sessions WHERE expires_at < NOW() OR revoked_at IS NOT NULL;
```

### providers

Registry of supported shopping websites.

```sql
CREATE TABLE providers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    base_url VARCHAR(500) NOT NULL,
    logo_url VARCHAR(500),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'maintenance', 'deprecated')),
    capabilities JSONB NOT NULL DEFAULT '{}',
    config JSONB DEFAULT '{}',
    rate_limit INTEGER DEFAULT 60,  -- requests per minute
    last_check TIMESTAMPTZ,
    health_status VARCHAR(20) DEFAULT 'unknown',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_providers_slug ON providers(slug);
CREATE INDEX idx_providers_status ON providers(status);
CREATE INDEX idx_providers_health ON providers(health_status);

-- Insert initial providers
INSERT INTO providers (name, slug, base_url, capabilities) VALUES
('Amazon', 'amazon', 'https://www.amazon.com', '{"price": true, "stock": true, "variants": true, "coupons": false, "delivery": true}'),
('Flipkart', 'flipkart', 'https://www.flipkart.com', '{"price": true, "stock": true, "variants": true, "coupons": true, "delivery": true}'),
('Myntra', 'myntra', 'https://www.myntra.com', '{"price": true, "stock": true, "variants": true, "coupons": true, "delivery": true}'),
('Ajio', 'ajio', 'https://www.ajio.com', '{"price": true, "stock": true, "variants": true, "coupons": true, "delivery": true}'),
('Snitch', 'snitch', 'https://www.thesnitch.co.in', '{"price": true, "stock": true, "variants": true, "coupons": false, "delivery": true}'),
('Nike', 'nike', 'https://www.nike.com', '{"price": true, "stock": true, "variants": true, "coupons": false, "delivery": true}'),
('Adidas', 'adidas', 'https://www.adidas.com', '{"price": true, "stock": true, "variants": true, "coupons": false, "delivery": true}'),
('H&M', 'hnm', 'https://www.hm.com', '{"price": true, "stock": true, "variants": true, "coupons": false, "delivery": true}'),
('Zara', 'zara', 'https://www.zara.com', '{"price": true, "stock": true, "variants": true, "coupons": false, "delivery": true}'),
('Levi''s', 'levis', 'https://www.levi.com', '{"price": true, "stock": true, "variants": true, "coupons": false, "delivery": true}'),
('Decathlon', 'decathlon', 'https://www.decathlon.in', '{"price": true, "stock": true, "variants": true, "coupons": false, "delivery": true}'),
('Croma', 'croma', 'https://www.croma.com', '{"price": true, "stock": true, "variants": true, "coupons": true, "delivery": true}'),
('Apple', 'apple', 'https://www.apple.com', '{"price": true, "stock": true, "variants": true, "coupons": false, "delivery": true}'),
('Samsung', 'samsung', 'https://www.samsung.com', '{"price": true, "stock": true, "variants": true, "coupons": false, "delivery": true}'),
('Boat', 'boat', 'https://www.boat-lifestyle.com', '{"price": true, "stock": true, "variants": true, "coupons": true, "delivery": true}'),
('Nothing', 'nothing', 'https://nothing.tech', '{"price": true, "stock": true, "variants": true, "coupons": false, "delivery": true}'),
('Swiggy', 'swiggy', 'https://www.swiggy.com', '{"price": true, "stock": false, "variants": false, "coupons": true, "delivery": true, "food": true}'),
('Zomato', 'zomato', 'https://www.zomato.com', '{"price": true, "stock": false, "variants": false, "coupons": true, "delivery": true, "food": true}'),
('Blinkit', 'blinkit', 'https://blinkit.com', '{"price": true, "stock": true, "variants": false, "coupons": true, "delivery": true, "food": true}'),
('Zepto', 'zepto', 'https://www.zeptonow.com', '{"price": true, "stock": true, "variants": false, "coupons": true, "delivery": true, "food": true}'),
('BigBasket', 'bigbasket', 'https://www.bigbasket.com', '{"price": true, "stock": true, "variants": false, "coupons": true, "delivery": true, "food": true}');
```

### products

Core product information from shopping websites.

```sql
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider_id UUID NOT NULL REFERENCES providers(id) ON DELETE RESTRICT,
    external_id VARCHAR(255) NOT NULL,  -- Product ID from the website
    url VARCHAR(2000) NOT NULL,
    title VARCHAR(1000) NOT NULL,
    brand VARCHAR(255),
    category VARCHAR(255),
    subcategory VARCHAR(255),
    description TEXT,
    main_image_url VARCHAR(500),
    rating DECIMAL(3,2) CHECK (rating >= 0 AND rating <= 5),
    review_count INTEGER DEFAULT 0,
    seller VARCHAR(255),
    shipping_cost DECIMAL(10,2) DEFAULT 0,
    return_policy TEXT,
    metadata JSONB DEFAULT '{}',  -- Provider-specific data
    last_crawled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    UNIQUE(provider_id, external_id)
);

-- Indexes
CREATE INDEX idx_products_provider_id ON products(provider_id);
CREATE INDEX idx_products_external_id ON products(external_id);
CREATE INDEX idx_products_url ON products(url);
CREATE INDEX idx_products_brand ON products(brand);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_title_trgm ON products USING gin(title gin_trgm_ops);
CREATE INDEX idx_products_created_at ON products(created_at);

-- Trigger for updated_at
CREATE TRIGGER update_products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

### product_variants

Product variants (sizes, colors, storage, etc.).

```sql
CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(255),
    variant_type VARCHAR(50) NOT NULL,  -- 'size', 'color', 'storage', etc.
    variant_value VARCHAR(100) NOT NULL,  -- 'M', 'Red', '128GB', etc.
    variant_key VARCHAR(100),  -- Normalized value for comparison ('M', 'red', etc.)
    price DECIMAL(12,2) NOT NULL,
    mrp DECIMAL(12,2),
    discount_percent DECIMAL(5,2) DEFAULT 0,
    available BOOLEAN DEFAULT TRUE,
    stock_level VARCHAR(50),  -- 'in_stock', 'low_stock', 'out_of_stock'
    quantity_available INTEGER,
    image_url VARCHAR(500),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);
CREATE INDEX idx_product_variants_type_value ON product_variants(variant_type, variant_value);
CREATE INDEX idx_product_variants_available ON product_variants(available);
CREATE INDEX idx_product_variants_price ON product_variants(price);

-- Trigger for updated_at
CREATE TRIGGER update_product_variants_updated_at
    BEFORE UPDATE ON product_variants
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

### product_images

Product images gallery.

```sql
CREATE TABLE product_images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url VARCHAR(500) NOT NULL,
    alt_text VARCHAR(255),
    position INTEGER DEFAULT 0,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_product_images_product_id ON product_images(product_id);
CREATE INDEX idx_product_images_position ON product_images(position);
```

### user_products

User's wishlist/favorites and tracked products.

```sql
CREATE TABLE user_products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    nickname VARCHAR(255),  -- User's custom name for the product
    tags TEXT[] DEFAULT '{}',
    notes TEXT,
    is_favorite BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    UNIQUE(user_id, product_id)
);

-- Indexes
CREATE INDEX idx_user_products_user_id ON user_products(user_id);
CREATE INDEX idx_user_products_product_id ON user_products(product_id);
CREATE INDEX idx_user_products_favorite ON user_products(is_favorite);
CREATE INDEX idx_user_products_tags ON user_products USING gin(tags);
```

### monitors

Active monitoring configurations for products.

```sql
CREATE TABLE monitors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'paused', 'stopped', 'error')),
    check_interval INTEGER NOT NULL DEFAULT 3600,  -- seconds (1 hour default)
    target_price DECIMAL(12,2),  -- Alert when price <= this
    target_discount DECIMAL(5,2),  -- Alert when discount >= this %
    desired_sizes TEXT[] DEFAULT '{}',  -- ['M', 'L', 'XL']
    desired_colors TEXT[] DEFAULT '{}',  -- ['Red', 'Blue']
    desired_variants JSONB DEFAULT '[]',  -- Complex variant preferences
    max_price DECIMAL(12,2),  -- Stop monitoring if price > this
    min_discount DECIMAL(5,2),  -- Minimum discount threshold
    delivery_pincode VARCHAR(10),  -- Indian pincode for delivery check
    seller_preference TEXT,  -- Preferred seller
    notification_channels TEXT[] DEFAULT '{"email"}',  -- ['email', 'telegram', 'discord']
    notify_on_price_drop BOOLEAN DEFAULT TRUE,
    notify_on_stock BOOLEAN DEFAULT TRUE,
    notify_on_variant BOOLEAN DEFAULT TRUE,
    notify_on_coupon BOOLEAN DEFAULT TRUE,
    notify_on_delivery BOOLEAN DEFAULT TRUE,
    last_check_at TIMESTAMPTZ,
    next_check_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    paused_at TIMESTAMPTZ,
    error_message TEXT,
    consecutive_failures INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_monitors_user_id ON monitors(user_id);
CREATE INDEX idx_monitors_product_id ON monitors(product_id);
CREATE INDEX idx_monitors_status ON monitors(status);
CREATE INDEX idx_monitors_next_check ON monitors(next_check_at) WHERE status = 'active';
CREATE INDEX idx_monitors_user_status ON monitors(user_id, status);

-- Trigger for updated_at
CREATE TRIGGER update_monitors_updated_at
    BEFORE UPDATE ON monitors
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

### price_history

Historical price observations.

```sql
CREATE TABLE price_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    monitor_id UUID NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    variant_id UUID REFERENCES product_variants(id) ON DELETE SET NULL,
    price DECIMAL(12,2) NOT NULL,
    mrp DECIMAL(12,2),
    discount_percent DECIMAL(5,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',
    seller VARCHAR(255),
    in_stock BOOLEAN DEFAULT TRUE,
    condition VARCHAR(50) DEFAULT 'new',  -- 'new', 'refurbished', 'used'
    observed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes - Critical for performance
CREATE INDEX idx_price_history_monitor_id ON price_history(monitor_id);
CREATE INDEX idx_price_history_product_id ON price_history(product_id);
CREATE INDEX idx_price_history_observed_at ON price_history(observed_at DESC);
CREATE INDEX idx_price_history_monitor_observed ON price_history(monitor_id, observed_at DESC);
CREATE INDEX idx_price_history_product_observed ON price_history(product_id, observed_at DESC);

-- Partitioning for large datasets (optional, for production)
-- CREATE TABLE price_history_y2024 PARTITION OF price_history
--     FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');
```

### stock_history

Historical stock availability.

```sql
CREATE TABLE stock_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    monitor_id UUID NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    variant_id UUID REFERENCES product_variants(id) ON DELETE SET NULL,
    available BOOLEAN NOT NULL,
    quantity_available INTEGER,
    stock_level VARCHAR(50),  -- 'in_stock', 'low_stock', 'out_of_stock'
    observed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_stock_history_monitor_id ON stock_history(monitor_id);
CREATE INDEX idx_stock_history_product_id ON stock_history(product_id);
CREATE INDEX idx_stock_history_observed_at ON stock_history(observed_at DESC);
```

### coupons

Discovered coupon codes.

```sql
CREATE TABLE coupons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    code VARCHAR(100) NOT NULL,
    title VARCHAR(255),
    description TEXT,
    discount_type VARCHAR(50),  -- 'percentage', 'fixed', 'free_shipping'
    discount_value DECIMAL(10,2),
    min_order_value DECIMAL(12,2),
    max_discount DECIMAL(12,2),
    valid_from TIMESTAMPTZ,
    valid_until TIMESTAMPTZ,
    terms TEXT[] DEFAULT '{}',
    applicable_categories TEXT[] DEFAULT '{}',
    applicable_products TEXT[] DEFAULT '{}',
    bank_offers TEXT[] DEFAULT '{}',
    active BOOLEAN DEFAULT TRUE,
    verified BOOLEAN DEFAULT FALSE,
    usage_count INTEGER DEFAULT 0,
    success_rate DECIMAL(5,2),
    discovered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_verified_at TIMESTAMPTZ,
    
    UNIQUE(provider_id, code)
);

-- Indexes
CREATE INDEX idx_coupons_product_id ON coupons(product_id);
CREATE INDEX idx_coupons_provider_id ON coupons(provider_id);
CREATE INDEX idx_coupons_active ON coupons(active);
CREATE INDEX idx_coupons_valid_until ON coupons(valid_until);
CREATE INDEX idx_coupons_code ON coupons(code);
```

### delivery_history

Delivery availability history.

```sql
CREATE TABLE delivery_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    monitor_id UUID NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    pincode VARCHAR(10) NOT NULL,
    available BOOLEAN NOT NULL,
    delivery_type VARCHAR(50),  -- 'standard', 'express', 'same_day', 'one_day'
    estimated_days INTEGER,
    estimated_date DATE,
    shipping_cost DECIMAL(10,2) DEFAULT 0,
    free_shipping_threshold DECIMAL(12,2),
    observed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_delivery_history_monitor_id ON delivery_history(monitor_id);
CREATE INDEX idx_delivery_history_product_id ON delivery_history(product_id);
CREATE INDEX idx_delivery_history_pincode ON delivery_history(pincode);
CREATE INDEX idx_delivery_history_observed_at ON delivery_history(observed_at DESC);
```

### alerts

Generated alerts for users.

```sql
CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    monitor_id UUID NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,  -- 'price_drop', 'back_in_stock', 'coupon', 'delivery'
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    severity VARCHAR(20) DEFAULT 'info' CHECK (severity IN ('info', 'warning', 'critical')),
    data JSONB DEFAULT '{}',  -- Alert-specific data
    previous_value DECIMAL(12,2),  -- Previous price/stock level
    current_value DECIMAL(12,2),  -- Current price/stock level
    read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMPTZ,
    dismissed BOOLEAN DEFAULT FALSE,
    dismissed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_alerts_user_id ON alerts(user_id);
CREATE INDEX idx_alerts_monitor_id ON alerts(monitor_id);
CREATE INDEX idx_alerts_type ON alerts(type);
CREATE INDEX idx_alerts_read ON alerts(read);
CREATE INDEX idx_alerts_created_at ON alerts(created_at DESC);
CREATE INDEX idx_alerts_user_unread ON alerts(user_id, read) WHERE read = FALSE;
```

### notifications

Notification queue for delivery to channels.

```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    alert_id UUID REFERENCES alerts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel VARCHAR(50) NOT NULL,  -- 'email', 'telegram', 'discord', 'slack', 'desktop', 'webhook'
    recipient VARCHAR(500),  -- Email, chat ID, webhook URL
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'sending', 'sent', 'failed', 'skipped')),
    payload JSONB NOT NULL,  -- Notification content
    error_message TEXT,
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    scheduled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sent_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_channel ON notifications(channel);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_scheduled_at ON notifications(scheduled_at) WHERE status = 'pending';
```

### notification_logs

Audit log of notification deliveries.

```sql
CREATE TABLE notification_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    notification_id UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    channel VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    request_data JSONB,
    response_data JSONB,
    error_message TEXT,
    latency_ms INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_notification_logs_notification_id ON notification_logs(notification_id);
CREATE INDEX idx_notification_logs_channel ON notification_logs(channel);
CREATE INDEX idx_notification_logs_created_at ON notification_logs(created_at DESC);
```

### jobs

Background job queue.

```sql
CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(100) NOT NULL,  -- 'check_price', 'check_stock', 'send_notification', etc.
    priority INTEGER DEFAULT 0,  -- Higher = more urgent
    payload JSONB NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')),
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    scheduled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    error_message TEXT,
    worker_id UUID,
    result JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_jobs_type ON jobs(type);
CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_priority_scheduled ON jobs(priority DESC, scheduled_at ASC) WHERE status = 'pending';
CREATE INDEX idx_jobs_scheduled_at ON jobs(scheduled_at) WHERE status = 'pending';
CREATE INDEX idx_jobs_worker_id ON jobs(worker_id);
```

### workers

Worker process registry.

```sql
CREATE TABLE workers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hostname VARCHAR(255) NOT NULL,
    pid INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'idle' CHECK (status IN ('idle', 'busy', 'stopping', 'dead')),
    current_job_id UUID REFERENCES jobs(id) ON DELETE SET NULL,
    jobs_completed BIGINT DEFAULT 0,
    jobs_failed BIGINT DEFAULT 0,
    last_heartbeat TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Indexes
CREATE INDEX idx_workers_status ON workers(status);
CREATE INDEX idx_workers_heartbeat ON workers(last_heartbeat);

-- Unique constraint for hostname+pid (prevent duplicate workers)
CREATE UNIQUE INDEX idx_workers_hostname_pid ON workers(hostname, pid);
```

### user_settings

User preferences and settings.

```sql
CREATE TABLE user_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    theme VARCHAR(20) DEFAULT 'system' CHECK (theme IN ('light', 'dark', 'system')),
    timezone VARCHAR(50) DEFAULT 'UTC',
    currency VARCHAR(3) DEFAULT 'USD',
    language VARCHAR(10) DEFAULT 'en',
    date_format VARCHAR(20) DEFAULT 'YYYY-MM-DD',
    time_format VARCHAR(20) DEFAULT '24h',
    notification_prefs JSONB DEFAULT '{
        "email": true,
        "desktop": true,
        "telegram": false,
        "discord": false,
        "slack": false,
        "webhook": false
    }',
    alert_frequency VARCHAR(20) DEFAULT 'instant' CHECK (alert_frequency IN ('instant', 'hourly', 'daily', 'weekly')),
    quiet_hours_start TIME,
    quiet_hours_end TIME,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Trigger for updated_at
CREATE TRIGGER update_user_settings_updated_at
    BEFORE UPDATE ON user_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

### saved_searches

User's saved product searches.

```sql
CREATE TABLE saved_searches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    query JSONB NOT NULL,  -- Search parameters
    filters JSONB DEFAULT '{}',
    notify_on_new BOOLEAN DEFAULT FALSE,
    last_run_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_saved_searches_user_id ON saved_searches(user_id);

-- Trigger for updated_at
CREATE TRIGGER update_saved_searches_updated_at
    BEFORE UPDATE ON saved_searches
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

### audit_logs

Security and activity audit trail.

```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    ip_address INET,
    user_agent VARCHAR(500),
    request_data JSONB,
    response_status INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_ip ON audit_logs(ip_address);
```

---

## Views

### v_product_summary

Convenient view for product information with current price.

```sql
CREATE VIEW v_product_summary AS
SELECT 
    p.id,
    p.provider_id,
    pr.name as provider_name,
    pr.slug as provider_slug,
    p.external_id,
    p.url,
    p.title,
    p.brand,
    p.category,
    p.main_image_url,
    p.rating,
    p.review_count,
    COALESCE(v.price, 0) as current_price,
    COALESCE(v.mrp, 0) as mrp,
    COALESCE(v.discount_percent, 0) as discount_percent,
    v.available as in_stock,
    (SELECT COUNT(*) FROM monitors m WHERE m.product_id = p.id) as monitor_count,
    (SELECT MIN(ph.price) FROM price_history ph WHERE ph.product_id = p.id) as lowest_price,
    (SELECT MAX(ph.price) FROM price_history ph WHERE ph.product_id = p.id) as highest_price,
    p.created_at,
    p.updated_at
FROM products p
LEFT JOIN providers pr ON p.provider_id = pr.id
LEFT JOIN LATERAL (
    SELECT price, mrp, discount_percent, available
    FROM product_variants
    WHERE product_id = p.id
    ORDER BY price ASC
    LIMIT 1
) v ON true;
```

### v_monitor_status

Monitor status with last observation.

```sql
CREATE VIEW v_monitor_status AS
SELECT 
    m.id,
    m.user_id,
    m.product_id,
    p.title as product_title,
    p.main_image_url,
    m.status,
    m.check_interval,
    m.target_price,
    m.target_discount,
    m.last_check_at,
    m.next_check_at,
    m.consecutive_failures,
    m.error_message,
    COALESCE(latest.price, 0) as last_price,
    COALESCE(latest.discount_percent, 0) as last_discount,
    COALESCE(latest.in_stock, false) as last_stock_status,
    CASE 
        WHEN latest.price > 0 AND prev.price > 0 THEN 
            ((prev.price - latest.price) / prev.price * 100)
        ELSE 0
    END as price_change_percent,
    m.created_at
FROM monitors m
JOIN products p ON m.product_id = p.id
LEFT JOIN LATERAL (
    SELECT price, discount_percent, in_stock, observed_at
    FROM price_history
    WHERE monitor_id = m.id
    ORDER BY observed_at DESC
    LIMIT 1
) latest ON true
LEFT JOIN LATERAL (
    SELECT price
    FROM price_history
    WHERE monitor_id = m.id
    ORDER BY observed_at DESC
    OFFSET 1
    LIMIT 1
) prev ON true;
```

---

## Functions

### fn_calculate_savings

Calculate total savings for a user.

```sql
CREATE OR REPLACE FUNCTION fn_calculate_savings(p_user_id UUID)
RETURNS TABLE (
    total_savings DECIMAL(12,2),
    alerts_triggered BIGINT,
    best_deal_product UUID,
    best_deal_savings DECIMAL(12,2)
) AS $$
BEGIN
    RETURN QUERY
    WITH price_drops AS (
        SELECT 
            a.product_id,
            (a.previous_value - a.current_value) as savings
        FROM alerts a
        WHERE a.user_id = p_user_id
        AND a.type = 'price_drop'
        AND a.previous_value > a.current_value
    )
    SELECT 
        COALESCE(SUM(savings), 0)::DECIMAL(12,2) as total_savings,
        COUNT(*)::BIGINT as alerts_triggered,
        (SELECT product_id FROM price_drops ORDER BY savings DESC LIMIT 1) as best_deal_product,
        (SELECT MAX(savings) FROM price_drops)::DECIMAL(12,2) as best_deal_savings
    FROM price_drops;
END;
$$ LANGUAGE plpgsql STABLE;
```

### fn_cleanup_old_data

Maintenance function to clean up old historical data.

```sql
CREATE OR REPLACE FUNCTION fn_cleanup_old_data(p_retention_days INTEGER DEFAULT 365)
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    -- Delete old price history
    DELETE FROM price_history
    WHERE observed_at < NOW() - (p_retention_days || ' days')::INTERVAL;
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    
    -- Delete old stock history
    DELETE FROM stock_history
    WHERE observed_at < NOW() - (p_retention_days || ' days')::INTERVAL;
    
    -- Delete old delivery history
    DELETE FROM delivery_history
    WHERE observed_at < NOW() - (p_retention_days || ' days')::INTERVAL;
    
    -- Delete old read alerts
    DELETE FROM alerts
    WHERE read = TRUE
    AND created_at < NOW() - (p_retention_days / 2 || ' days')::INTERVAL;
    
    -- Delete dead workers (no heartbeat in 5 minutes)
    DELETE FROM workers
    WHERE last_heartbeat < NOW() - INTERVAL '5 minutes';
    
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;
```

---

## Initial Data

### Default Admin User (for development)

```sql
-- Password: admin123 (hash this properly in production!)
INSERT INTO users (email, password_hash, name, email_verified) VALUES
('admin@shopmonitor.local', '$2a$10$...', 'Admin User', TRUE);

INSERT INTO user_settings (user_id, theme, timezone, currency) VALUES
((SELECT id FROM users WHERE email = 'admin@shopmonitor.local'), 'dark', 'Asia/Kolkata', 'INR');
```

---

## Migration Strategy

Use GORM AutoMigrate for development and generate SQL migrations for production:

```bash
# Generate migration
go run cmd/migrate/main.go create create_initial_schema

# Run migrations
go run cmd/migrate/main.go up
```

See `backend/internal/database/migrations/` for migration files.

---

## Performance Considerations

1. **Partitioning**: For high-volume installations, partition `price_history`, `stock_history`, and `audit_logs` by date.

2. **Connection Pooling**: Use PgBouncer in production for connection pooling.

3. **Read Replicas**: Configure read replicas for heavy read workloads (dashboard, charts).

4. **Vacuum**: Regular VACUUM ANALYZE on high-churn tables.

5. **Indexes**: Monitor index usage with `pg_stat_user_indexes`.

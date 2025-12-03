-- WMS Database Initialization Script
-- Created for WMS (Warehouse Management System)

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set timezone
SET timezone = 'Asia/Shanghai';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE wms_db TO wms_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO wms_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO wms_user;

-- Create schema version table for migrations tracking
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Log initialization
INSERT INTO schema_migrations (version)
VALUES ('initial_setup_' || to_char(CURRENT_TIMESTAMP, 'YYYYMMDD_HH24MISS'))
ON CONFLICT DO NOTHING;

-- Success message
DO $$
BEGIN
    RAISE NOTICE 'WMS database initialized successfully';
END $$;

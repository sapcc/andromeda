-- Database initialization script for Andromeda Akamai Reports testing

-- Ensure we're using the correct database
USE andromeda_db;

-- Note: The datacenter table will be created by the migration tool
-- We'll just insert test data after the migration runs

-- Insert test data for Akamai properties
-- These are example Akamai property names that would be retrieved from the Akamai GTM API
-- The collector will look for these in the database to map to project_ids for metrics
INSERT INTO datacenter (id, project_id, provider, admin_state_up) VALUES
    ('andromeda-mena.ccee.sapcloud.io', 'test-project-001', 'akamai', 1),
    ('andromeda-eu.ccee.sapcloud.io', 'test-project-002', 'akamai', 1),
    ('andromeda-us.ccee.sapcloud.io', 'test-project-003', 'akamai', 1)
ON DUPLICATE KEY UPDATE 
    project_id = VALUES(project_id),
    admin_state_up = VALUES(admin_state_up);

-- Verify the data was inserted
SELECT * FROM datacenter WHERE provider = 'akamai';
#!/bin/bash

# Script to insert test data after migrations have run

echo "Waiting for database tables to be ready..."
sleep 2

echo "Inserting test data into datacenter table..."
docker exec andromeda-mariadb mariadb -uandromeda -pandromeda andromeda_db -e "
INSERT INTO datacenter (id, project_id, provider, admin_state_up) VALUES
    ('fa471c2a-4b6d-11f0-a028-26af4eb7f05b', 'test-project-001', 'akamai', 1),
    ('a10d9971-4a8a-11f0-a028-26af4eb7f05b', 'test-project-002', 'akamai', 1),
    ('c6176476-1509-11f0-bc67-3eda432567b3', 'test-project-003', 'akamai', 1)
ON DUPLICATE KEY UPDATE 
    project_id = VALUES(project_id),
    admin_state_up = VALUES(admin_state_up);
"

echo "Verifying inserted data..."
docker exec andromeda-mariadb mariadb -uandromeda -pandromeda andromeda_db -e "SELECT id, project_id, provider FROM datacenter WHERE provider = 'akamai';"

echo "Test data inserted successfully!"
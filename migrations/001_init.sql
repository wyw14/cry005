CREATE TABLE IF NOT EXISTS workflow_items (id VARCHAR(64) PRIMARY KEY, scope VARCHAR(128) NOT NULL, state VARCHAR(32) NOT NULL, version INTEGER NOT NULL);
CREATE INDEX IF NOT EXISTS idx_workflow_items_scope ON workflow_items(scope);

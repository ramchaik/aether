CREATE TABLE logs (
    logid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    projectId UUID NOT NULL,
    log TEXT NOT NULL,
    timestamp BIGINT NOT NULL
);

CREATE INDEX idx_logs_projectid ON logs(projectId);
CREATE INDEX idx_logs_timestamp ON logs(timestamp);
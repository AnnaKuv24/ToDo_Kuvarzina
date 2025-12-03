ALTER TABLE tasks 
ADD COLUMN priority VARCHAR(10);

UPDATE tasks SET priority = 'MEDIUM' WHERE priority IS NULL;

ALTER TABLE tasks 
ALTER COLUMN priority SET DEFAULT 'MEDIUM';

ALTER TABLE tasks 
ADD CONSTRAINT tasks_priority_check 
CHECK (priority IN ('LOW', 'MEDIUM', 'HIGH'));
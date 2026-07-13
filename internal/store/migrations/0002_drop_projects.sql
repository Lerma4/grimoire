-- 0002_drop_projects.sql: remove the Projects feature entirely.
-- Drops the dead indexes and the projects table, then drops the now-orphan
-- project_id columns from tasks and notes. The columns must be dropped AFTER
-- the projects table so their REFERENCES clause no longer resolves.

DROP INDEX IF EXISTS idx_tasks_project_id;
DROP INDEX IF EXISTS idx_notes_project_id;
DROP TABLE projects;
ALTER TABLE tasks DROP COLUMN project_id;
ALTER TABLE notes DROP COLUMN project_id;

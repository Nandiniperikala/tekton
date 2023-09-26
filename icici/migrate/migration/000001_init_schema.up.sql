CREATE TABLE IF NOT EXISTS "event" (
  "id" integer generated always as identity,
  "title" varchar,
  "description" varchar
);

INSERT INTO event (title, description) VALUES
  ('Event 1', 'Description of Event 1'),
  ('Event 2', 'Description of Event 2'),
  ('Event 3', 'Description of Event 3');
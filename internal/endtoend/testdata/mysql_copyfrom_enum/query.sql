-- name: UpsertExperienceLocations :copyfrom
INSERT INTO experience_locations (location_id, type)
VALUES (?, ?);

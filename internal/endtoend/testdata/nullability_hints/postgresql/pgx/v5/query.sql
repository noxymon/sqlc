-- name: ListAuthors :many
SELECT 
    id,
    sqlc.nullable(name) as name_nullable,
    sqlc.notnull(bio) as bio_notnull,
    sqlc.nullable(count(*) over()) as total_count
FROM authors;

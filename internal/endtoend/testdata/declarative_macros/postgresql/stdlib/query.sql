/*
name: ListAuthorsByIDs :many
slice: ids
*/
SELECT * FROM authors WHERE id IN (@ids);

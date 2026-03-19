-- name: GetMainSequence :one
SELECT ID FROM Sequence
WHERE SeriesID = ?
ORDER BY Name = ? DESC, ID
LIMIT 1;

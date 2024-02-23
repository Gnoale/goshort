-- name: InsertURL :exec
INSERT INTO urls (url) VALUES (?);

-- name: GetURL :one
SELECT id FROM urls WHERE url = ?;


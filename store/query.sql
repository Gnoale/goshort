-- name: InsertURL :one
INSERT INTO urls (url) VALUES (?)
RETURNING id;

-- name: GetURL :one
SELECT url FROM urls WHERE id = ?;


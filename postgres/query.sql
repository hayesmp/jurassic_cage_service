-- name: GetDinosaur :one
SELECT * FROM dinosaur WHERE id = $1;

-- name: GetDinosaurByName :one
SELECT * FROM dinosaur WHERE name = $1;

-- name: GetDinosaurs :many
SELECT * FROM dinosaur;

-- name: GetDinosaursByCage :many
SELECT * FROM dinosaur WHERE cage_id = $1;

-- name: UpsertDinosaur :exec
INSERT INTO dinosaur (name, eating_habit, species, cage_id)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (name) DO UPDATE SET (eating_habit, species, cage_id) = ($2, $3, $4);

-- name: DeleteDinosaur :exec
DELETE FROM dinosaur WHERE id = $1;

-- name: GetCage :one
SELECT * FROM cage WHERE id = $1;

-- name: GetCages :many
SELECT * FROM cage;

-- name: UpsertCage :exec
INSERT INTO cage (name, status) VALUES ($1, $2)
    ON CONFLICT (name) DO UPDATE SET (status) = ($2);

-- name: UpdateCageStatus :one
UPDATE cage SET status = $1 WHERE id = $2 RETURNING *;

-- name: UpdateCageStatusByName :one
UPDATE cage SET status = $1 WHERE name = $2 RETURNING *;
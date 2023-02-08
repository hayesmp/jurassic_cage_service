-- name: GetDinosaur :one
SELECT * FROM dinosaur WHERE id = $1;

-- name: GetDinosaurByName :one
SELECT * FROM dinosaur WHERE name = $1;

-- name: GetDinosaurs :many
SELECT * FROM dinosaur;

-- name: GetDinosaursBySpecies :many
SELECT * FROM dinosaur WHERE species = $1;

-- name: GetCageDinosaurCount :one
SELECT COUNT(*) FROM dinosaur WHERE cage_id = $1;


-- name: GetDinosaurAndCage :one
SELECT dinosaur.*, cage.*
FROM dinosaur
    LEFT JOIN cage ON dinosaur.cage_id = cage.id
    WHERE dinosaur.id = $1;

-- name: GetDinosaursByCage :many
SELECT * FROM dinosaur WHERE cage_id = $1;

-- name: UpsertDinosaur :one
INSERT INTO dinosaur (name, eating_habit, species, cage_id)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (name) DO UPDATE SET (eating_habit, species, cage_id) = ($2, $3, $4)
    RETURNING *;

-- name: UpdateDinosaurCage :one
UPDATE dinosaur SET cage_id = $1 WHERE id = $2 RETURNING *;

-- name: DeleteDinosaur :exec
DELETE FROM dinosaur WHERE id = $1;

-- name: GetCage :one
SELECT * FROM cage WHERE id = $1;

-- name: GetCages :many
SELECT * FROM cage;

-- name: GetCagesByStatus :many
SELECT * FROM cage WHERE status = $1;

-- name: GetCageAndDinosaurs :one
SELECT cage.*, dinosaur.*
FROM cage
    LEFT JOIN dinosaur ON cage.id = dinosaur.cage_id
    WHERE cage.id = $1;

-- name: UpsertCage :one
INSERT INTO cage (name, status, predominate_eating_habit) VALUES ($1, $2, $3)
    ON CONFLICT (name) DO UPDATE SET (status, predominate_eating_habit) = ($2, $3)
    RETURNING *;

-- name: UpdateCagePredominateEatingHabit :exec
UPDATE cage SET predominate_eating_habit = $1 WHERE id = $2;

-- name: UpdateCageStatus :one
UPDATE cage SET status = $1 WHERE id = $2 RETURNING *;

-- name: UpdateCageStatusByName :one
UPDATE cage SET status = $1 WHERE name = $2 RETURNING *;

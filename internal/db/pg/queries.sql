-- name: GetCarrier :one
SELECT * FROM carriers WHERE id = $1;

-- name: CreateCarrier :one
INSERT INTO carriers (driver_category) VALUES ($1) RETURNING id;

-- name: UpdateCarrier :one
UPDATE carriers SET driver_category = $2 WHERE id = $1 RETURNING id;

-- name: DeleteCarrier :exec
DELETE FROM carriers WHERE id = $1;

-- name: GetCar :one
SELECT 
    id, 
    type, 
    length, 
    width, 
    height, 
    max_weight, 
    number, 
    owner 
FROM cars 
WHERE id = $1;

-- name: CreateCar :one
INSERT INTO cars (
    type, length, width, height, max_weight, number, owner
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING id;

-- name: UpdateCar :one
UPDATE cars 
SET 
    type = $2, 
    length = $3, 
    width = $4, 
    height = $5, 
    max_weight = $6, 
    number = $7, 
    owner = $8
WHERE id = $1 
RETURNING id;

-- name: DeleteCar :exec
DELETE FROM cars WHERE id = $1;

-- name: ListCarsByOwner :many
SELECT 
    id, type, length, width, height, max_weight, number, owner
FROM cars 
WHERE owner = $1 
ORDER BY id;
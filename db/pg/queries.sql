-- name: GetCarrier :one
SELECT * FROM carriers WHERE id = $1;

-- name: CreateCarrier :one
INSERT INTO carriers (id, driver_category) VALUES ($1, $2) RETURNING id;

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

-- name: GetRecipient :one
SELECT id, first_name, second_name, third_name, phone, email
FROM recipients
WHERE id = $1;

-- name: CreateRecipient :one
INSERT INTO recipients (first_name, second_name, third_name, phone, email)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdateRecipient :one
UPDATE recipients
SET first_name = $2,
    second_name = $3,
    third_name = $4,
    phone = $5,
    email = $6
WHERE id = $1
RETURNING id;

-- name: DeleteRecipient :exec
DELETE FROM recipients WHERE id = $1;

-- name: ListRecipients :many
SELECT id, first_name, second_name, third_name, phone, email
FROM recipients
ORDER BY id;

-- name: GetCargoRequest :many
select * from cargo_requests where id = $1;

-- name: InsertCargoRequest :one
insert into cargo_requests (id, consigner_id, recipient_id, from_station, to_station, created_at, deadline, route_id, trip_id, price, status, receive_code)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) returning id;

-- name: SelectStation :one
select id, address, lat, lon from stations s where s.id = $1;

-- name: SelectStations :many
select id, address, lat, lon from stations s where s.id in ($1, $2);

-- name: InsertStation :one
insert into stations (id, address, lat, lon) values ($1, $2, $3, $4) returning id;

-- name: GetCargoTypes :many
SELECT id, type, fragile FROM cargo_types;

-- name: CreateCargo :one
INSERT INTO cargo (length, width, height, weight, cargo_type, worth, request_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: GetTripByCargoRequest :one
SELECT t.*
FROM trips t
JOIN cargo_requests c ON c.trip_id = t.id
WHERE c.id = $1;

-- name: UpdateCargoRequestTrip :one
UPDATE cargo_requests 
SET trip_id = $2
WHERE id = $1 
RETURNING id;

-- name: GetSuitableTripsForCargoRequest :many
SELECT 
    t.route_id
FROM cargo_requests cr
JOIN cargo ON cargo.request_id = cr.id
JOIN cargo_types ct ON cargo.cargo_type = ct.id
CROSS JOIN trips t
JOIN cars c ON t.car = c.id
WHERE cr.id = $1
  AND cr.deadline >= t.calculate_end_at
  AND cargo.weight <= c.max_weight
  AND cargo.length <= c.length
  AND cargo.width <= c.width
  AND cargo.height <= c.height;

-- name: GetTripsByIDsWithPagination :many
SELECT 
    t.id,
    t.started_at,
    t.calculate_end_at,
    t.actual_end_at,
    t.price,
    t.status,
    t.carrier as carrier_id,
    t.car as car_id,
    fs.address as from_address,
    fs.lat as from_lat,
    fs.lon as from_lon,
    ts.address as to_address,
    ts.lat as to_lat,
    ts.lon as to_lon
FROM trips t
JOIN stations fs ON t.from_station = fs.id
JOIN stations ts ON t.to_station = ts.id
WHERE t.id = ANY($1::uuid[])
LIMIT $2 OFFSET $3;

-- name: InsertConsigner :exec
INSERT INTO consigners (id) VALUES ($1);
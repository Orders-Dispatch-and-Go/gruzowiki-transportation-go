-- name: GetCarrier :one
select * from carriers where id = $1;

-- name: CreateCarrier :one
insert into carriers (driver_category) values ($1) returning id;

-- name: GetCargoRequest :one
select * from cargo_requests where id = $1;

-- name: InsertCargoRequest :one
insert into cargo_requests (id, consigner_id, recipient_id, created_at, deadline, route_id, trip_id, price, status)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id;

-- name: GetCargoTypes :many
SELECT id, type, fragile FROM cargo_types;

-- name: CreateCargo :one
INSERT INTO cargo (length, width, height, weight, cargo_type, worth, request_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;
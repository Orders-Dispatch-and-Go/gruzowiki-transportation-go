-- name: GetCarrier :one
select * from carriers where id = $1;

-- name: CreateCarrier :one
insert into carriers (driver_category) values ($1) returning id;

-- name: GetCargoRequest :many
select * from cargo_requests where id = $1;

-- name: InsertCargoRequest :one
insert into cargo_requests (id, consigner_id, recipient_id, created_at, deadline, route_id, trip_id, price, status)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id;
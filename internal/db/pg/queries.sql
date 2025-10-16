-- name: GetCarrier :one
select * from carriers where id = $1;

-- name: CreateCarrier :one
insert into carriers (driver_category) values ($1) returning id;
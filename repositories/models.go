package repositories

type Carrier struct {
    Id             int    `db:"id"`
    DriverCategory string `db:"driver_category"`
}
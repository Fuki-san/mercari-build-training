package main

const (
	ImgDir = "images"
	JsonFile = "items.json"
	Sqlpath = "db/mercari.sqlite3"
)

type Response struct {
	Message string `json:"message"`
}

type Item struct {
	ID int `db:"id"`
	Name string `db:"name"`
	ImageName string `db:"image_name"`
	CategoryID int `db:"category_id"`
	CategoryName string `db:"category_name"`
}

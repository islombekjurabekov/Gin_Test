package domain

import (
	"database/sql"
)

type DBHouse struct {
	ID               string         `db:"id"`
	Address          string         `db:"address"`
	Number           int            `db:"number"`
	RoomName         string         `db:"room_name"`
	ColorOfBookshelf string         `db:"color_of_bookshelf"`
	Image            sql.NullString `db:"image"`
	ImageName        sql.NullString `db:"image_name"`
}

type HouseRest struct {
	ID               string `json:"id"`
	Address          string `json:"address"`
	Number           int    `json:"number"`
	RoomName         string `json:"room_name"`
	ColorOfBookshelf string `json:"color_of_bookshelf"`
	Image            string `json:"image"`
	ImageName        string `json:"image_name"`
}

type DBHome struct {
	Address          string         `db:"address"`
	Number           int            `db:"number"`
	RoomName         string         `db:"room_name"`
	ColorOfBookshelf string         `db:"color_of_bookshelf"`
	Image            sql.NullString `db:"image"`
	ImageName        sql.NullString `db:"image_name"`
}

type HomeRest struct {
	Address          string `json:"address"`
	Number           int    `json:"number"`
	RoomName         string `json:"room_name"`
	ColorOfBookshelf string `json:"color_of_bookshelf"`
	ImageKey         string `json:"image_key"`
	Image            string `json:"image"`
	ImageName        string `json:"image_name"`
}

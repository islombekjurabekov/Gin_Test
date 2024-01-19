package database

import (
	"Gin_test/domain"
	"bytes"
	"database/sql"
	"image"
	"image/png"
	"log"
	"os"
)

type DBStore struct {
	DB *sql.DB
}

func New(db *sql.DB) DBStore {
	return DBStore{
		DB: db,
	}
}

func DBHouseToModel(h []domain.DBHouse) []domain.HouseRest {
	v := make([]domain.HouseRest, len(h))
	for i := range h {
		v[i] = HouseToModel(h[i])
	}
	return v
}

func HouseToModel(h domain.DBHouse) domain.HouseRest {
	v := domain.HouseRest{
		ID:               h.ID,
		Address:          h.Address,
		Number:           h.Number,
		RoomName:         h.RoomName,
		ColorOfBookshelf: h.ColorOfBookshelf,
	}
	return v
}

func DBHomeRestToModel(h domain.HomeRest) domain.DBHome {
	view := domain.DBHome{
		Address:          h.Address,
		Number:           h.Number,
		RoomName:         h.RoomName,
		ColorOfBookshelf: h.ColorOfBookshelf,
		Image: sql.NullString{
			String: h.Image,
			Valid:  h.Image != "",
		},
		ImageName: sql.NullString{
			String: h.ImageName,
			Valid:  h.ImageName != "",
		},
	}
	return view
}

func (d *DBStore) GetHouseByID(id string) (domain.HouseRest, error) {
	rows := d.DB.QueryRow("SELECT id, address, number, room_name, color_of_bookshelf FROM my_house WHERE id = $1", id)
	var s domain.DBHouse
	err := rows.Scan(&s.ID, &s.Address, &s.Number, &s.RoomName, &s.ColorOfBookshelf)
	if err != nil {
		return domain.HouseRest{}, err
	}
	return HouseToModel(s), nil
}

func (d *DBStore) GetAllHouses() ([]domain.HouseRest, error) {
	var res []domain.DBHouse

	rows, err := d.DB.Query("SELECT id, address, number, room_name, color_of_bookshelf FROM my_house")
	if err != nil {
		log.Fatal(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var r domain.DBHouse
		err := rows.Scan(&r.ID, &r.Address, &r.Number, &r.RoomName, &r.ColorOfBookshelf)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return DBHouseToModel(res), nil
}

func (d *DBStore) AddNewHouse(param domain.HomeRest) {
	DBHomeRestToModel(param)
	_, err := d.DB.Exec("INSERT INTO my_house(address, number, room_name, color_of_bookshelf, image, image_name) VALUES ($1, $2, $3, $4, $5, $6)", param.Address, param.Number, param.RoomName, param.ColorOfBookshelf, param.Image, param.ImageName)
	if err != nil {
		log.Fatal(err)
	}

}

//func (d *DBStore) UpdateNewImage(buffer string, id string, res string) {
//	_, err := d.DB.Exec("UPDATE my_house SET image = $1, image_name = $2 WHERE id = $3", buffer.String(), res.Filename, id)
//	if err != nil {
//		log.Fatal(err)
//	}
//}

func (d *DBStore) GetImageByID(id string) (string, error) {
	var (
		images, imagesName string
		err                error
		img                image.Image
	)
	row := d.DB.QueryRow("SELECT image, image_name FROM my_house WHERE id = $1", id)
	err = row.Scan(&images, &imagesName)
	if err != nil {
		return "", err
	}
	bite := []byte(images)

	img, _, err = image.Decode(bytes.NewReader(bite))
	if err != nil {
		return "", err
	}

	out, _ := os.Create("./picture/" + imagesName)
	defer func(out *os.File) {
		err = out.Close()
		if err != nil {
			return
		}
	}(out)

	err = png.Encode(out, img)
	if err != nil {
		return "", err
	}
	return imagesName, nil
}

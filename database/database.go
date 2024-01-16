package database

import (
	"Gin_test/ConnectDatabase"
	"Gin_test/domain"
	"bytes"
	"database/sql"
	"image"
	"image/png"
	"log"
	"mime/multipart"
	"os"
)

func DBHouseToModel(h []domain.DBHouse) []domain.HouseRest {
	v := make([]domain.HouseRest, len(h))
	for i := range h {
		v[i] = domain.HouseRest{
			ID:               h[i].ID,
			Address:          h[i].Address,
			Number:           h[i].Number,
			RoomName:         h[i].RoomName,
			ColorOfBookshelf: h[i].ColorOfBookshelf,
			Image:            h[i].Image,
			ImageName:        h[i].ImageName,
		}
	}
	return v
}

func DBHomeRestToModel(h domain.HomeRest) domain.DBHome {
	view := domain.DBHome{
		Address:          h.Address,
		Number:           h.Number,
		RoomName:         h.RoomName,
		ColorOfBookshelf: h.ColorOfBookshelf,
		Image:            h.Image,
		ImageName:        h.ImageName,
	}
	return view
}

func GetHouseByID(id string) ([]domain.HouseRest, error) {
	var res []domain.DBHouse
	rows, err := ConnectDatabase.DB.Query("SELECT id, address, number, room_name, color_of_bookshelf, image_name FROM my_house WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	for rows.Next() {
		var s domain.DBHouse
		err := rows.Scan(&s.ID, &s.Address, &s.Number, &s.RoomName, &s.ColorOfBookshelf, &s.ImageName)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return DBHouseToModel(res), nil
}

func GetAllHouses() ([]domain.HouseRest, error) {
	var res []domain.DBHouse

	rows, err := ConnectDatabase.DB.Query("SELECT id, address, number, room_name, color_of_bookshelf, image_name FROM my_house")
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
		err := rows.Scan(&r.ID, &r.Address, &r.Number, &r.RoomName, &r.ColorOfBookshelf, &r.ImageName)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return DBHouseToModel(res), nil
}

func AddNewHouse(param domain.HomeRest) {
	DBHomeRestToModel(param)
	_, err := ConnectDatabase.DB.Exec("INSERT INTO my_house(address, number, room_name, color_of_bookshelf) VALUES ($1, $2, $3, $4)", param.Address, param.Number, param.RoomName, param.ColorOfBookshelf)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateNewImage(buffer bytes.Buffer, id string, res *multipart.FileHeader) {
	_, err := ConnectDatabase.DB.Exec("UPDATE my_house SET image = $1, image_name = $2 WHERE id = $3", buffer.String(), res.Filename, id)
	if err != nil {
		log.Fatal(err)
	}
}

func GetImageByID(id string) error {
	var (
		images, imagesName string
		err                error
		img                image.Image
	)
	row := ConnectDatabase.DB.QueryRow("SELECT image, image_name FROM my_house WHERE id = $1", id)
	err = row.Scan(&images, &imagesName)
	if err != nil {
		return err
	}
	bite := []byte(images)

	img, _, err = image.Decode(bytes.NewReader(bite))
	if err != nil {
		return err
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
		return err
	}
	return nil
}

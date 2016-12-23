package models

import (
	"database/sql"
	"encoding/json"
	"github.com/iot-spark/sparkmanager/db/sqlite"
	"github.com/iot-spark/sparkmanager/utils"
	"io"
	"log"
	"strconv"
)

type (
	DeviceIdentity struct {
		Id int `json:"id,omitempty"`
	}

	DeviceInfo struct {
		Name       string `json:"name"`
		Identity   string `json:"identity"`
		Password   string `json:"token"`
		CreateDate string `json:"created_date,omitempty"`
	}

	Device struct {
		DeviceIdentity
		DeviceInfo
	}

	DeviceModel struct {
		Name string `default:"iotspark.db"`
		Conn *sql.DB
	}
)

func NewDeviceModel(DbPath string) *DeviceModel {
	Conn := sqlite.InitDB(DbPath)

	dm := &DeviceModel{DbPath, Conn}
	dm.CreateTable()

	return dm
}

func (d *DeviceModel) CreateTable() bool {
	// create table if not exists
	sqlStmt := `
		CREATE TABLE IF NOT EXISTS devices(
			Id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
			Name char(32),
			PskIdentity char(32),
			PskPassword char(32),
			CreateDate DATETIME
		);`

	_, err := d.Conn.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

/**
 * Retrieves all items
 * @return - returns bytes array
 */
func (d *DeviceModel) GetAllDevices() ([]Device, bool) {
	var (
		devices []Device
		err     error
		rows    *sql.Rows
		status  = true
	)

	sqlStmt := `
		SELECT Id, Name, PskIdentity, PskPassword, CreateDate
		FROM devices
		ORDER BY datetime(CreateDate) DESC;`

	rows, err = d.Conn.Query(sqlStmt)
	if err != nil {
		log.Printf("GET ALL DEVICES: %#v\n", err)
		status = false
	}
	defer rows.Close()

	for rows.Next() {
		item := Device{}
		err = rows.Scan(&item.Id, &item.Name,
			&item.Identity, &item.Password, &item.CreateDate)

		if err != nil {
			log.Printf("GET ALL DEVICES: %#v\n", err)
		}

		devices = append(devices, item)
	}

	return devices, status
}

func (d *DeviceModel) GetDeviceById(id int) (DeviceInfo, bool) {
	var device DeviceInfo

	sqlStmt := `
		SELECT Name, PskIdentity, PskPassword, CreateDate
		FROM devices
		WHERE Id = ?`

	stmt, err := d.Conn.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return device, false
	}

	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&device.Name, &device.Identity,
		&device.Password, &device.CreateDate)
	if err != nil {
		log.Println(err)
		return device, false
	}

	return device, true
}

func (d *DeviceModel) GetDeviceByIdentity(identity string) (DeviceInfo, bool) {
	var device DeviceInfo

	sqlStmt := `
		SELECT Name, PskIdentity, PskPassword, CreateDate
		FROM devices
		WHERE PskIdentity = ?`

	stmt, err := d.Conn.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return device, false
	}
	defer stmt.Close()

	err = stmt.QueryRow(identity).Scan(&device.Name, &device.Identity,
		&device.Password, &device.CreateDate)

	if err != nil {
		log.Println(err)
		return device, false
	}

	return device, true
}

func (d *DeviceModel) RemoveItem(id int) error {
	tx, err := d.Conn.Begin()
	if err != nil {
		return err
	}

	sqlStmt := `DELETE from devices WHERE Id=?`

	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Println(err)
		return err
	}

	tx.Commit()

	return nil
}

/**
 * Adds a new device
 * @param io.Reader params - json
 * @return ([]Device, bool)
 */
func (d *DeviceModel) AddDevice(params io.Reader) ([]Device, bool) {
	var (
		err     error
		devices []Device
		token   string
		device  = Device{}
	)

	// Populate params
	err = json.NewDecoder(params).Decode(&device)
	if err != nil {
		log.Printf("INCOMING DATA: Json Decode error: %#v\n", err)
		return devices, false
	}

	// Generate PSK identity
	identity := utils.GeneratePskIdentity(device.Name)

	// is empty name
	if identity == "" {
		log.Printf("Device name is not specified!\n")
		return devices, false
	}

	// already in DB
	if _, ok := d.GetDeviceByIdentity(identity); ok {
		log.Printf("Device %#v has already registered in database!\n", device.Name)
		return devices, false
	}

	// Generate PSK password
	token, err = utils.GeneratePskKey(32)
	if err != nil {
		log.Printf("Generate PSK key error!\n")
		return devices, false
	}

	device.Identity = identity
	device.Password = token

	err = d.AddItem(device)
	if err != nil {
		log.Printf("Could not add device into database: %#v\n", err)
		return devices, false
	}

	devices = append(devices, device)

	return devices, true
}

func (d *DeviceModel) AddItem(device Device) error {
	tx, err := d.Conn.Begin()
	if err != nil {
		return err
	}

	sqlStmt := `INSERT INTO devices(Name, PskIdentity, PskPassword, CreateDate)
				VALUES (?, ?, ?, CURRENT_TIMESTAMP)`

	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(device.Name, device.Identity, device.Password)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (d *DeviceModel) DeleteDevice(id string) bool {
	status := false
	if id, err := strconv.Atoi(id); err == nil {
		_, ok := d.GetDeviceById(id)
		if ok {
			if err := d.RemoveItem(id); err == nil {
				status = true
				log.Printf("DELETE ITEM: Item id=%#v has been deleted!", id)
			}
		}
	}

	return status
}

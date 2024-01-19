package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	number, err := res.LastInsertId()

	return int(number), err
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("SELECT client, status, address, created_at FROM parcel WHERE number = :number",
		sql.Named("number", number))
	p := Parcel{Number: number}
	err := row.Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)

	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel

	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client",
		sql.Named("client", client))
	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return res, err
		}

		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	p, err := s.Get(number)

	if err != nil {
		return err
	}

	if p.Status == ParcelStatusRegistered {
		_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
			sql.Named("address", address),
			sql.Named("number", number))
	}

	return err
}

func (s ParcelStore) Delete(number int) error {
	p, err := s.Get(number)

	if err != nil {
		return err
	}

	if p.Status == ParcelStatusRegistered {
		_, err = s.db.Exec("DELETE FROM parcel  WHERE number = :number",
			sql.Named("number", number))
	}

	return err
}

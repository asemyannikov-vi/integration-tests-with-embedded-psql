package methods

import (
	"database/sql"
	"log"
	"moon/models"

	"github.com/jmoiron/sqlx"
)

const (
	getAllRowsFromAstroCatalogueTable   = "SELECT * FROM astro_catalogue;"
	getRowFromAstroCatalogueTableByName = "SELECT * FROM astro_catalogue WHERE astro_catalogue.name ~* :name;"
	addPlanetToAstroCatalogueTable      = "INSERT INTO astro_catalogue (name) VALUES (:name) RETURNING id;"
)

func GetRowsFromAstroCatalogueTable(db *sqlx.DB) []models.AstroCatalogueTableRow {
	data := make([]models.AstroCatalogueTableRow, 0)
	if err := db.Select(&data, getAllRowsFromAstroCatalogueTable); err != nil {
		log.Fatalf("Failed to execute a query `%s`.", getAllRowsFromAstroCatalogueTable)
	}
	return data
}

func GetRowFromAstroCatalogueTableByName(db *sqlx.DB, name string) models.AstroCatalogueTableRow {
	m := map[string]interface{}{"name": name}
	rows, err := db.NamedQuery(getRowFromAstroCatalogueTableByName, m)
	if err != nil {
		log.Fatalf("Failed to execute a query `%s`.", getRowFromAstroCatalogueTableByName)
	}
	var data models.AstroCatalogueTableRow
	for rows.Next() {
		err = rows.Scan(&data.Id, &data.Name)
		if err != nil {
			log.Fatalf("Failed to scan a row `%v`", err)
		}
	}
	return data
}

func AddPlanetToAstroCatalogueTable(db *sqlx.DB, name string) bool {
	m := map[string]interface{}{"name": name}
	_, err := db.NamedQuery(addPlanetToAstroCatalogueTable, m)
	if err != nil {
		log.Printf("Failed to execute a query `%s`.\n", addPlanetToAstroCatalogueTable)
		return false
	}
	if err != sql.ErrNoRows {
		return true
	}
	return false
}

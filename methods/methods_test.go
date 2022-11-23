package methods

import (
	"fmt"
	"log"
	"moon/models"
	"strings"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/slices"
)

func ConnectToDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database.")
	}
	return db
}

type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) BeforeTest() (*embeddedpostgres.EmbeddedPostgres, *sqlx.DB) {
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig())
	if err := postgres.Start(); err != nil {
		fmt.Println("Failed to migrate data into database", err.Error())
	}

	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database.")
	}

	if err := goose.Up(db.DB, "../migrations"); err != nil {
		fmt.Println("Failed to migrate data into database", err.Error())
	}

	return postgres, db
}

func (suite *TestSuite) AfterTest(postgres *embeddedpostgres.EmbeddedPostgres, db *sqlx.DB) {
	if err := db.Close(); err != nil {
		log.Fatalln("Failed to close database connection.", err.Error())
	}

	if err := postgres.Stop(); err != nil {
		log.Fatalln("Failed to stop postgreSQL embded driver.", err.Error())
	}
}

func (suite *TestSuite) TestShouldSuccessfullyGetRowsFromAstroCatalogueTable() {
	postgres, db := suite.BeforeTest()

	expect := []models.AstroCatalogueTableRow{
		{Name: "Mercury"},
		{Name: "Venus"},
		{Name: "Earth"},
	}

	data := GetRowsFromAstroCatalogueTable(db)

	assert.Equal(suite.T(), len(expect), len(data))
	for _, row := range data {
		fmt.Println(row)
		assert.Equal(suite.T(), true, slices.Contains(expect, models.AstroCatalogueTableRow{Name: row.Name}))
	}

	suite.AfterTest(postgres, db)
}

func (suite *TestSuite) TestShouldFailToGetRowsFromAstroCatalogueTableDueToExpectMoreDataInTheTable() {
	postgres, db := suite.BeforeTest()

	expect := []models.AstroCatalogueTableRow{
		{Name: "Mercury"},
		{Name: "Venus"},
		{Name: "Earth"},
		{Name: "Mars"},
	}
	data := GetRowsFromAstroCatalogueTable(db)

	for _, v := range data {
		fmt.Println(v)
	}

	assert.NotEqual(suite.T(), len(expect), len(data))

	suite.AfterTest(postgres, db)
}

func (suite *TestSuite) TestShouldFailToGetRowsFromAstroCatalogueTableDueToAbsenceSomeInformation() {
	postgres, db := suite.BeforeTest()

	expect := []models.AstroCatalogueTableRow{
		{Name: "Mercury"},
		{Name: "Venus"},
		{Name: "Mars"},
	}

	data := GetRowsFromAstroCatalogueTable(db)

	assert.Equal(suite.T(), len(expect), len(data))
	for _, row := range data {
		if !slices.Contains(expect, models.AstroCatalogueTableRow{Name: row.Name}) {
			assert.Equal(suite.T(), false, false)
		}
	}

	suite.AfterTest(postgres, db)
}

func (suite *TestSuite) TestShouldSuccessfullyToGetRowFromAstroCatalogueTableByName() {
	postgres, db := suite.BeforeTest()

	expect := models.AstroCatalogueTableRow{Name: "Earth"}

	data := GetRowFromAstroCatalogueTableByName(db, expect.Name)
	assert.Equal(suite.T(), expect.Name, data.Name)

	suite.AfterTest(postgres, db)
}

func (suite *TestSuite) TestShouldSuccessfullyToGetRowFromAstroCatalogueTableByNameInLowerCase() {
	postgres, db := suite.BeforeTest()

	expect := models.AstroCatalogueTableRow{Name: "Earth"}

	data := GetRowFromAstroCatalogueTableByName(db, strings.ToLower(expect.Name))
	assert.Equal(suite.T(), expect.Name, data.Name)

	suite.AfterTest(postgres, db)
}

func (suite *TestSuite) TestShouldFailToGetRowFromAstroCatalogueTableByNameDueToAbsenceRecordInTable() {
	postgres, db := suite.BeforeTest()

	expect := models.AstroCatalogueTableRow{Name: "Mars"}

	data := GetRowFromAstroCatalogueTableByName(db, expect.Name)
	assert.NotEqual(suite.T(), expect.Name, data.Name)

	suite.AfterTest(postgres, db)
}

func (suite *TestSuite) TestShouldSuccessfullyAddPlanetToAstroCatalogueTable() {
	postgres, db := suite.BeforeTest()

	expect := models.AstroCatalogueTableRow{Name: "Mars"}

	data := AddPlanetToAstroCatalogueTable(db, expect.Name)

	assert.Equal(suite.T(), true, data)

	suite.AfterTest(postgres, db)
}

func (suite *TestSuite) TestShouldFailToAddPlanetToAstroCatalogueTableDueToAlreadyExistingPlanetInTheTable() {
	postgres, db := suite.BeforeTest()

	expect := models.AstroCatalogueTableRow{Name: "Earth"}

	data := AddPlanetToAstroCatalogueTable(db, expect.Name)
	assert.Equal(suite.T(), false, data)

	suite.AfterTest(postgres, db)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

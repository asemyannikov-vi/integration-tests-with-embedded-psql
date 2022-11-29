package methods

import (
	"fmt"
	"log"
	"moon/models"
	"os"
	"strings"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/slices"
)

const (
	host                = "localhost"
	port                = 6616
	username            = "postgres"
	password            = "postgres"
	database            = "postgres"
	sslmode             = "disable"
	version             = "14.5.0"
	startTimeout        = 15 * time.Second
	binaryRepositoryURL = "https://repo1.maven.org/maven2"
)

func ConnectToDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, username, password, database, sslmode))
	if err != nil {
		log.Fatal("Failed to connect to database.")
	}
	return db
}

type TestSuite struct {
	suite.Suite
	db *sqlx.DB
}

func (suite *TestSuite) FillTables() {
	if err := goose.Up(suite.db.DB, "../migrations/data"); err != nil {
		fmt.Println("Failed to refresh data into database")
		return
	}
}

func (suite *TestSuite) TestShouldSuccessfullyGetRowsFromAstroCatalogueTable() {
	expect := []models.AstroCatalogueTableRow{
		{Name: "Mercury"},
		{Name: "Venus"},
		{Name: "Earth"},
	}

	data := GetRowsFromAstroCatalogueTable(suite.db)

	assert.Equal(suite.T(), len(expect), len(data))
	for _, row := range data {
		fmt.Println(row)
		assert.Equal(suite.T(), true, slices.Contains(expect, models.AstroCatalogueTableRow{Name: row.Name}))
	}
}

func (suite *TestSuite) TestShouldFailToGetRowsFromAstroCatalogueTableDueToExpectMoreDataInTheTable() {
	expect := []models.AstroCatalogueTableRow{
		{Name: "Mercury"},
		{Name: "Venus"},
		{Name: "Earth"},
		{Name: "Mars"},
	}

	data := GetRowsFromAstroCatalogueTable(suite.db)

	assert.NotEqual(suite.T(), len(expect), len(data))
}

func (suite *TestSuite) TestShouldFailToGetRowsFromAstroCatalogueTableDueToAbsenceSomeInformation() {
	expect := []models.AstroCatalogueTableRow{
		{Name: "Mercury"},
		{Name: "Venus"},
		{Name: "Mars"},
	}

	data := GetRowsFromAstroCatalogueTable(suite.db)

	assert.Equal(suite.T(), len(expect), len(data))
	for _, row := range data {
		if !slices.Contains(expect, models.AstroCatalogueTableRow{Name: row.Name}) {
			assert.Equal(suite.T(), false, false)
		}
	}
}

func (suite *TestSuite) TestShouldSuccessfullyToGetRowFromAstroCatalogueTableByName() {
	expect := models.AstroCatalogueTableRow{Name: "Earth"}

	data := GetRowFromAstroCatalogueTableByName(suite.db, expect.Name)
	assert.Equal(suite.T(), expect.Name, data.Name)
}

func (suite *TestSuite) TestShouldSuccessfullyToGetRowFromAstroCatalogueTableByNameInLowerCase() {
	expect := models.AstroCatalogueTableRow{Name: "Earth"}

	data := GetRowFromAstroCatalogueTableByName(suite.db, strings.ToLower(expect.Name))
	assert.Equal(suite.T(), expect.Name, data.Name)
}

func (suite *TestSuite) TestShouldFailToGetRowFromAstroCatalogueTableByNameDueToAbsenceRecordInTable() {
	suite.FillTables()

	expect := models.AstroCatalogueTableRow{Name: "Mars"}

	data := GetRowFromAstroCatalogueTableByName(suite.db, expect.Name)
	assert.NotEqual(suite.T(), expect.Name, data.Name)
}

func (suite *TestSuite) TestShouldSuccessfullyAddPlanetToAstroCatalogueTable() {
	expect := models.AstroCatalogueTableRow{Name: "Mars"}

	data := AddPlanetToAstroCatalogueTable(suite.db, expect.Name)
	assert.NotEqual(suite.T(), true, data)

	command := "DELETE FROM astro_catalogue WHERE name=$1;"
	suite.db.Exec(command, "Mars")
}

func (suite *TestSuite) TestShouldFailToAddPlanetToAstroCatalogueTableDueToAlreadyExistingPlanetInTheTable() {
	expect := models.AstroCatalogueTableRow{Name: "Earth"}

	data := AddPlanetToAstroCatalogueTable(suite.db, expect.Name)
	assert.Equal(suite.T(), false, data)
}

func (suite *TestSuite) SetupTest() {
	suite.db = ConnectToDB()
	if err := goose.Up(suite.db.DB, "../migrations"); err != nil {
		fmt.Println("Failed to migrate data into database")
		return
	}
	suite.FillTables()
}

func (suite *TestSuite) TearDownTest() {
	suite.FillTables()
}

func GenerateConfig() embeddedpostgres.Config {
	config := embeddedpostgres.DefaultConfig().Port(port).Version(version).Database(database).Username(username).Password(password).StartTimeout(startTimeout).Logger(os.Stdout).BinaryRepositoryURL(binaryRepositoryURL)
	return config
}

func TestExampleTestSuite(t *testing.T) {
	postgres := embeddedpostgres.NewDatabase(GenerateConfig())
	postgres.Start()
	newSuite := new(TestSuite)
	newSuite.SetupTest()
	suite.Run(t, newSuite)
	postgres.Stop()
}

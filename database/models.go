package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"log"
)

type ContributorModel struct {
	gorm.Model

	Name           string
	Current_bounty int `gorm:"default:0"`
}

type ContributorRecordModel struct {
	gorm.Model

	Contributor_name string
	Maintainer_name  string
	Pullreq_url      string
	Points_allotted  int
}

// TODO Implement method to connect GORM based on connection
// String
// Return GORM instance to store on main struct

// Manager struct
type DBManager struct {
	db *gorm.DB
}

func (manager *DBManager) Init(connection_string string) error {

	log.Println("[DBMANAGER] Initializing Database")
	// Initialize The GORM DB interface
	db, err := gorm.Open(sqlite.Open(connection_string), &gorm.Config{})
	if err != nil {
		log.Println("[ERROR][DBMANAGER] Could not initialize Database ->", err)
		return err
	} else {
		manager.db = db
		log.Println("[DBMANAGER] Successfully Initialized Database")
	}

	log.Println("[DBMANAGER] Beginning Model Automigration")

	err = manager.db.AutoMigrate(&ContributorModel{})
	if err != nil {
		log.Println("[ERROR][DBMANAGER] Could not AutoMigrate ContributorModel ->", err)
		return err
	} else {
		log.Println("[DBMANAGER] Successfully AutoMigrated ContributorModel")
	}

	err = manager.db.AutoMigrate(&ContributorRecordModel{})
	if err != nil {
		log.Println("[ERROR][DBMANAGER] Could not AutoMigrate ContributorRecordModel ->", err)
		return err
	} else {
		log.Println("[DBMANAGER] Successfully AutoMigrated ContributorRecordModel")
	}

	return nil
}

func (manager *DBManager) AssignBounty(
	maintainer string,
	contributor string,
	pr_html_url string,
	bounty_points int,
) error {

	// TODO Handle for Re-assignment
	// Start a New Transaction to create this object

	log.Println("[DBMANAGER][BOUNTY] Beginning Transaction to Assign Bounty")
	// Create the dummy record for the contributor_model
	// contributor_model := ContributorModel{name: contributor}

	// Create the time-series record of this transaction
	log.Println("[DBMANAGER][BOUNTY] Creating Contributor Record Model")

	crm := ContributorRecordModel{
		Maintainer_name:  maintainer,
		Contributor_name: contributor,
		Pullreq_url:      pr_html_url,
		Points_allotted:  bounty_points,
	}

	// Create the user struct
	// contributor_temp_representation := ContributorModel{
	// 	Name:           contributor,
	// 	Current_bounty: bounty_points,
	// }

	log.Println("[DBMANAGER][BOUNTY] Creating Contributor Record Model -> ", crm)
	log.Println("[DBMANAGER][BOUNTY] Beginning Transaction -> ", crm)

	manager.db.Transaction(func(tx *gorm.DB) error {

		// Create the time-series record
		result := tx.Create(&crm)
		if result.Error != nil {

			// Edge Case - User record already exists in time-series data
			// In that case, update that

			log.Println("[ERROR][DBMANAGER][BOUNTY] Could Not Create ContributorRecordModel ->", result.Error)
			return result.Error
		} else {
			log.Println("[DBMANAGER][BOUNTY] Successfully Created Contributor Record Model")
		}

		// default case - assume the user does not exist

		/*
			// Test if the user exists by attempting to create the user as
			// a new record
			user_create_result := tx.Create(&contributor_temp_representation)

			if user_create_result.Error != nil {
				// Check for the case where the user already exists

				// if that's the case, update the bounty with the new points

				// Else, report the error -> We found somethin unexpected

			} else {
				// Set the Bounty values
				// No Error, you can use this newly created user
				return nil
			}
		*/

		// commit the transaction
		return nil
	})

	return nil
}

func (manager *DBManager) Get_all_records() ([]ContributorRecordModel, error) {

	// Declare the array of all records
	var records []ContributorRecordModel

	// Fetch from the database
	log.Println("[DBMANAGER|RECORDS] Fetching All Records")
	fetch_result := manager.db.Find(&records)
	if fetch_result.Error != nil {
		log.Println("[ERROR][DBMANAGER|RECORDS] Could not fetch all records ->", fetch_result.Error)
		return nil, fetch_result.Error
	} else {
		log.Println("[DBMANAGER|RECORDS] Successfully Fetched all records")
		return records, nil
	}

}

func (manager *DBManager) Get_leaderboard() ([]ContributorModel, error) {

	// Declare the array of all records
	var records []ContributorModel

	// Fetch from the database
	log.Println("[DBMANAGER|LEADERBOARD] Fetching All Records")
	fetch_result := manager.db.Find(&records)

	if fetch_result.Error != nil {
		log.Println("[ERROR][DBMANAGER|LEADERBOARD] Could not fetch all records ->", fetch_result.Error)
		return nil, fetch_result.Error
	} else {
		log.Println("[DBMANAGER|LEADERBOARD] Successfully Fetched all records")
		return records, nil
	}

}

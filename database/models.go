package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"log"
)

type ContributorModel struct {
	gorm.Model
	name           string
	current_bounty int
}

type ContributorRecordModel struct {
	gorm.Model
	contributor_name   string
	maintainer_name    string
	pullreq_url        string
	points_allotted    int
	timestamp_allotted string
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

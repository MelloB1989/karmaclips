package services

import (
	"karmaclips/database"
	"karmaclips/utils"
	"log"
)

func CreateService(service *database.AiServices) (*database.AiServices, error) {
	db, err := database.DBConn()
	if err != nil {
		log.Fatal("DB connection failed!")
		return nil, err
	}
	defer db.Close()

	service.Aid = utils.GenerateID()

	err = database.InsertStruct(db, "ai_services", service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func GetServiceById(id string) (*database.AiServices, error) {
	db, err := database.DBConn()
	if err != nil {
		log.Fatal("DB connection failed!")
		return nil, err
	}
	defer db.Close()

	var services []*database.AiServices
	query := "SELECT * FROM ai_services WHERE aid = $1"

	rows, err := db.Query(query)
	if err != nil {

		return nil, err
	}
	defer rows.Close()

	if err := database.ParseRows(rows, &services); err != nil {
		return nil, err
	}

	return services[0], nil
}

func GetServices() ([]*database.AiServices, error) {
	db, err := database.DBConn()
	if err != nil {
		log.Fatal("DB connection failed!")
		return nil, err
	}
	defer db.Close()

	var services []*database.AiServices
	query := "SELECT * FROM ai_services"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := database.ParseRows(rows, &services); err != nil {
		return nil, err
	}

	return services, nil
}

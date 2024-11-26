package generations

import (
	"karmaclips/database"
	"karmaclips/utils"
)

func CreateGeneration(g *database.Generation) (*database.Generation, error) {
	db, err := database.DBConn()
	g.Id = utils.GenerateID()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = database.InsertStruct(db, "generations", g)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func GetGenerationById(gid string) (*database.Generation, error) {
	db, err := database.DBConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT * FROM generations WHERE id = $1"
	rows, err := db.Query(query, gid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var generation []*database.Generation

	if err := database.ParseRows(rows, &generation); err != nil {
		return nil, err
	}

	return generation[0], nil
}

func GetGenerationsByUserId(user_id string) ([]*database.Generation, error) {
	db, err := database.DBConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var generations []*database.Generation

	query := "SELECT * FROM generations WHERE user_id = $1"
	rows, err := db.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ParseRows(rows, &generations)
	if err != nil {
		return nil, err
	}

	return generations, nil
}

func GetGenerationsByUserIdAndType(user_id string, mtype string) ([]*database.Generation, error) {
	db, err := database.DBConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var generations []*database.Generation

	query := "SELECT * FROM generations WHERE user_id = $1 AND type = $2"
	rows, err := db.Query(query, user_id, mtype)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = database.ParseRows(rows, &generations)
	if err != nil {
		return nil, err
	}

	return generations, nil
}

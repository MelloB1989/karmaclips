package generations

import (
	"encoding/json"
	"karmaclips/database"
	"karmaclips/utils"
	"time"
)

func CreateGeneration(g *database.Generation) (*database.Generation, error) {
	db, err := database.DBConn()
	g.Id = utils.GenerateID()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	s, err := json.Marshal(g.Meta)

	gg := &database.GenerationDB{
		Id:          g.Id,
		CreatedBy:   g.CreatedBy,
		CreditsUsed: g.CreditsUsed,
		Timestamp:   g.Timestamp,
		MediaUri:    g.MediaUri,
		Type:        g.Type,
		Meta:        string(s),
	}

	err = database.InsertStruct(db, "generations", gg)
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

	query := "SELECT * FROM generations WHERE created_by = $1"
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

	query := "SELECT * FROM generations WHERE created_by = $1 AND type = $2"
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

func GetGenerationsByUserIdAndDate(user_id string, date time.Time) ([]*database.Generation, error) {
	db, err := database.DBConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var generations []*database.Generation

	query := "SELECT * FROM generations WHERE created_by = $1 AND DATE(timestamp) = $2"
	rows, err := db.Query(query, user_id, date.Format("2006-01-02"))
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

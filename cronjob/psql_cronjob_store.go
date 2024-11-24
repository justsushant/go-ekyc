package cronjob

import (
	"database/sql"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/db"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/types"
)

type PsqlCrobJobStore struct {
	db *sql.DB
}

func NewPsqlCrobJobStore(dsn string) PsqlCrobJobStore {
	psqlClient := db.NewPsqlClient(dsn)
	return PsqlCrobJobStore{
		db: psqlClient,
	}
}

func (s PsqlCrobJobStore) GetReportData(date string) ([]*types.ClientReport, error) {
	var report []*types.ClientReport
	err := s.db.QueryRow("").Scan()
	if err != nil {
		return nil, err
	}

	return report, nil
}

// SELECT  c.id AS clientID
// 		,c.name AS clientName
// 		,p.name AS planName
// 		,count(f.id) AS totalFaceMatch
// 		,count(o.id) AS totalOcrMatch
// 		,sum(u.file_size_kb)/1000 AS totalImgStorageMB
// 		,(count(f.id) * p.per_call_cost) + (count(o.id) * p.per_call_cost) AS totalAPIUsageCost
// 		,(count(f.id) * p.upload_cost_per_mb) + (count(o.id) * p.upload_cost_per_mb) AS totalStorageCost
// FROM client c
// JOIN plan p ON c.plan_id = p.id
// LEFT JOIN face_match f ON f.client_id = c.id
// LEFT JOIN ocr o ON o.client_id = c.id
// LEFT JOIN upload u ON u.client_id = c.id
// GROUP BY c.id, c.name, p.name, p.per_call_cost, p.upload_cost_per_mb

// SELECT
//     c.id AS clientID,
//     c.name AS clientName,
//     p.name AS planName,
//     CAST(f.created_at AS DATE) AS reportDate, -- Extract the date
//     COUNT(DISTINCT f.id) AS totalFaceMatch,
//     COUNT(DISTINCT o.id) AS totalOcrMatch,
//     COALESCE(SUM(u.file_size_kb), 0) / 1000 AS totalImgStorageMB,
//     (COUNT(DISTINCT f.id) * p.per_call_cost) + (COUNT(DISTINCT o.id) * p.per_call_cost) AS totalAPIUsageCost,
//     (COALESCE(SUM(u.file_size_kb), 0) / 1000) * p.upload_cost_per_mb AS totalStorageCost
// FROM client c
// JOIN plan p ON c.plan_id = p.id
// LEFT JOIN face_match f ON f.client_id = c.id
// LEFT JOIN ocr o ON o.client_id = c.id
// LEFT JOIN upload u ON u.client_id = c.id
// GROUP BY
//     c.id, c.name, p.name, p.per_call_cost, p.upload_cost_per_mb, CAST(f.created_at AS DATE);

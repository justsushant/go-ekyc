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
	query := `
        SELECT  
			c.id AS clientID,
			c.name AS clientName,
			p.name AS planName,
			CAST(f.created_at AS DATE) AS reportDate, -- Extract the date
			COUNT(DISTINCT f.id) AS totalFaceMatch,
			COUNT(DISTINCT o.id) AS totalOcrMatch,
			COALESCE(SUM(u.file_size_kb), 0) / 1000 AS totalImgStorageMB,
			(COUNT(DISTINCT f.id) * p.per_call_cost) + (COUNT(DISTINCT o.id) * p.per_call_cost) AS totalAPIUsageCost,
			(COALESCE(SUM(u.file_size_kb), 0) / 1000) * p.upload_cost_per_mb AS totalStorageCost
		FROM client c
		JOIN plan p ON c.plan_id = p.id
		LEFT JOIN face_match f ON f.client_id = c.id AND CAST(f.created_at AS DATE) = $1
		LEFT JOIN ocr o ON o.client_id = c.id AND CAST(o.created_at AS DATE) = $1
		LEFT JOIN upload u ON u.client_id = c.id
		GROUP BY 
			c.id, c.name, p.name, p.per_call_cost, p.upload_cost_per_mb, CAST(f.created_at AS DATE);
    `
	rows, err := s.db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r types.ClientReport
		if err := rows.Scan(&r.ClientID, &r.Name, &r.Plan, &r.Date, &r.TotalFaceMatch, &r.TotalOcr, &r.TotalImgStorageMB, &r.TotalAPIUsageCost, &r.TotalStorageCost); err != nil {
			return nil, err
		}
		report = append(report, &r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return report, nil
}

func (s PsqlCrobJobStore) GetMonthlyReport(currentMonth, currentYear int) ([][]*types.ClientReportMonthly, error) {
	return nil, nil
}

// SELECT
// 	c.id AS clientID
// 	,c.name AS clientName
// 	,p.name AS planName
// 	,CAST(f.created_at AS DATE) AS reportDate
// 	,COUNT(DISTINCT f.id) AS totalFaceMatch
// 	,COUNT(DISTINCT o.id) AS totalOcrMatch
// 	,COALESCE(SUM(u.file_size_kb), 0) / 1000 AS totalImgStorageMB
// 	,(COUNT(DISTINCT f.id) * p.per_call_cost) + (COUNT(DISTINCT o.id) * p.per_call_cost) AS totalAPIUsageCost
// 	,(COALESCE(SUM(u.file_size_kb), 0) / 1000) * p.upload_cost_per_mb AS totalStorageCost
// FROM client c
// JOIN plan p ON c.plan_id = p.id
// LEFT JOIN face_match f ON f.client_id = c.id AND CAST(f.created_at AS DATE) = $1
// LEFT JOIN ocr o ON o.client_id = c.id AND CAST(o.created_at AS DATE) = $1
// LEFT JOIN upload u ON u.client_id = c.id
// GROUP BY
// 	c.id, c.name, p.name, p.per_call_cost, p.upload_cost_per_mb, CAST(f.created_at AS DATE)

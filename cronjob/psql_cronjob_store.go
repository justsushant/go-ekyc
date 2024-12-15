package cronjob

import (
	"database/sql"
	"fmt"

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
			CAST(f.created_at AS DATE) AS reportDate,
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
	// fetch all client IDs
	clientQuery := `
		SELECT id 
		FROM client;
	`

	clientRows, err := s.db.Query(clientQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch client IDs: %w", err)
	}
	defer clientRows.Close()

	var clientIDs []int
	for clientRows.Next() {
		var clientID int
		if err := clientRows.Scan(&clientID); err != nil {
			return nil, fmt.Errorf("failed to scan client ID: %w", err)
		}
		clientIDs = append(clientIDs, clientID)
	}
	if err := clientRows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	// holds the final report
	var finalReports [][]*types.ClientReportMonthly

	// query to fetch daily reports for a specific client
	reportQuery := `
		WITH daily_usage AS (
			SELECT
				DATE(fm.created_at) AS report_date,
				COUNT(DISTINCT fm.id) AS total_face_match,
				COUNT(DISTINCT o.id) AS total_ocr,
				COALESCE(SUM(u.file_size_kb) / 1024.0, 0) AS total_image_storage_mb,
				p.per_call_cost,
				p.upload_cost_per_mb
			FROM
				client c
			LEFT JOIN plan p ON c.plan_id = p.id
			LEFT JOIN face_match fm ON c.id = fm.client_id AND DATE_PART('month', fm.created_at) = $1 AND DATE_PART('year', fm.created_at) = $2
			LEFT JOIN ocr o ON c.id = o.client_id AND DATE_PART('month', o.created_at) = $1 AND DATE_PART('year', o.created_at) = $2
			LEFT JOIN upload u ON c.id = u.client_id AND DATE_PART('month', u.created_at) = $1 AND DATE_PART('year', u.created_at) = $2
			WHERE c.id = $3
			GROUP BY DATE(fm.created_at), p.per_call_cost, p.upload_cost_per_mb
		)
		SELECT
			report_date,
			SUM(total_face_match) AS total_face_match_for_day,
			SUM(total_ocr) AS total_ocr_for_day,
			SUM(total_image_storage_mb) AS total_image_storage_in_mb,
			SUM(total_face_match * per_call_cost + total_ocr * per_call_cost) AS api_usage_cost_usd,
			SUM(total_image_storage_mb * upload_cost_per_mb) AS storage_cost_usd
		FROM
			daily_usage
		GROUP BY report_date;
	`

	// iterate over client ids to fetch reports
	for _, clientID := range clientIDs {
		rows, err := s.db.Query(reportQuery, currentMonth, currentYear, clientID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch report for client %d: %w", clientID, err)
		}
		defer rows.Close()

		var clientReports []*types.ClientReportMonthly
		for rows.Next() {
			report := &types.ClientReportMonthly{
				ClientID: fmt.Sprintf("%d", clientID),
			}
			err := rows.Scan(
				&report.Date,
				&report.TotalFaceMatch,
				&report.TotalOcr,
				&report.TotalImgStorageMB,
				&report.TotalAPIUsageCost,
				&report.TotalStorageCost,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to scan report for client %d: %w", clientID, err)
			}
			clientReports = append(clientReports, report)
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("rows error for client %d: %w", clientID, err)
		}

		// append to final report
		finalReports = append(finalReports, clientReports)
	}

	return finalReports, nil
}

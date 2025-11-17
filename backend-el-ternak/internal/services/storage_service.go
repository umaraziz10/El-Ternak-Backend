package services

import (
	"backend-el-ternak/internal/models"
	"backend-el-ternak/internal/repository"
)

func GetCurrentStock() (*models.StorageResponse, error) {
	stocks, err := repository.GetCurrentStock()
	if err != nil {
		return nil, err
	}

	return stocks, nil
}

func CheckPakanStock() (*models.CheckPakanResponse, error) {
	data, err := repository.CheckPakanStock()
	if err != nil {
		return data, err
	}

	return data, nil
}

func GetYearlyReport(tahun string) (*models.StorageReport, error) {
	reports, err := repository.GetYearlyReport(tahun)
	if err != nil {
		return nil, err
	}

	return reports, nil
}
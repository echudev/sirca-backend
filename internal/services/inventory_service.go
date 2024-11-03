package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

func GenerateInventaryCode(reqTypeName string, reqBrandName string, reqModelName string, reqItemAdquisitonDate pgtype.Date, reqItemId int32) (string, error) {
	// Validar que los strings tengan al menos 3 caracteres
	if len(reqTypeName) < 3 {
		return "", fmt.Errorf("type name must be at least 3 characters long")
	}
	if len(reqBrandName) < 3 {
		return "", fmt.Errorf("brand name must be at least 3 characters long")
	}
	if len(reqModelName) < 3 {
		return "", fmt.Errorf("model name must be at least 3 characters long")
	}

	// Validar fecha
	if !reqItemAdquisitonDate.Valid {
		return "", fmt.Errorf("invalid date")
	}

	// Validar ID
	if reqItemId <= 0 {
		return "", fmt.Errorf("invalid item ID")
	}

	year_code := strconv.Itoa(reqItemAdquisitonDate.Time.Year())
	id_code := strconv.Itoa(int(reqItemId))

	item_code := strings.ToUpper(reqTypeName[:3] + "-" + reqBrandName[:3] + "-" + reqModelName[:3] + "-" + year_code + "-" + id_code)

	return item_code, nil
}

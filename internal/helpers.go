package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

var ErrNoParameter = errors.New("no parameter given")

func ParseIntParam(w http.ResponseWriter, r *http.Request, key string) (int, error) {
	val := chi.URLParam(r, key)
	return ParseInt(w, r, val, fmt.Sprintf("Invalid %v '%v' on request parameter. Must be an integer.", key, val))
}

func ParseIntQuery(w http.ResponseWriter, r *http.Request, key string) (int, error) {
	val := r.URL.Query().Get(key)
	if val != "" {
		return ParseInt(w, r, val, fmt.Sprintf("Invalid %v '%v' on query parameter. Must be an integer.", key, val))
	} else {
		return -1, ErrNoParameter
	}
}

func ParseInt(w http.ResponseWriter, r *http.Request, val string, errMessage string) (int, error) {
	num, err := strconv.ParseInt(val, 10, 0)
	if err != nil {
		http.Error(w, errMessage, http.StatusBadRequest)
		return -1, err
	} else {
		return int(num), nil
	}
}

func QueryName(w http.ResponseWriter, db *gorm.DB, name string) (query *gorm.DB, err error) {
	name = strings.ToLower(name)
	names := strings.Split(name, " ")
	if len(names) != 2 {
		http.Error(w, "Name must be of format 'First Last'.", http.StatusBadRequest)
		return query, errors.New("name must be of format 'First Last'")
	}
	first, last := names[0], names[1]
	query = db.Where("LOWER(first_name) = ? AND LOWER(last_name) = ?", first, last)
	return query, nil
}

func CheckJSON(w http.ResponseWriter, r *http.Request, v any) error {
	ct := r.Header.Get("Content-Type")
	mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
	switch mediaType {
	case "application/json":
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err := dec.Decode(&v)
		if err != nil {
			Out(err)
			http.Error(w, "Request body is not valid JSON", http.StatusBadRequest)
		}
		return err
	default:
		msg := "Content-Type header is not application/json"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return errors.New(msg)
	}
}

func HandleDBError(w http.ResponseWriter, err error) error {
	if err != nil {
		switch err {
		case gorm.ErrDuplicatedKey:
			http.Error(w, "JSON body conflicts with existing course data", http.StatusBadRequest)
		}
	}
	return err
}

func HandleDBErrorGeneric(w http.ResponseWriter, err error) error {
	if err != nil {
		Out(err)
		http.Error(w, "Internal SQL Exception", http.StatusInternalServerError)
	}
	return err
}

package internal

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetPersons(w http.ResponseWriter, r *http.Request) {
	query := DB

	age, err := ParseIntQuery(w, r, "age")
	if (err != nil) && (err != ErrNoParameter) {
		return
	} else if !errors.Is(err, ErrNoParameter) {
		query = query.Where("age = ?", age)
	}
	name := r.URL.Query().Get("name")
	if name != "" {
		query, err = QueryName(w, DB, name)
		if err != nil {
			return
		}
	}

	var persons []Person
	if err = LoadAllPersonCourses(query, &persons); err != nil {
		HandleDBErrorGeneric(w, err)
	}
	render.JSON(w, r, persons)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	query, err := QueryName(w, DB, name)
	if err != nil {
		return
	}
	var person Person
	if err = LoadPersonCourses(query, &person); err != nil {
		http.Error(w, fmt.Sprintf("Person with name '%v' not found.", name), http.StatusNotFound)
		return
	}
	render.JSON(w, r, person)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var newPerson Person
	if err = CheckJSON(w, r, &newPerson); err != nil {
		return
	}

	if err = DB.Create(&newPerson).Error; errors.Is(err, gorm.ErrDuplicatedKey) {
		http.Error(w, fmt.Sprintf("JSON id '%v' conflicts with existing person data.", newPerson.ID), http.StatusConflict)
		return
	} else if errors.Is(err, gorm.ErrCheckConstraintViolated) {
		http.Error(w, fmt.Sprintf("Invalid type '%v'. Type must be either 'student' or 'professor'.", newPerson.Type), http.StatusBadRequest)
		return
	} else if err != nil {
		HandleDBErrorGeneric(w, err)
		return
	}

	output := map[string]int{"id": newPerson.ID}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, output)

}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	var newPerson Person
	if err = CheckJSON(w, r, &newPerson); err != nil {
		return
	}

	name := chi.URLParam(r, "name")
	query, err := QueryName(w, DB, name)
	if err != nil {
		return
	}

	var person Person
	if err = LoadPersonCourses(query, &person); errors.Is(err, logger.ErrRecordNotFound) {
		http.Error(w, fmt.Sprintf("Course with name '%v' not found.", name), http.StatusNotFound)
		return
	} else if err != nil {
		HandleDBErrorGeneric(w, err)
		return
	} else {

		err = DB.Transaction(func(db *gorm.DB) error {
			if len(person.Courses) > 0 {
				if err = db.Model(&person).Association("Courses").Delete(person.Courses); err != nil {
					HandleDBErrorGeneric(w, err)
					return err
				}
			}

			newPerson.ID = person.ID
			if err := db.Updates(&newPerson).Error; errors.Is(err, gorm.ErrCheckConstraintViolated) {
				http.Error(w, fmt.Sprintf("Invalid type '%v'. Type must be either 'student' or 'professor'.", newPerson.Type), http.StatusBadRequest)
				return err
			} else if err != nil {
				HandleDBErrorGeneric(w, err)
				return err
			}

			newPerson, err = LoadPerson(db, newPerson)
			if err != nil {
				HandleDBErrorGeneric(w, err)
				return err
			}
			return nil
		})
		if err != nil {
			return
		}

	}

	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, newPerson)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	query, err := QueryName(w, DB, name)
	if err != nil {
		return
	}

	var person Person
	var msg string
	if err = LoadPersonCourses(query, &person); errors.Is(err, logger.ErrRecordNotFound) {
		msg = fmt.Sprintf("No person found with name '%v'", name)
		// render.Status(r, http.StatusNoContent)
	} else if err != nil {
		HandleDBErrorGeneric(w, err)
		return
	} else {
		err = DB.Transaction(func(db *gorm.DB) error {
			if len(person.Courses) > 0 {
				if err = db.Model(&person).Association("Courses").Delete(person.Courses); err != nil {
					HandleDBErrorGeneric(w, err)
					return err
				}
			}

			if err := db.Delete(&person).Error; err != nil {
				HandleDBErrorGeneric(w, err)
				return err
			}
			return nil
		})
		if err != nil {
			return
		}
		msg = "Deletion Successful."
	}

	output := map[string]string{"message": msg}
	render.JSON(w, r, output)
}

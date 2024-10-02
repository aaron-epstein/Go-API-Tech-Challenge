package internal

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetCourses(w http.ResponseWriter, r *http.Request) {
	var courses []Course
	DB.Find(&courses)
	render.JSON(w, r, courses)
}

func GetCourse(w http.ResponseWriter, r *http.Request) {
	id, err := ParseIntParam(w, r, "id")
	if err != nil {
		return
	}
	course := Course{ID: id}
	err = DB.First(&course).Error
	if err != nil {
		http.Error(w, fmt.Sprintf("Course with id '%v' not found.", id), http.StatusNotFound)
		return
	}
	render.JSON(w, r, course)
}

func CreateCourse(w http.ResponseWriter, r *http.Request) {
	var newCourse Course
	if err = CheckJSON(w, r, &newCourse); err != nil {
		return
	}

	err = DB.Create(&newCourse).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		http.Error(w, fmt.Sprintf("JSON id '%v' conflicts with existing course data.", newCourse.ID), http.StatusConflict)
		return
	} else if err != nil {
		HandleDBErrorGeneric(w, err)
		return
	}

	output := map[string]int{"id": newCourse.ID}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, output)
}

func UpdateCourse(w http.ResponseWriter, r *http.Request) {
	var newCourse Course
	if err = CheckJSON(w, r, &newCourse); err != nil {
		return
	}

	id, err := ParseIntParam(w, r, "id")
	if err != nil {
		return
	}

	course := Course{ID: id}
	if err = DB.First(&course).Error; err != nil {
		http.Error(w, fmt.Sprintf("Course with id '%v' not found.", id), http.StatusNotFound)
		return
	}

	newCourse.ID = id
	if err = DB.Updates(newCourse).Error; err != nil {
		HandleDBErrorGeneric(w, err)
		return
	}
	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, newCourse)
}

func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	id, err := ParseIntParam(w, r, "id")
	if err != nil {
		return
	}

	course := Course{ID: id}

	var msg string
	if course, err = LoadCourse(DB, &course); errors.Is(err, logger.ErrRecordNotFound) {
		msg = fmt.Sprintf("No course found with id '%v'", id)
		// render.Status(r, http.StatusNoContent)
	} else if err != nil {
		HandleDBErrorGeneric(w, err)
		return
	} else {
		if len(course.Persons) > 0 {
			if err = DB.Model(&course).Association("Persons").Delete(course.Persons); err != nil {
				HandleDBErrorGeneric(w, err)
				return
			}
		}

		if err := DB.Delete(&course).Error; err != nil {
			HandleDBErrorGeneric(w, err)
			return
		}

		msg = "Deletion Successful."
	}

	output := map[string]string{"message": msg}
	render.JSON(w, r, output)
}

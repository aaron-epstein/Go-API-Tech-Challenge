package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

/*
Person definitions.
*/
type Person struct {
	ID        int      `json:"id,omitempty" gorm:"column:id;primaryKey;autoIncrement"`
	FirstName string   `json:"first_name,omitempty" gorm:"column:first_name"`
	LastName  string   `json:"last_name,omitempty" gorm:"column:last_name"`
	Type      string   `json:"type,omitempty" gorm:"column:type;check:type IN ('professor', 'student')"`
	Age       int      `json:"age,omitempty" gorm:"column:age"`
	Courses   []Course `json:"courses,omitempty" gorm:"many2many:person_course"`
}

func (p *Person) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == `""` {
		return nil
	}

	var person PersonJSON
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&person); err != nil {
		return err
	}

	courses := make([]Course, len(person.Courses))
	for i, id := range person.Courses {
		courses[i] = Course{ID: id}
	}
	*p = Person{
		ID:        person.ID,
		FirstName: person.FirstName,
		LastName:  person.LastName,
		Type:      person.Type,
		Age:       person.Age,
		Courses:   courses,
	}
	return nil
}

type PersonJSON struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Type      string `json:"type,omitempty"`
	Age       int    `json:"age,omitempty"`
	Courses   []int  `json:"courses,omitempty"`
}

func (s Person) String() string {
	courses := ""
	for _, course := range s.Courses {
		courses += strings.ReplaceAll(course.String(), "\n", "\n    ") + ", "
	}
	return fmt.Sprintf(
		`{
  ID: %d
  FirstName: %v
  LastName: %v
  Type: %v
  Age: %d
  Courses: [
    %v
  ]
}`, s.ID, s.FirstName, s.LastName, s.Type, s.Age, courses)
}

func (Person) TableName() string {
	return "person"
}

func LoadAllPersons(db *gorm.DB) ([]Person, error) {
	var person []Person
	err := db.Model(&Person{}).Preload("Courses").Find(&person).Error
	return person, err
}

func LoadPerson(db *gorm.DB, query any) (Person, error) {
	var person Person
	err := db.Model(&Person{}).Where(query).Preload("Courses").First(&person).Error
	return person, err
}

func LoadAllPersonCourses(dbQuery *gorm.DB, persons *[]Person) error {
	return dbQuery.Preload("Courses").Find(&persons).Error
}

func LoadPersonCourses(dbQuery *gorm.DB, person *Person) error {
	return dbQuery.Preload("Courses").First(&person).Error
}

/*
Course definitions.
*/
type Course struct {
	ID      int      `json:"id,omitempty"   gorm:"column:id;primaryKey;autoIncrement"`
	Name    string   `json:"name,omitempty" gorm:"column:name"`
	Persons []Person `json:"-"              gorm:"many2many:person_course"`
}

func (s Course) String() string {
	persons := ""
	for _, person := range s.Persons {
		persons += strings.ReplaceAll(person.String(), "\n", "\n    ") + ", "
	}
	return fmt.Sprintf(
		`{
  ID: %d
  Name: %v
  (Users): [
    %v
  ]
}`, s.ID, s.Name, persons)
}

func (Course) TableName() string {
	return "course"
}

func LoadCourse(db *gorm.DB, query any) (Course, error) {
	var course Course
	err := db.Model(&Course{}).Where(query).Preload("Persons").First(&course).Error
	return course, err
}

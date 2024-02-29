package model

import "errors"

type Course struct {
	CourseId       int    `json:"course_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	CourseDuration string `json:"courseDuration"`
}

var courses = []Course{
	{
		CourseId:       1,
		Title:          "Programming Principles I",
		Description:    "Starter programming courses",
		CourseDuration: "3 month",
	}, {
		CourseId:       2,
		Title:          "Programming Principles II",
		Description:    "Advanced programming courses",
		CourseDuration: "4 month",
	}, {
		CourseId:       3,
		Title:          "Algorithms and Data Structures",
		Description:    "This course for student who accepted their fate",
		CourseDuration: "3 month",
	}, {
		CourseId:       4,
		Title:          "Computer Networks",
		Description:    "This course will provide good topic for understanding world web",
		CourseDuration: "5 month",
	},
}

func GetCourses() []Course {
	return courses
}

func getCoursesById(id int) (*Course, error) {
	for _, c := range courses {
		if c.CourseId == id {
			return &c, nil
		}
	}
	return nil, errors.New("Courses not Found")
}

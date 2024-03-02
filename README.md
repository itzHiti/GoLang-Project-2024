# Online course management

## Features of our project and small description

### Features:
- [x] REST API for Courses
- [x] REST API for Users (Students)
- [ ] API structure
- [x] Structure by Standart Layout
- [x] File with handlers according to the API description
- [x] Migrations to the main entities and relationship tables


### Description:
> The project is intended to be an online learning website of school of courses.
> Project name: WIP, but for now it's `Online Course Management`.
> Project entities: `courses, users`.

## Members

| Full Name | Student ID |
| --------- | ---------- |
| Abzalkhanuly Alan | 22B030505 |
| Karim Madi | 22B031181 |
| Baynazarov Ramadan | 22B030523 |
| Abdullaev Shakhzod |  |

## OCM's Rest API

```
POST /courses
GET /courses/:id
PUT /courses/:id
DELETE /courses/:id
```

## Database structure

```

Table courses {
    course_id int [primary key, unique, increment]     
    title varchar(255) [not null]
    description text    course_duration varchar(50)
    student_id int [ref: > students.student_id] 
    }

Table students {    
    student_id int [primary key, unique, increment] 
    name varchar(50) [not null]    
    age int [not null]
    gpa float
    }
Table courses_and_students {
    id int [primary key, unique, increment]    
    description text
    student_id int [ref: > students.student_id] // Ссылается на student_id в таблице students    
    student_name varchar(50) [not null]
    course_id int [ref: > courses.course_id] // Ссылается на course_id в таблице courses    
    course_title varchar(255) [not null]
    }

```
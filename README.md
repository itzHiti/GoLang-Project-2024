# Online course management

## Features of our project and small description

### Features:
- [x] [REST API for Courses](#ocms-rest-api) ([SS](#rest-api))
- [x] API structure
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
| Abdullaev Shakhzod | 22B031601 |

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
    course_duration varchar(50) 
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

## Screenshots

### REST API:

| REST API method | Screenshot | Description |
| ----- | ---------- | --------- |
| GET | ![image](https://github.com/itzHiti/GoLang-Project-2024/assets/81374715/3a288d93-2414-4403-9b40-ce471b02a521) | Get method returns query's row in JSON format if it exists in our database, if not it will return `"Courses not Found"` or `"404 page not found"` | 
| POST | ![image](https://github.com/itzHiti/GoLang-Project-2024/assets/81374715/db001542-a970-46c0-8236-52a6468e3933) ![image](https://github.com/itzHiti/GoLang-Project-2024/assets/81374715/d37184bb-6590-4521-a988-8b45958d7fa2) | Post method adds new data in database (2nd screenshot) which was inputed and posted by user and returning new course's id (1st screenshot) |
| PUT (UPDATE) | ![image](https://github.com/itzHiti/GoLang-Project-2024/assets/81374715/4d4121ae-54d0-4cea-9c22-cc7ef8cdfc6f) ![image](https://github.com/itzHiti/GoLang-Project-2024/assets/81374715/c2d81096-d9be-4da9-bb4d-37f5d17688d9) | Put method updates data in our database and returns its new value after updating data if it exists (2nd screenshot). In our example (1st screenshot) we changed `"physics"` to `"physics II"` and `"3 month"` to `"4 month"` |
| DELETE | ![image](https://github.com/itzHiti/GoLang-Project-2024/assets/81374715/c997d581-bdd9-4bb5-8b78-b99ccfae20b3) | Delete method delete's data from database and returns text in JSON format if it succeeded or not |
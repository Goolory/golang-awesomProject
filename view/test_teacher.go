package view

import (
	"awesomeProject/dbmodel"
	"github.com/jinzhu/gorm"
)

type TestTeacher struct {
	dbmodel.Test
	TeacherName string `json:"teacher_name""`
}

func (TestTeacher) TableName() string {
	return "test_teacher"
}

func (v TestTeacher) CreateView(db *gorm.DB) error {
	err := db.Exec("CREATE ALGORITHM = MERGE VIEW " + v.TableName() + " AS" +
		" SELECT test.*, th.teacher_name AS teacher_name" +
		" FROM test" +
		" LEFT JOIN teacher as th ON test.user_id = th.id").Error

	return err
}

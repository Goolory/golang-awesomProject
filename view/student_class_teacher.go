package view

import (
	"awesomeProject/dbmodel"
	"github.com/jinzhu/gorm"
)

type StudentClassTeacher struct {
	dbmodel.User
	ClassId     uint32 `json:"class_id"`
	ClassName   string `json:"class_name"`
	TeacherId   uint32 `json:"teacher_id"`
	TeacherNo   string `json:"teacher_no"`
	TeacherName string `json:"teacher_name"`
	TeacherType string `json:"teacher_type"`
}

func (StudentClassTeacher) TableName() string {
	return "student_class_teacher"
}

func (v StudentClassTeacher) CreateView(db *gorm.DB) error {
	err := db.Exec("CREATE ALGORITHM = MERGE VIEW " + v.TableName() + " AS" +
		" SELECT user.*, class_name, t.id AS teacher_id, teacher_no, teacher_name, t.type AS teacher_type" +
		" FROM user" +
		" LEFT JOIN class AS c ON user.class_id = c.id" +
		" LEFT JOIN teacher AS t ON c.teacher_id = t.id").Error

	return err
}

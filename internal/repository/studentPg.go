package repository

import (
	"context"
	"errors"

	"CRUD_Go_Backend/internal/handlers/models"
	"CRUD_Go_Backend/internal/pkg/connection"
	"CRUD_Go_Backend/internal/pkg/pkgErrors"
	"CRUD_Go_Backend/internal/repository/entities"

	"github.com/jackc/pgx"
)

type StudentStorage struct {
	db connection.DBops
}

func NewStudentStorage(database connection.DBops) StudentStorage {
	return StudentStorage{db: database}
}

func ToStudentStorage(s models.StudentRequest) entities.Student {
	return entities.Student{
		StudentID:   s.StudentID,
		StudentName: s.StudentName,
		Grade:       s.Grade,
	}
}

func (r *StudentStorage) Add(ctx context.Context, studentReq models.StudentRequest) (int64, error) {
	student := ToStudentStorage(studentReq)
	var studentID int64
	err := r.db.ExecQueryRow(ctx,
		`INSERT INTO student(student_name, grade) VALUES($1, $2) RETURNING student_id;`,
		student.StudentName,
		student.Grade,
	).Scan(&studentID)

	if err != nil {
		return -1, err
	}

	return studentID, err
}

func (r *StudentStorage) GetByID(ctx context.Context, studentID int64) (models.StudentRequest, error) {
	var student entities.Student

	err := r.db.Get(
		ctx,
		&student,
		`SELECT student_id, student_name, grade, created_at FROM student WHERE student_id=$1;`,
		studentID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.StudentRequest{}, pkgErrors.ErrNotFound
		}

		return models.StudentRequest{}, err
	}

	return student.ToStudentDomain(), nil
}

func (r *StudentStorage) Delete(ctx context.Context, studentID int64) error {
	command, err := r.db.Exec(ctx, "DELETE FROM student WHERE student_id = $1", studentID)
	if err != nil {
		return err
	}

	if command.RowsAffected() == 0 {
		return pkgErrors.ErrNotFound
	}

	return nil
}

func (r *StudentStorage) Update(ctx context.Context, studentID int64, studentReq models.StudentRequest) error {
	student := ToStudentStorage(studentReq)

	command, err := r.db.Exec(ctx, `
		UPDATE student
		SET student_name = $2, grade = $3
		WHERE student_id = $1
	`, studentID, student.StudentName, student.Grade)

	if err != nil {
		return err
	}

	if command.RowsAffected() == 0 {
		return pkgErrors.ErrNotFound
	}

	return nil
}

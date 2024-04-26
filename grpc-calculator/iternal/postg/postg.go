package postg

import (
	"context"
	"culc/iternal/lib/postfix"
	"culc/iternal/model"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrUserExist    = errors.New("User  exist")
	ErrUserNotFound = errors.New("User not exist")
)

type Storage struct {
	db *gorm.DB
}

func ConnectDB(user string, password string, host string, port string, dbname string) *Storage {
	connectString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, password, host, port, dbname)
	db, err := gorm.Open(postgres.Open(connectString), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	if _, err := db.DB(); err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.User{}, &model.Expression{})
	return &Storage{
		db: db,
	}
}
func (s *Storage) SaveUser(ctx context.Context, email string, password string) (uint, error) {
	NewUser := model.User{Email: email, Password: password}
	if res := s.db.WithContext(ctx).Create(&NewUser); res.Error != nil {
		return 0, fmt.Errorf("error in registry")
	}
	return NewUser.ID, nil
}
func (s *Storage) UserProvider(ctx context.Context, email string) (model.User, error) {
	var GetUser model.User
	if res := s.db.WithContext(ctx).Where("email= ?", email).First(&GetUser); res.Error != nil {
		return model.User{}, fmt.Errorf("error in login")
	}

	return GetUser, nil
}
func (s *Storage) SaveEx(ctx context.Context, ex string, uid float64) (int64, []string, error) {
	subEx, err := postfix.Postfix(ex)
	if err != nil {
		return 0, nil, fmt.Errorf("error in expression")
	}
	StringSubEx := postfix.SaveSlice(subEx)
	NewEx := model.Expression{UserId: uint(uid), EX: ex, SubEx: StringSubEx, TimeUpdate: time.Now()}
	if res := s.db.WithContext(ctx).Create(&NewEx); res.Error != nil {
		return 0, nil, fmt.Errorf("error in save ex in db")
	}
	return int64(NewEx.ID), subEx, nil
}
func (s *Storage) UpdateSubEx(ctx context.Context, idex int64, userid float64, subEx []string) error {
	ex := &model.Expression{}

	if err := s.db.WithContext(ctx).Where("user_id = ? AND id = ?", uint(userid), uint(idex)).Find(ex).Error; err != nil {
		return fmt.Errorf("failed to find expression: %v", err)

	}

	StringSubEx := postfix.SaveSlice(subEx)

	ex.SubEx = StringSubEx
	ex.TimeUpdate = time.Now()
	if err := s.db.WithContext(ctx).Save(ex).Error; err != nil {
		return fmt.Errorf("failed to update expression: %v", err)
	}
	return nil
}
func (s *Storage) UpdateExTimeReady(ctx context.Context, idex int64, userid float64) error {
	ex := &model.Expression{}
	if err := s.db.WithContext(ctx).Where("user_id = ? AND id = ?", uint(userid), uint(idex)).Find(ex).Error; err != nil {
		return fmt.Errorf("failed to find expression: %v", err)

	}
	ex.Timestamp = time.Now()
	ex.Ready = true
	if err := s.db.WithContext(ctx).Save(ex).Error; err != nil {
		return fmt.Errorf("failed to update expression: %v", err)
	}
	return nil
}
func (s *Storage) GetExHistory(ctx context.Context, userid float64) ([]model.Expression, error) {
	var expressions []model.Expression
	if err := s.db.Where("user_id = ?", uint(userid)).Order("timestamp desc").Find(&expressions).Error; err != nil {
		return nil, err
	}

	return expressions, nil
}
func (s *Storage) GetUnREadyEx() ([]model.Expression, error) {
	var expressions []model.Expression
	if err := s.db.Where("ready = ? AND error = ?", false, false).Find(&expressions).Error; err != nil {
		return nil, err
	}
	return expressions, nil

}
func (s *Storage) ErrorinEx(uid float64, idex int64) error {
	ex := &model.Expression{}
	fmt.Println("он тут?")
	if err := s.db.Where("user_id = ? AND id = ?", uint(uid), uint(idex)).Find(ex).Error; err != nil {
		return fmt.Errorf("failed to find expression: %v", err)

	}

	ex.Error = true
	if err := s.db.Save(ex).Error; err != nil {
		return fmt.Errorf("failed to update expression: %v", err)
	}
	return nil
}

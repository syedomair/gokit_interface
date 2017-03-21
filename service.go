package main

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/syedomair/kit2/models"
	"os"
	"strconv"
)

type Service interface {
	PostUser(ctx context.Context, u models.User) error
	GetUser(ctx context.Context, id string) (models.User, error)
	PutUser(ctx context.Context, id string, u models.User) error
	GetUserBooks(ctx context.Context, id string, offset string, limit string, orderby string, sort string) (interface{}, string, string, string)
}

/*
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}
*/

type BookResponse struct {
	Id          int64  `json:"id"`
	UserId      int64  `json:"user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Name        string `json:"book_name" `
	Description string `json:"description" `
	Publish     bool   `json:"publish" `
}

type dbService struct {
	db *gorm.DB
}

func DBService() Service {

	var Db *gorm.DB
	/*
		db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			logger.Log("err", "The data source arguments are not valid")
			return nil, err
		}

		err = db.Ping()
		if err != nil {
			logger.Log("err", "Database connection error")
			return nil, err
		}
		return db, err
	*/
	var err error
	Db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Println("The data source arguments are not valid")
		panic(err)
	}
	Db.SingularTable(true)

	Db.DB().SetMaxIdleConns(3)
	Db.DB().SetMaxOpenConns(10)
	Db.LogMode(true)

	//The following lines are needed when Heroku drops all the tables periodically.
	if !Db.HasTable(&models.Client{}) {
		Db.CreateTable(&models.Client{})
		client := models.Client{Name: "Test",
			ApiKey:    "dHb%e@Bg0f8-API_KEY-&bE71jKoH=2",
			ApiSecret: "g$5%6kQ56-API_SECRET-6gE@7&EbR2",
			Active:    true}
		Db.NewRecord(client)
		Db.Create(&client)
	}
	if !Db.HasTable(&models.User{}) {
		Db.CreateTable(&models.User{})
		user := models.User{FirstName: "John",
			LastName: "Smith",
			Email:    "john@gmail.com",
			Password: "123"}
		Db.NewRecord(user)
		Db.Create(&user)
	}
	if !Db.HasTable(&models.Book{}) {
		Db.CreateTable(&models.Book{})
		book := models.Book{UserId: 1,
			Name:        "Test Book",
			Description: "Test Book Description",
			Publish:     true}
		Db.NewRecord(book)
		Db.Create(&book)
	}

	return &dbService{
		db: Db,
	}
}

func (d *dbService) PostUser(ctx context.Context, u models.User) error {
	d.db.NewRecord(u)
	d.db.FirstOrCreate(&u, u)
	return nil
}

func (d *dbService) GetUser(ctx context.Context, id string) (models.User, error) {
	u := models.User{}
	d.db.Table("public.user as u").
		Select("*").
		Where("u.id = ?", id).
		Scan(&u)

	return u, nil
}

func (d *dbService) GetUserBooks(ctx context.Context, id string, offset string, limit string, orderby string, sort string) (interface{}, string, string, string) {

	orderby = "book." + orderby
	count := 0
	bookResponse := []BookResponse{}
	d.db.Table("book").
		Select("*").
		Joins("join public.user as u on book.user_id = u.id").
		Count(&count).
		Limit(limit).
		Offset(offset).
		Order(orderby+" "+sort).
		Where("book.user_id = ?", id).
		Scan(&bookResponse)

	return bookResponse, offset, limit, strconv.Itoa(count)
}

func (d *dbService) PutUser(ctx context.Context, id string, u models.User) error {
	d.db.First(&u, id)
	d.db.Model(&u).Updates(&u)

	return nil
}

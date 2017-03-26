package main

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/syedomair/gokit_interface/models"
	"os"
	"strconv"
	"strings"
)

type Service interface {
	PostUser(ctx context.Context, u models.User) error
	PostAuthenticate(ctx context.Context, u models.User) (UserResponse, error)
	GetUser(ctx context.Context, id string) (UserResponse, error)
	PutUser(ctx context.Context, id string, u models.User) error
	PatchBook(ctx context.Context, id string, u models.Book) error
	GetUserBooks(ctx context.Context, id string, offset string, limit string, orderby string, sort string) (interface{}, string, string, string)
	AuthProvider(xkey string, xjwt string, url_path string) map[string]interface{}
	PostBook(ctx context.Context, b models.Book) error
	GetBook(ctx context.Context, id string) (BookResponse, error)
	GetBooks(ctx context.Context, offset string, limit string, orderby string, sort string) (interface{}, string, string, string)
	GetPublicBooks(ctx context.Context, offset string, limit string, orderby string, sort string) (interface{}, string, string, string)
}

type dbService struct {
	db *gorm.DB
}

func DBService() Service {

	var Db *gorm.DB
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
			Password: "MQ=="}
		Db.NewRecord(user)
		Db.FirstOrCreate(&user, models.User{Email: user.Email})
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

/* USER  */
type UserResponse struct {
	Id        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (d *dbService) PostUser(ctx context.Context, u models.User) error {
	d.db.NewRecord(u)
	d.db.FirstOrCreate(&u, models.User{Email: u.Email})
	return nil
}

func (d *dbService) PostAuthenticate(ctx context.Context, u models.User) (UserResponse, error) {

	userResponse := UserResponse{}
	d.db.Table("public.user as u").
		Select("*").
		Where("email = ? AND password = ?", u.Email, u.Password).
		Scan(&userResponse)

	return userResponse, nil
}

func (d *dbService) GetUser(ctx context.Context, id string) (UserResponse, error) {
	u := UserResponse{}
	d.db.Table("public.user as u").
		Select("*").
		Where("u.id = ?", id).
		Scan(&u)

	return u, nil
}

func (d *dbService) PutUser(ctx context.Context, id string, u models.User) error {
	d.db.First(&u, id)
	d.db.Model(&u).Updates(&u)

	return nil
}

/* BOOK */
type BookResponse struct {
	Id          int64  `json:"id"`
	UserId      int64  `json:"user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Name        string `json:"book_name" `
	Description string `json:"description" `
	Publish     bool   `json:"publish" `
}

func (d *dbService) PostBook(ctx context.Context, b models.Book) error {
	d.db.NewRecord(b)
	d.db.FirstOrCreate(&b, b)
	return nil
}

func (d *dbService) PatchBook(ctx context.Context, id string, b models.Book) error {
	cb := models.Book{}

	d.db.First(&cb, id)
	cb.Name = b.Name
	cb.Description = b.Description
	cb.Publish = b.Publish
	d.db.Save(&cb)

	return nil
}

func (d *dbService) GetBook(ctx context.Context, id string) (BookResponse, error) {
	bookResponse := BookResponse{}
	d.db.Table("book").
		Select("*").
		Joins("join public.user as u on book.user_id = u.id").
		Where("book.id = ?", id).
		Scan(&bookResponse)

	return bookResponse, nil

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

func (d *dbService) GetBooks(ctx context.Context, offset string, limit string, orderby string, sort string) (interface{}, string, string, string) {

	orderby = "book." + orderby
	count := 0
	bookResponse := []BookResponse{}
	d.db.Table("book").
		Select("*").
		Joins("join public.user as u on book.user_id = u.id").
		Count(&count).
		Limit(limit).
		Offset(offset).
		Order(orderby + " " + sort).
		Scan(&bookResponse)

	return bookResponse, offset, limit, strconv.Itoa(count)
}

func (d *dbService) GetPublicBooks(ctx context.Context, offset string, limit string, orderby string, sort string) (interface{}, string, string, string) {

	orderby = "book." + orderby
	count := 0
	bookResponse := []BookResponse{}
	d.db.Table("book").
		Select("*").
		Joins("join public.user as u on book.user_id = u.id").
		Where("book.publish = ?", true).
		Count(&count).
		Limit(limit).
		Offset(offset).
		Order(orderby + " " + sort).
		Scan(&bookResponse)

	return bookResponse, offset, limit, strconv.Itoa(count)
}
func (d *dbService) AuthProvider(xkey string, xjwt string, url_path string) map[string]interface{} {

	apiKey := xkey
	jwtToken := xjwt

	publicEndPoint := false
	if strings.Contains(url_path, "/public/") {
		publicEndPoint = true
	}

	if apiKey == "" {
		return errorAuthResponse("Header missing: x-key ")
	}
	if jwtToken == "" {
		return errorAuthResponse("Header missing: x-jwt ")
	}

	client := models.Client{}
	d.db.Table("public.client as c").
		Select("*").
		Where("c.api_key = ?", apiKey).
		Scan(&client)

	if client.ApiSecret == "" {
		return errorAuthResponse("Invalid api_key ")
	}

	type Claims struct {
		Username string `json:"username"`
		Password string `json:"password"`
		jwt.StandardClaims
	}

	tokenClaims := Claims{}

	_, err := jwt.ParseWithClaims(jwtToken, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(client.ApiSecret), nil
	})

	if err != nil {
		return errorAuthResponse("Invalid JWT Signature")
	}

	if !publicEndPoint {
		if (tokenClaims.Username != "new_registration") && (tokenClaims.Password != "new_registration") {
			user := models.User{}

			d.db.Where("email = ? AND password = ?", tokenClaims.Username, tokenClaims.Password).Find(&user)
			if user.FirstName == "" {
				return errorAuthResponse("Invalid Email or Password ")
			}
		}
	}

	return nil
}

package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
)

const (
	ImgDir   = "../images"
	DataBase = "../../db/mercari.sqlite3"
)

type Response struct {
	Message string `json:"message"`
}

type Item struct {
	Name      string `json:"name"`
	Category  string `json:"category"`
	ImageName string `json:"image_name"`
}

type ItemsList struct {
	Items []Item `json:"items"`
}

func root(c echo.Context) error {
	res := Response{Message: "Hello, world!"}
	return c.JSON(http.StatusOK, res)
}

func addItem(c echo.Context, db *sql.DB) error {
	// Get form data
	name := c.FormValue("name")
	category := c.FormValue("category")

	// Receive image files
	file, err := c.FormFile("imageName")
	if err != nil {
		c.Logger().Errorf("Error receiving image: %s", err)
		return err
	}
	c.Logger().Infof("Receive item: %s", name)

	// Open file
	src, err := file.Open()
	if err != nil {
		c.Logger().Errorf("Error opening file: %s", err)
		return err
	}
	defer src.Close()

	// Read file and calculate hash value
	hash := sha256.New()
	if _, err := io.Copy(hash, src); err != nil {
		c.Logger().Errorf("Error reading file: %s", err)
		return err
	}
	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	// Generate file names from hash values
	imageName := hashString + ".jpg"

	// Save images in the images directory
	dst, err := os.Create(filepath.Join(ImgDir, imageName))
	if err != nil {
		c.Logger().Errorf("Error creating image: %s", err)
		return err
	}
	defer dst.Close()

	// move the file pointer back to the beginning
	src.Seek(0, io.SeekStart)
	if _, err := io.Copy(dst, src); err != nil {
		c.Logger().Errorf("Error  moving file pointer: %s", err)
		return err
	}

	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO items(name, category, image_name) VALUES (?, ?, ?)")
	if err != nil {
		c.Logger().Errorf("Error preparing statement: %s", err)
		return err
	}
	defer stmt.Close()
	// Execute
	if _, err := stmt.Exec(name, category, imageName); err != nil {
		c.Logger().Errorf("Error inserting items: %s", err)
		return err
	}

	c.Logger().Infof("Receive item: %s", name)
	message := fmt.Sprintf("Item added to database: %s; Category: %s", name, category)
	res := Response{Message: message}

	return c.JSON(http.StatusOK, res)
}

func getItems(c echo.Context, db *sql.DB) error {
	rows, err := db.Query(`SELECT * FROM items`)
	if err != nil {
		c.Logger().Errorf("Error retrieving items: %s", err)
		return err
	}
	var itemsList ItemsList
	for rows.Next() {
		var item Item
		var id int
		if err := rows.Scan(&id, &item.Name, &item.Category, &item.ImageName); err != nil {
			c.Logger().Errorf("Error scanning items: %s", err)
			return err
		}
		itemsList.Items = append(itemsList.Items, item)
	}
	return c.JSON(http.StatusOK, itemsList)
}

func getItemById(c echo.Context, db *sql.DB) error {
	id, _ := strconv.Atoi(c.Param("id"))
	row := db.QueryRow(`SELECT * FROM items WHERE id = ?;`, id)
	var item Item
	if err := row.Scan(&id, &item.Name, &item.Category, &item.ImageName); err != nil {
		if err == sql.ErrNoRows {
			c.Logger().Errorf("Error: no row")
			return err
		}
		c.Logger().Errorf("Error in scanning: %s", err)
	}
	return c.JSON(http.StatusOK, item)
}

func getImg(c echo.Context) error {
	// Create image path
	imgPath := path.Join(ImgDir, c.Param("imageFilename"))

	if !strings.HasSuffix(imgPath, ".jpg") {
		res := Response{Message: "Image path does not end with .jpg"}
		return c.JSON(http.StatusBadRequest, res)
	}
	if _, err := os.Stat(imgPath); err != nil {
		c.Logger().Debugf("Image not found: %s", imgPath)
		imgPath = path.Join(ImgDir, "default.jpg")
	}
	return c.File(imgPath)
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.DEBUG)

	frontURL := os.Getenv("FRONT_URL")
	if frontURL == "" {
		frontURL = "http://localhost:3000"
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{frontURL},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Open database
	db, err := sql.Open("sqlite3", DataBase)
	if err != nil {
		log.Fatal("unable to use data source name", err)
	}
	defer db.Close()

	// Routes
	e.GET("/", root)
	e.GET("/items", func(c echo.Context) error {
		return getItems(c, db)
	})
	e.POST("/items", func(c echo.Context) error {
		return addItem(c, db)
	})
	e.GET("/items/:id", func(c echo.Context) error {
		return getItemById(c, db)
	})
	e.GET("/image/:imageFilename", getImg)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}

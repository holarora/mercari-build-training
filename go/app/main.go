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
	DataBase = "../db/mercari.sqlite3"
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

func getCategoryId(db *sql.DB, category string) (int, error) {
	var categoryId int
	query := `
		SELECT id
		FROM categories
		WHERE name = ?
	`
	err := db.QueryRow(query, category).Scan(&categoryId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, err
		}
		return 0, err
	}
	return categoryId, nil
}

func addItem(c echo.Context, db *sql.DB) error {
	// Get form data
	name := c.FormValue("name")
	category := c.FormValue("category")
	categoryId, err := getCategoryId(db, category)
	if err != nil {
		c.Logger().Errorf("Category does not exist: %s", err)
		return err
	}

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

	// Move the file pointer back to the beginning
	src.Seek(0, io.SeekStart)
	if _, err := io.Copy(dst, src); err != nil {
		c.Logger().Errorf("Error moving file pointer: %s", err)
		return err
	}

	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO items(name, category_id, image_name) VALUES (?, ?, ?)")
	if err != nil {
		c.Logger().Errorf("Error preparing statement: %s", err)
		return err
	}
	defer stmt.Close()
	// Execute
	if _, err := stmt.Exec(name, categoryId, imageName); err != nil {
		c.Logger().Errorf("Error inserting items: %s", err)
		return err
	}

	c.Logger().Infof("Receive item: %s", name)
	message := fmt.Sprintf("Item added to database: %s; Category: %s", name, category)
	res := Response{Message: message}

	return c.JSON(http.StatusOK, res)
}

func getItems(c echo.Context, db *sql.DB) error {
	query := `
		SELECT items.name, categories.name, items.image_name
		FROM items
		INNER JOIN categories ON items.category_id = categories.id;
	`
	rows, err := db.Query(query)
	if err != nil {
		c.Logger().Errorf("Error retrieving items: %s", err)
		return err
	}
	defer rows.Close()

	var itemsList ItemsList
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.Name, &item.Category, &item.ImageName); err != nil {
			c.Logger().Errorf("Error scanning items: %s", err)
			return err
		}
		itemsList.Items = append(itemsList.Items, item)
	}
	if err := rows.Err(); err != nil {
		c.Logger().Errorf("Error retrieving items: %s", err)
		return err
	}
	return c.JSON(http.StatusOK, itemsList)
}

func getItemById(c echo.Context, db *sql.DB) error {
	id, _ := strconv.Atoi(c.Param("id"))
	query := `
		SELECT items.name, categories.name, items.image_name
		FROM items
		INNER JOIN categories ON items.category_id=categories.id
		WHERE items.id = ?;
	`
	row := db.QueryRow(query, id)
	var item Item
	if err := row.Scan(&item.Name, &item.Category, &item.ImageName); err != nil {
		if err == sql.ErrNoRows {
			c.Logger().Errorf("Error: no row")
			return err
		}
		c.Logger().Errorf("Error in scanning: %s", err)
		return err
	}
	c.Logger().Infof("Item found: %+v", item)
	return c.JSON(http.StatusOK, item)
}

func getItemByKeyWord(c echo.Context, db *sql.DB) error {
	keyword := c.QueryParam("keyword")
	query := `
		SELECT items.name, categories.name, items.image_name
		FROM items
		INNER JOIN categories ON items.category_id = categories.id
		WHERE items.name LIKE ?;
	`
	rows, err := db.Query(query, keyword)
	if err != nil {
		c.Logger().Errorf("Error during query: %s", err)
		return err
	}
	defer rows.Close()

	var itemsList ItemsList
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.Name, &item.Category, &item.ImageName); err != nil {
			c.Logger().Errorf("Error scanning items: %s", err)
			return err
		}
		itemsList.Items = append(itemsList.Items, item)
	}
	if err := rows.Err(); err != nil {
		c.Logger().Errorf("Error retrieving items: %s", err)
		return err
	}
	return c.JSON(http.StatusOK, itemsList)
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
	e.GET("/search", func(c echo.Context) error {
		return getItemByKeyWord(c, db)
	})

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}

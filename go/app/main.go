package main

import (
	"net/http"
	"os"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func insertNewclm(db *sql.DB, categoryName string, c echo.Context) error {
		name, image, err := getNewItemsForClm(c)
		if err != nil {
			c.Logger().Errorf("getItemFromRequest error: %v", err)
			return err
		}
		image_name,err := ImgSave(image)
		if err != nil {
			c.Logger().Errorf("ImgSave error: %v",err)
			return err
		}
		cmd := "INSERT INTO categories (name) VALUES(?)"
		stmt, err := db.Prepare(cmd)
		if err != nil {
			c.Logger().Errorf("db.Prepare error : %v\n",err)
			return err
		}
		addQuery, err := stmt.Exec(categoryName)
		if err != nil {
			c.Logger().Errorf("addQuery error : %v", err)
			return err
		}
		category_id, err := addQuery.LastInsertId()
		if err != nil {
			c.Logger().Errorf("getting last insert ID error : %v",err)
			return err
		}

		ItemQuery := "INSERT INTO items (name,category_id,image_name) VALUES (?,?,?)"
		insertQuery, err := db.Prepare(ItemQuery)
		if err != nil {
			c.Logger().Errorf("db.Prepare error: %v", err)
			return err
		}
		defer insertQuery.Close()
		//you can change param easily nothing to change sql-statement
		if _,err = insertQuery.Exec(name,category_id,image_name); err != nil{
			c.Logger().Errorf("insertQuery.Exec error", err)
			return err
		}

		return nil
}

func addItem(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		name, category, image, err := getItemFromRequest(c)
		// fmt.Printf("getItemFrom!!! %v\n", category)
		if err != nil {
			c.Logger().Errorf("getItemFromRequest error")
			return err
		}
		image_name,err := ImgSave(image)
		if err != nil {
			c.Logger().Errorf("ImgSave error: %v\n", err)
			return err
		}
		// fmt.Printf("getItemFrom??? %v\n", category)
		category_id, err := getCategoryID(db, category,c)
		
		if err != nil {
			c.Logger().Errorf("getCategoryID error %v:",err)
			return err
		}
		//insert date into a table
		//name(categories)->id(categories)->category_id(items)
		cmd := "INSERT INTO items (name, category_id, image_name) VALUES(?, ?, ?)"
		
		//prepare sql-statement
		stmt, err := db.Prepare(cmd)
		if err != nil {
			c.Logger().Errorf("db.Prepare error: %v", err)
			return err
		}
		defer stmt.Close()
		//Exec(sql-statement is executed)(get Exec's arguments into ↑?,?,?)
		//you can change param easily nothing to change sql-statement
		if _,err = stmt.Exec(name,category_id,image_name); err != nil{
			c.Logger().Errorf("stmt.Exec error", err)
			return err
		}
		// defer db.Close()
		message := "You've successfully added items"
		res := Response{Message: message}
		return c.JSON(http.StatusOK, res)
	}
}


func main() {
	e := echo.New()

	// Middleware 
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.INFO)

	frontURL := os.Getenv("FRONT_URL")
	if frontURL == "" {
		frontURL = "http://localhost:3000"
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{frontURL},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	///db := InitDB() <- good connect-DB-code better than sqlOpen() in more complex situation
	db, err := sqlOpen()
	if err != nil {
		//%v can be putted any data-type,%s is string-type only though
		e.Logger.Infof("Failed to open the database: %v",err)
	}
	defer db.Close()
	// e.GET("/", root)
	//curl -X POST --url "http://localhost:9000/items" -F "name=jacket" -F "category=fashion" -F "image=@/mnt/c/Users/ff102/mercari/mercari-build-training//go/images/default.jpg"
	e.POST("/items", addItem(db))
	//curl http://localhost:9000/items
	e.GET("/items",getItems)
	//curl -X GET 'http://localhost:9000/image/default.jpg' -o output.jpg
	e.GET("/image/:imageFilename", getImg)
	//curl http://localhost:9000/items/3
	e.GET("/items/:item_id", getItemByItemId(db))
	// curl http://localhost:9000/search?keyword=jacket
	e.GET("/search",searchItems)
	// Start server
	e.Logger.Fatal(e.Start(":9000")) 
}


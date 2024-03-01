package main

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/labstack/echo/v4"
	"mime/multipart"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func getItemFromRequest(c echo.Context) (string,string,*multipart.FileHeader, error){
	//Get date from request
	name := c.FormValue("name")
	category := c.FormValue("category")
	image, err := c.FormFile("image")
	if err != nil{
		c.Logger().Errorf("FormFile error : %v\n", err)
		return "", "", nil, err
	}
	return name, category, image, nil
}

func getNewItemsForClm(c echo.Context) (string,*multipart.FileHeader, error) {
	//Get date from request
	name := c.FormValue("name")
	image, err := c.FormFile("image")
	if err != nil{
		c.Logger().Errorf("FormFile error : %v\n", err)
		return "", nil, err
	}
	return name, image, nil
}

func getItemByItemId(db *sql.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        itemId, err := strconv.Atoi(c.Param("item_id"))
        if err != nil {
            // Use echo's HTTP error to return the error to the client properly.
            return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid item ID: %v", err))
        }
        var items Item
		//PLEASE BE CAREFULL TO BE THE RIGHT SEQUENCE
        cmd := `SELECT items.id, items.name, categories.name, items.category_id ,items.image_name FROM items LEFT OUTER JOIN categories ON items.category_id = categories.id WHERE items.id = ?`

        err = db.QueryRow(cmd, itemId).Scan(&items.ID, &items.Name, &items.CategoryName,&items.CategoryID,&items.ImageName)
		// fmt.Printf("item.ID %v",items.ID)
        if err != nil {
            if err == sql.ErrNoRows {
                // Item not found
                return echo.NewHTTPError(http.StatusNotFound, "Item not found")
            }
            // Internal server error for other types of errors
            return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
        }

        return c.JSON(http.StatusOK, items)
    }
}

func getItems (c echo.Context) error {
		//connect to db
		db, err := sqlOpen()
		if err != nil {
			c.Logger().Errorf("sql.Open error")
			return err
		}
		defer db.Close()
		//get db-date(return sql.Rows object)
		// cmd := "SELECT items.id ,items.name, categories.name, items.image_name FROM items JOIN categories ON items.category_id = categories.id"
		cmd := "SELECT items.id, items.name, categories.name AS category, items.image_name FROM items JOIN categories ON items.category_id = categories.id"

		// cmd := "SELECT items.name, categories.name, items.image_name FROM items JOIN categories ON items.category_id = categories.id"
		rows, err := db.Query(cmd)
		if err != nil {
			c.Logger().Errorf("db.Query error")
			return err
		}
		defer rows.Close()
		//get db-date into item
		// var items []Item
		var items []ResponseItem
		for rows.Next() {
			//declare item instance for storage
			// var item Item
			var item ResponseItem
			//store rows-date to item-field
			// if err := rows.Scan(&item.Name,&item.CategoryName,&item.ImageName); err != nil{
			if err := rows.Scan(&item.Id, &item.Name, &item.Category, &item.Image_name); err != nil{
				c.Logger().Errorf("rows.Scan: rows-date storing error")
				return err
			}
			items = append(items,item)
		}
		//check roup error
		if err = rows.Err(); err != nil {
			c.Logger().Errorf("for roup error(rows.Next())")
			return err
		}
		
		// return rows-date
		// c.Logger().Errorf("Array %v", ResponseItems{Items: items})
		// fmt.Printf("Array %v", ResponseItems{Items: items})
		//こっちだとPrototypeがobjectで送っていて、
		//Uncaught TypeErrorが起きない
		return c.JSON(http.StatusOK, ResponseItems{Items: items})
		//↓こっちだとPrototypeが配列で送っている
		// return c.JSON(http.StatusOK, items)
}

func searchItems(c echo.Context) error {
	//receive keywords from request
	keyword := c.QueryParam("keyword")
	//connect to db
	db, err := sqlOpen()
	if err != nil {
		c.Logger().Errorf("sql.Open error")
		return err
	}
	defer db.Close()
	//search keywords to items in db
	cmd := "SELECT items.id ,items.name , items.category_id ,categories.name, items.image_name FROM items JOIN categories ON items.category_id = categories.id WHERE items.name LIKE ?"
	// cmd := "SELECT categories.name FROM categories WHERE name LIKE ?"
	rows, err := db.Query(cmd,"%"+keyword+"%")
	// fmt.Printf("err! %v",rows)
	// fmt.Printf("err? %v", err)
	// fmt.Printf("err# %v", keyword)

	if err != nil {
		c.Logger().Errorf("db.Query error")
		return err
	}
	defer rows.Close()
	//get items matched it
	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID,&item.Name, &item.CategoryID, &item.CategoryName, &item.ImageName); err != nil {
			c.Logger().Errorf("rows.Scan error")
			return err
		}
		items = append(items, item)
	}
	//return its items
	return c.JSON(http.StatusOK, items)
}

func getCategoryID(db *sql.DB, categoryName string, c echo.Context) (int, error) {
	//prepare variable for id
	var categoryId int
	//Query.Row id
	cmd := "SELECT id FROM categories where name = ?"
	//run sql-stmt and assigning the retreived data(because got data by sql-stmt already) to variable
	row := db.QueryRow(cmd, categoryName)
	err := row.Scan(&categoryId)
	if err != nil {
		//no wanted rows error
		if err == sql.ErrNoRows {
			err := insertNewclm(db, categoryName, c)
			if err != nil {
				return 0, fmt.Errorf("insertNewclm error : %v", err)
			}
		}
		//ex.sql-stmt error & connection error
		addItem(db)
	}
	return categoryId, nil
}

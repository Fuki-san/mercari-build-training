package main

import(
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"github.com/labstack/echo/v4"
	"crypto/sha256"
	"io"
	"path/filepath"
	"encoding/hex"
	"mime/multipart"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func sqlOpen()(*sql.DB, error){
	db, err := sql.Open("sqlite3",Sqlpath)
	if err != nil {
		return nil, fmt.Errorf("sql.Open error: %v",err)
	}
	return db, nil
}

func ImgSave(image *multipart.FileHeader)(string, error) {
		src, err := image.Open()
		if err != nil {
			return "", fmt.Errorf("image.Open error: %v", err)
		}
		defer src.Close()
		//receive hashedImage
		hashedString, err := getHash(src)
		if err != nil{
			return "", fmt.Errorf("getHash error: %v",err)
		}
		
		image_name := hashedString + ".jpg"
		if _, err := src.Seek(0,0); err != nil {
			return "", fmt.Errorf("failed to seek file: %v", err)
		}
		//create file path
		fmt.Printf("ImgDir: %v",ImgDir)
		fmt.Printf("Image_name: %v",image_name)
		fmt.Printf("file_name: %v",filepath.Join(ImgDir, image_name))
		dst, err := os.Create(filepath.Join(ImgDir, image_name))
		if err != nil{
			return "", fmt.Errorf("os.Create filepath.Join error: %v",err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst,src); err != nil {
			return "", fmt.Errorf("io.Copy error: %v",err)
		}

		return image_name, nil
} 

func getHash(src io.Reader) (string, error) {
	hash := sha256.New()
	if _,err := io.Copy(hash,src); err != nil{
		return "", fmt.Errorf("io.Copy error : %v",err)
	}
	HashedString := hex.EncodeToString(hash.Sum(nil))
	return HashedString, nil
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

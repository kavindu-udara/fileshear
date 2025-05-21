package fileserver

import (
	"flag"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Get current working directory
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

func getMainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "File Shaar",
	})
}

func getUploadPage(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", gin.H{
		"title": "File Shaar",
	})
}

func getFiles(c *gin.Context) {
	currentDir := getCurrentDir()

	// read the contents of the current dir
	files, err := os.ReadDir(currentDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to read directory",
			"error":   err.Error(),
		})
		return
	}

	// Create a slice to store filenames
	var fileList []map[string]interface{}

	// collect all file information
	for _, file := range files {
		if !file.IsDir() { // only include files, not dirs
			info, err := file.Info()
			if err != nil {
				continue
			}
			fileList = append(fileList, map[string]interface{}{
				"name":     file.Name(),
				"size":     info.Size(),
				"modified": info.ModTime(),
			})
		}
	}

	// return the filenames as JSON
	c.JSON(http.StatusOK, gin.H{
		"files": fileList,
	})

}

func getFileByName(c *gin.Context) {
	fileName := c.Param("filename")
	currentDir := getCurrentDir()

	// construct file path
	filePath := filepath.Join(currentDir, fileName)

	// check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "File not found",
		})
		return
	}

	// File exists, serve it
	c.File(filePath)
}

// upload file
func uploadFile(c *gin.Context) {
	currentDir := getCurrentDir()
	// Set a reasonable max file size (e.g., 32 MB)
	c.Request.ParseMultipartForm(32 << 20)

	// Get the file from the form data
	file, header, err := c.Request.FormFile("uploadFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No file uploaded",
			"error":   err.Error(),
		})
		return
	}
	defer file.Close()

	// Ensure the files directory exists
	// Create the destination file in current directory
	dst := filepath.Join(currentDir, header.Filename)
	out, err := os.Create(dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create file",
			"error":   err.Error(),
		})
		return
	}
	defer out.Close()

	// Copy the uploaded file to the destination file
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to save file",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": header.Filename,
	})

}

func API(LanIp string) {
	// Define command line flags
	dev := flag.Bool("dev", false, "Run in development mode")
	flag.Parse()

	router := gin.Default()

	// Set web directory based on environment
	var webDir string
	if *dev {
		webDir = "web" // Development path
		gin.SetMode(gin.DebugMode)
	} else {
		webDir = "/usr/local/share/fileshear/web" // Production path
		gin.SetMode(gin.ReleaseMode)
	}

	// Load HTML templates
	router.LoadHTMLGlob(filepath.Join(webDir, "*.html"))

	// serve static files
	router.Static("/scripts", filepath.Join(webDir, "scripts"))
	router.Static("/styles", filepath.Join(webDir, "styles"))

	router.GET("/", getMainPage)
	router.GET("/upload", getUploadPage)

	// get files list
	router.GET("/files", getFiles)

	router.GET("/files/:filename", getFileByName)

	// upload file
	router.POST("/files", uploadFile)

	router.Run(LanIp + ":8080")

}

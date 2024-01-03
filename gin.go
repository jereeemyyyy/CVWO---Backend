package main

import (
        "database/sql"
        "strconv"
        "fmt"
        "log"
        "os"
        "net/http"
        "github.com/gin-gonic/gin"

        
        "github.com/go-sql-driver/mysql"
        
)
var db *sql.DB

// Structure of album
type Album struct {
        ID     int64
        Title  string
        Artist string
        Price  float32
    }

func main() {
        // Capture connection properties.
        cfg := mysql.Config{
                User:   os.Getenv("DBUSER"),
                Passwd: os.Getenv("DBPASS"),
                Net:    "tcp",
                Addr:   "127.0.0.1:3306",
                DBName: "recordings",
        }
        // Get a database handle.
        var err error
        db, err = sql.Open("mysql", cfg.FormatDSN())
        if err != nil {
                log.Fatal(err)
        }

        pingErr := db.Ping()
        if pingErr != nil {
                log.Fatal(pingErr)
        }
        fmt.Println("Connected!")

       

// TODO: move to functions
        router := gin.Default()
        router.GET("/albums", getAlbum)
        router.GET("/albums/:id", getAlbumByID)
        router.POST("/albums", addAlbum)
        router.Run("localhost:8080")


}

// getAlbums responds with the list of all albums as JSON.
func getAlbum(c *gin.Context) {
//
        // An albums slice to hold data from returned rows.
        var albums []Album
    
        rows, err := db.Query("SELECT * FROM album")
        if err != nil {
             c.IndentedJSON(http.StatusBadRequest, err)
        }
        defer rows.Close()
        // Loop through rows, using Scan to assign column data to struct fields.
        for rows.Next() {
            var alb Album
            if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
                c.IndentedJSON(http.StatusBadRequest, err)
            }
            albums = append(albums, alb)
        }
        if err := rows.Err(); err != nil {
                c.IndentedJSON(http.StatusBadRequest, err)
        }

        c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumByID(c *gin.Context) {
        id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
        alb, err := albumByIDQuery(id)

        if err != nil {
                c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
        } else {
                c.IndentedJSON(http.StatusOK, alb)
        }
                
                
                

        
}
        

func albumsByArtist(c *gin.Context) {
        artist := c.Param("artist")
        alb, err := albumsByArtistQuery(artist)

        if err != nil {
                c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
        } else {
                c.IndentedJSON(http.StatusOK, alb)     
        }
        
}



func addAlbum(c *gin.Context) {
        var newAlbum Album

        if err := c.BindJSON(&newAlbum); err != nil {
                return
        }
        newAlbum.ID , _ = addAlbumQuery(newAlbum)
        c.IndentedJSON(http.StatusCreated, newAlbum)

}




//QUERY FUNCTIONS

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtistQuery(name string) ([]Album, error) {
        // An albums slice to hold data from returned rows.
        var albums []Album
    
        rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
        if err != nil {
            return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
        }
        defer rows.Close()
        // Loop through rows, using Scan to assign column data to struct fields.
        for rows.Next() {
            var alb Album
            if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
                return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
            }
            albums = append(albums, alb)
        }
        if err := rows.Err(); err != nil {
            return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
        }
        return albums, nil
    }

// albumByID queries for the album with the specified ID.
func albumByIDQuery(id int64) (Album, error) {
    // An album to hold data from the returned row.
    var alb Album

    row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
    // Declaring err variable, and if statement in 1 line
    if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
        if err == sql.ErrNoRows {
            return alb, fmt.Errorf("albumsById %d: no such album", id)
        }
        return alb, fmt.Errorf("albumsById %d: %v", id, err)
    }
    return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func addAlbumQuery(alb Album) (int64, error) {
        result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
        if err != nil {
            return 0, fmt.Errorf("addAlbum: %v", err)
        }
        id, err := result.LastInsertId()
        if err != nil {
            return 0, fmt.Errorf("addAlbum: %v", err)
        }
        return id, nil
    }

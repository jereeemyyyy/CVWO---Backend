package main

import (
        "database/sql"
        "strconv"
        "fmt"
        "log"
        "os"
        "net/http"
        "github.com/gin-gonic/gin"

        
        _ "github.com/lib/pq"
        
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
        username := os.Getenv("DBUSER")
        password := os.Getenv("DBPASS")
        album := os.Getenv("ALBUM")

        connStr := fmt.Sprintf("user=%s name=%s password=%s sslmode=disable", username, album, password)

        // Get a database handle.
        var err error
        db, err = sql.Open("postgres", connStr)
        if err != nil {
                log.Fatal(err)
        }
        //defer db.Close()

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
        router.Run("0.0.0.0:8080")


}

// getAlbums responds with the list of all albums as JSON.
func getAlbum(c *gin.Context) {
//
        // An albums slice to hold data from returned rows.
        var albums []Album
    
        rows, err := db.Query("SELECT * FROM yourmom")
        if err != nil {
             c.IndentedJSON(http.StatusBadRequest, err)
             return;
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
        var err error
        newAlbum.ID , err = addAlbumQuery(newAlbum)
        if err != nil {
                fmt.Println(err)
        }
        c.IndentedJSON(http.StatusCreated, newAlbum)

}




//QUERY FUNCTIONS

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtistQuery(name string) ([]Album, error) {
        // An albums slice to hold data from returned rows.
        var albums []Album
    
        rows, err := db.Query("SELECT * FROM yourmom WHERE artist = ?", name)
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

    row := db.QueryRow("SELECT * FROM yourmom WHERE id = ?", id)
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
        var id int64
        err := db.QueryRow(`INSERT INTO yourmom (title, price, artist) VALUES ($1, $2, $3) RETURNING id`, alb.Title, alb.Price, alb.Artist).Scan(&id)
        if err != nil {
                fmt.Print(err)
            return -1, fmt.Errorf("addAlbum: %v", err)
        }
        
        return id, nil
    }

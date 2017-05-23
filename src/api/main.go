package main

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	//"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	//"github.com/gocraft/dbr/dialect"
)

type (
	mind struct {
		ID      int    `db:"id"`
		Created string `db:"created"`
		Name    string `db:"name"`
		Mind    string `db:"mind"`
	}

	mindJSON struct {
		ID      int    `json:"id"`
		Created string `json:"created"`
		Name    string `json:"name"`
		Mind    string `json:"lastName"`
	}

	responseData struct {
		Minds []mind `json:"names"`
	}
)

var (
	tablename = "mind_data"
	seq       = 1
	conn, _   = dbr.Open("mysql", "root:thach208@tcp(localhost:3306)/mindmap", nil)
	sess      = conn.NewSession(nil)
)

//----------
// Handlers
//----------

func insertMind(c echo.Context) error {
	u := new(mindJSON)
	if err := c.Bind(u); err != nil {
		return err
	}

	sess.InsertInto(tablename).Columns("id", "created", "name", "mind").Values(u.ID, u.Created, u.Name, u.Mind).Exec()

	return c.NoContent(http.StatusOK)
}

func selectMinds(c echo.Context) error {
	var u []mind

	sess.Select("name").From(tablename).Load(&u)
	// response := new(responseData)
	// response.Minds = u
	fmt.Println(u)
	// fmt.Println(response)
	return c.JSON(http.StatusOK, u)
}

func selectMind(c echo.Context) error {
	var m mind
	id := c.Param("id")
	sess.Select("*").From(tablename).Where("id = ?", id).Load(&m)
	//idはPrimary Key属性のため重複はありえない。そのため一件のみ取得できる。そのため配列である必要はない
	return c.JSON(http.StatusOK, m)

}

func updateMind(c echo.Context) error {
	u := new(mindJSON)
	if err := c.Bind(u); err != nil {
		return err
	}

	attrsMap := map[string]interface{}{"id": u.ID, "created": u.Created, "name": u.Name, "mind": u.Mind}
	sess.Update(tablename).SetMap(attrsMap).Where("id = ?", u.ID).Exec()
	return c.NoContent(http.StatusOK)
}

func deleteMind(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	sess.DeleteFrom(tablename).
		Where("id = ?", id).
		Exec()

	return c.NoContent(http.StatusNoContent)
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes

	e.POST("/minds/", insertMind)
	e.GET("/mind/:id", selectMind)
	e.GET("/minds/", selectMinds)
	e.PUT("/minds/", updateMind)
	e.DELETE("/minds/:id", deleteMind)

	// Start server
	// e.Run(standard.New(":1323"))
	e.Logger.Fatal(e.Start(":4000"))
}
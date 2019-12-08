package web

import (
	"encoding/base64"
	"git.bn4t.me/bn4t/dynamic-qr/app/utils"
	"git.bn4t.me/bn4t/dynamic-qr/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"strconv"
)
import "github.com/skip2/go-qrcode"

type managementData struct {
	Target   string
	Password string
	Id       int
	QrCode   string
	Link     string
}

// create a new dynamic QR code
func handleCreateQr(c echo.Context) error {
	target := c.FormValue("target")

	// check if target link is supplied
	if target == "" {
		return c.String(http.StatusBadRequest, "Bad request. Target link is missing.")
	}

	// check if the url supplied is valid
	targetUri, err := url.Parse(target)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid URL supplied.")
	}

	// concat three uuids to create a secure enough string
	password := uuid.NewV4().String() + uuid.NewV4().String() + uuid.NewV4().String()

	// add link to database
	_, err = db.Db.Exec("INSERT INTO qrcodes (password, target) VALUES (?,?)", password, targetUri.String())
	if err != nil {
		log.Print(err)
		return c.String(http.StatusInternalServerError, "Internal server error.")
	}

	// redirect the user to the manage page
	// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Redirections
	return c.Redirect(303, "/manage/"+password)
}

// update an existing QR code
func handleUpdateQr(c echo.Context) error {
	target := c.FormValue("target")
	password := c.FormValue("password")
	id := c.FormValue("id")

	// check if target link is supplied
	if target == "" {
		return c.String(http.StatusBadRequest, "Bad request. Target link is missing.")
	}

	// check if id is supplied
	if id == "" {
		return c.String(http.StatusBadRequest, "Bad request. Id is missing.")
	}

	// check if password is supplied
	if password == "" {
		return c.String(http.StatusBadRequest, "Bad request. Password is missing.")
	}

	// check if the url supplied is valid
	_, err := url.Parse(target)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid URL supplied.")
	}

	// check if id/password combination is correct and discard the value
	var discard string
	err = db.Db.QueryRow("SELECT id FROM qrcodes WHERE id=? AND password=?", id, password).Scan(&discard)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid password or id supplied.")
	}

	// add link to database
	_, err = db.Db.Exec("UPDATE qrcodes SET target = ? WHERE id = ?", target, id)
	if err != nil {
		log.Print(err)
		return c.String(http.StatusInternalServerError, "Internal server error.")
	}

	// redirect the user to the manage page
	return c.Redirect(303, "/manage/"+password)
}

// return the manage page for a QR code
func handleManage(c echo.Context) error {
	password := c.Param("password")

	// check if target link is supplied
	if password == "" {
		return c.String(http.StatusBadRequest, "Bad request.")
	}

	// retrieve necessary data from database
	var managementData managementData
	err := db.Db.QueryRow("SELECT target, id, password FROM qrcodes WHERE password=?", password).Scan(&managementData.Target, &managementData.Id, &managementData.Password)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid management link.")
	}

	link := utils.GetEnv("BASE_URL", "") + "link/" + strconv.Itoa(managementData.Id)
	managementData.Link = link

	// create QR code and encode it as base64
	qrPng, err := qrcode.Encode(link, qrcode.Medium, 256)
	if err != nil {
		log.Print(err)
		return c.String(500, "Internal server error")
	}
	managementData.QrCode = base64.StdEncoding.EncodeToString(qrPng)

	return c.Render(http.StatusOK, "manage", managementData)
}

// handle link redirects
func handleLink(c echo.Context) error {
	id := c.Param("id")

	// check if id is supplied
	if id == "" {
		return c.String(http.StatusBadRequest, "Bad request.")
	}

	// retrieve target from database
	var target string
	err := db.Db.QueryRow("SELECT target FROM qrcodes WHERE id=?", id).Scan(&target)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid id supplied.")
	}

	return c.Redirect(302, target)
}

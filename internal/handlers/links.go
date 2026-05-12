package handlers

import (
	"crypto/rand"
	"database/sql"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/paprykdev/urlshortener/internal/models"
	"github.com/paprykdev/urlshortener/internal/response"
)

const (
	Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	StringLength = 8
)

type LinksHandler struct {
	DB *sql.DB
}

func NewLinksHandler(db *sql.DB) *LinksHandler {
	return &LinksHandler{
		DB: db,
	}
}

func (h *LinksHandler) GetAll(c *gin.Context) {
	var links []models.Link

	qr := `
		SELECT * FROM links
	`

	rows, err := h.DB.Query(qr)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed querying database: "+err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		var link models.Link

		if err := rows.Scan(
			&link.ID,
			&link.Url,
			&link.ShortCode,
			&link.CreatedAt,
		); err != nil {
			response.Fail(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to scan links")
			return
		}

		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
	}

	response.Success(c, http.StatusOK, links)
}

func (h *LinksHandler) CreateLink(c *gin.Context) {
	var dto CreateLinkDTO

	if err := c.BindJSON(&dto); err != nil {
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	str, err := generateString()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	link := models.Link{
		ID: uuid.New().String(),
		Url: dto.Url,
		ShortCode: str,
		CreatedAt: time.Now(),
	}

	qr := `
		INSERT INTO links(id, url, short_code) VALUES (?, ?, ?)
	`

	res, err := h.DB.Exec(qr, link.ID, link.Url, link.ShortCode)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, res)
}

func generateString() (string, error) {
	result := make([]byte, StringLength)

	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(Charset))))
		if err != nil {
			return "", err
		}

		result[i] = Charset[n.Int64()]
	}

	return string(result), nil
}

type CreateLinkDTO struct {
	Url string `json:"url"`
}

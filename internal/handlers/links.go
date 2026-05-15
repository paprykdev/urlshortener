package handlers

import (
	"crypto/rand"
	"database/sql"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/paprykdev/urlshortener/internal/models"
	"github.com/paprykdev/urlshortener/internal/response"
)

const (
	charset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	stringLength = 8
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
		SELECT id, url, short_code, created_at FROM links
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
		return
	}

	response.Success(c, http.StatusOK, links)
}

func (h *LinksHandler) CreateLink(c *gin.Context) {
	var dto CreateLinkDTO
	var short_code string


	if err := c.BindJSON(&dto); err != nil {
		response.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	qrCheck := `
		SELECT short_code FROM Links WHERE url = ?
	`

	if err := h.DB.QueryRow(qrCheck, dto.Url).Scan(&short_code); err == nil {
		response.Success(c, http.StatusOK, gin.H{
			"short_code": short_code,
		})
		return
	} else if err != sql.ErrNoRows {
		response.Fail(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	qr := `
	INSERT INTO links(id, url, short_code) VALUES (?, ?, ?)
	`

	for {
		shortCode, err := generateString()
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
			return
		}

		_, err = h.DB.Exec(qr, uuid.New().String(), dto.Url, shortCode)

		if err == nil {
			response.Success(c, http.StatusCreated, gin.H{
				"short_code": shortCode,
			})
			return
		}
	}

}

func (h *LinksHandler) Redirect(c *gin.Context) {
	var url string

	shortCode := c.Param("code")
	if len(shortCode) != stringLength {
		response.Fail(c, http.StatusBadRequest, "INVALID_SHORT_CODE", "invalid short code")
		return
	}

	qr := `
		SELECT url FROM Links WHERE short_code = ?
	`

	if err := h.DB.QueryRow(qr, shortCode).Scan(&url); err != nil {
		response.Fail(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	c.Redirect(http.StatusMovedPermanently, url)
}

func generateString() (string, error) {
	result := make([]byte, stringLength)

	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		result[i] = charset[n.Int64()]
	}

	return string(result), nil
}

type CreateLinkDTO struct {
	Url string `json:"url" binding:"required,url"`
}


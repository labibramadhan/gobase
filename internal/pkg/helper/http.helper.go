package helper

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	modeldto "gobase/internal/model/dto"
)

type JSONIdResponse struct {
	Id string `json:"id"`
}

type JSONIdUUIDResponse struct {
	Id uuid.UUID `json:"id"`
}

type JSONResponse struct {
	Data           interface{}                `json:"data,omitempty"`
	Pagination     *modeldto.PageableDto      `json:"pagination,omitempty"`
	ResponseStatus modeldto.ResponseStatusDto `json:"response_status,omitempty"`
}

func NewJSONResponse() JSONResponse {
	// make default to handle copier error
	return JSONResponse{
		Data:           "",
		ResponseStatus: modeldto.ResponseStatusDto{},
		Pagination:     &modeldto.PageableDto{},
	}
}

func NewOkIdResponse(c *fiber.Ctx, id string) error {
	return c.JSON(fiber.Map{"data": &JSONIdResponse{Id: id}})
}

func NewOkIdCreatedResponse(c *fiber.Ctx, id string) error {
	return c.Status(http.StatusCreated).JSON(fiber.Map{"data": &JSONIdResponse{Id: id}})
}

func NewOkIdUUIDResponse(c *fiber.Ctx, id uuid.UUID) error {
	return c.JSON(fiber.Map{"data": &JSONIdUUIDResponse{Id: id}})
}

func NewOkIdUUIDCreatedResponse(c *fiber.Ctx, id uuid.UUID) error {
	return c.Status(http.StatusCreated).JSON(fiber.Map{"data": &JSONIdUUIDResponse{Id: id}})
}

func NewOkResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{"data": data})
}

func NewOkRawResponse(c *fiber.Ctx, dataWithPaging interface{}) error {
	return c.JSON(dataWithPaging)
}

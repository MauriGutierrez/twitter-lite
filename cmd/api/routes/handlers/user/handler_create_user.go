package user

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"ualaTwitter/internal/platform/httphelper"
	"ualaTwitter/internal/usecase/create_user"
)

const maxCreateUserBodySize = 2 * 1024 // 2kb

var (
	ErrInvalidCreateUserBody = errors.New("invalid request body")
	ErrEmptyUserName         = errors.New("user name cannot be empty")
	ErrInvalidDoc            = errors.New("invalid document: must be a of 7 or 8 digits")
	docRegex                 = regexp.MustCompile(`^\d{7,8}$`)
)

type createUserService interface {
	Execute(ctx context.Context, input create_user.Input) (create_user.Output, error)
}

type CreateUserHandler struct {
	service createUserService
}

func NewCreateUserHandler(service createUserService) *CreateUserHandler {
	return &CreateUserHandler{
		service: service,
	}
}

func (h *CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, maxCreateUserBodySize)

	input, err := h.parseRequest(r)
	if err != nil {
		httphelper.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	output, err := h.service.Execute(ctx, *input)
	if err != nil {
		httphelper.RenderError(w, httphelper.StatusFromError(err), err.Error())
		return
	}

	h.renderResponse(w, output.ID)

}

func (h *CreateUserHandler) parseRequest(r *http.Request) (*create_user.Input, error) {
	var req createUserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		return nil, ErrInvalidCreateUserBody
	}

	if req.Name == "" {
		return nil, ErrEmptyUserName
	}

	doc, err := validateDoc(req.Document)
	if err != nil {
		return nil, err
	}

	return &create_user.Input{
		Name:     req.Name,
		Document: doc,
	}, nil
}

func (h *CreateUserHandler) renderResponse(w http.ResponseWriter, userID string) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(createUserResponse{ID: userID}); err != nil {
		log.Printf("failed to encode post tweet response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func validateDoc(input string) (string, error) {
	doc := strings.TrimSpace(input)
	if doc == "" || !docRegex.MatchString(doc) {
		return "", ErrInvalidDoc
	}
	return doc, nil
}

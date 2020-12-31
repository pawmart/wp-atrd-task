package secrets

import (
		"context"
		"encoding/json"
		"errors"
		"github.com/go-kit/kit/endpoint"
		"github.com/gorilla/mux"
		"github.com/gorilla/schema"
		"net/http"
		"time"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
		NewSecret(ctx context.Context, payload newSecretRequest) (secretResponse, error)
		Fetch(ctx context.Context, hash string) (secretResponse, error)
}

type secretsService struct {
		repository SecretRepository
}

type newSecretRequest struct {
		Secret           string `json:"secret"`
		ExpireAfterViews int    `json:"expireAfterViews"`
		ExpireAfter      int    `json:"expireAfter"`
}

type fetchSecretRequest struct {
		Hash string `json:"hash"`
}

type secretResponse struct {
		Hash           string `json:"hash"`
		SecretText     string `json:"secretText"`
		CreatedAt      string `json:"createdAt"`
		ExpiresAt      string `json:"expiresAt"`
		RemainingViews int    `json:"remainingViews"`
}

func decodeNewSecretRequest(ctx context.Context, r *http.Request) (interface{}, error) {
		err := r.ParseForm()
		if err != nil {
				return nil, err
		}
		newSecret := &newSecretRequest{}
		decoder := schema.NewDecoder()

		err = decoder.Decode(newSecret, r.Form)
		if err != nil {
				return nil, err
		}

		return *newSecret, nil
}

func decodeFetchSecretRequest(ctx context.Context, r *http.Request) (interface{}, error) {
		vars := mux.Vars(r)
		hash, ok := vars["hash"]
		if !ok {
				return nil, errors.New("bad route")
		}
		return fetchSecretRequest{Hash: hash}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		return json.NewEncoder(w).Encode(response)
}

func NewService(repository SecretRepository) Service {
		return secretsService{repository: repository}
}

func (s secretsService) NewSecret(ctx context.Context, secretDTO newSecretRequest) (secretResponse, error) {
		sr := secretResponse{}
		if secretDTO.Secret == "" || secretDTO.ExpireAfterViews < 1 {
				return sr, ErrInvalidArgument
		}

		secret := NewSecret(
				secretDTO.Secret,
				secretDTO.ExpireAfterViews,
				secretDTO.ExpireAfter,
		)

		err := s.repository.Store(&secret)
		if err != nil {
				return sr, err
		}

		return secretResponse{
				Hash:           string(secret.Hash),
				SecretText:     secret.Secret,
				CreatedAt:      secret.Created.Format(time.UnixDate),
				ExpiresAt:      secret.Expire.Date.Format(time.UnixDate),
				RemainingViews: secret.Expire.ViewsLeft,
		}, nil
}

func (s secretsService) Fetch(ctx context.Context, hash string) (secretResponse, error) {
		var sr secretResponse
		secret, err := s.repository.Find(hash)

		if err != nil {
				return sr, err
		}

		return secretResponse{
				Hash:           hash,
				SecretText:     secret.Secret,
				CreatedAt:      secret.Created.Format(time.UnixDate),
				ExpiresAt:      secret.Expire.Date.Format(time.UnixDate),
				RemainingViews: secret.Expire.ViewsLeft,
		}, nil
}

func MakeNewSecretEndpoint(srv Service) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
				req := request.(newSecretRequest)
				ns, err := srv.NewSecret(ctx, req)
				if err != nil {
						return nil, err
				}

				return ns, nil
		}
}

func MakeFetchEndpoint(srv Service) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
				req := request.(fetchSecretRequest)
				s, err := srv.Fetch(ctx, req.Hash)
				if err != nil {
						return nil, err
				}

				return s, nil
		}
}

type Endpoints struct {
		NewSecretEndpoint endpoint.Endpoint
		FetchEndpoint     endpoint.Endpoint
}

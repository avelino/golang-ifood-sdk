package orders

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	httpadapter "github.com/arxdsilva/golang-ifood-sdk/adapters/http"
	"github.com/arxdsilva/golang-ifood-sdk/mocks"
	auth "github.com/arxdsilva/golang-ifood-sdk/services/authentication"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var orderDetails = `{
    "id": "REFERENCIA",
    "reference": "Id de referencia do pedido",
    "shortReference": "Extranet Id",
    "createdAt": "Timestamp do pedido",
    "type": "Tipo do pedido('DELIVERY' ou 'TOGO')",
    "merchant": {
        "id": "Identificador unico do restaurante⁎",
        "name": "Nome do restaurante",
        "phones": [
            "Telefone do restaurante"
        ],
        "address": {
            "formattedAddress": "Endereço formatado",
            "country": "Pais",
            "state": "Estado",
            "city": "Cidade",
            "neighborhood": "Bairro",
            "streetName": "Endereço (Tipo logradouro + Logradouro)",
            "streetNumber": "Numero",
            "postalCode": "CEP"
        }
    },
    "payments": [
        {
            "name": "Nome da forma de pagamento",
            "code": "Codigo da forma de pagamento⁎⁎⁎",
            "value": "Valor pago na forma",
            "prepaid": "Pedido pago ('true' ou 'false')",
            "issuer": "Bandeira"
        },
        {
            "name": "Nome da forma de pagamento",
            "code": "Codigo da forma de pagamento⁎⁎⁎",
            "value": "Valor pago na forma",
            "prepaid": "Pedido pago ('true' ou 'false')",
            "collector": "Recebedor da forma",
            "issuer": "Bandeira"
        }
    ],
    "customer": {
        "id": "Id do cliente",
        "uuid": "Id Único do cliente",
        "name": "Nome do cliente",
        "taxPayerIdentificationNumber": "CPF/CNPJ do cliente (opcional) ",
        "phone": "0800 + Localizador",
        "ordersCountOnRestaurant":"Qtde de pedidos do cliente nesse restaurante"
    },
    "items": [
        {
            "name": "Nome do item",
            "quantity": "Quantidade",
            "price": "Preço",
            "subItemsPrice": "Preço dos subitens",
            "totalPrice": "Preço total",
            "discount": "Desconto",
            "addition": "Adição",
            "externalCode": "Código do e-PDV",
            "subItems": [
                {
                    "name": "Nome do item",
                    "quantity": "Quantidade",
                    "price": "Preço",
                    "totalPrice": "Preço total",
                    "discount": "Desconto",
                    "addition": "Adição",
                    "externalCode": "Código do e-PDV"
                }
            ]
        },
        {
            "name": "Nome do item",
            "quantity": "Quantidade",
            "price": "Preço",
            "subItemsPrice": "Preço dos subitens",
            "totalPrice": "Preço total",
            "discount": "Desconto",
            "addition": "Adição",
            "subItems": [
                {
                    "name": "Nome do item",
                    "quantity": "Quantidade",
                    "price": "Preço",
                    "totalPrice": "Preço total",
                    "discount": "Desconto",
                    "addition": "Adição",
                    "externalCode": "Código e-PDV"
                }
            ]
        },
        {
            "name": "Nome do item",
            "quantity": "Quantidade",
            "price": "Preço",
            "subItemsPrice": "Preço dos subitens",
            "totalPrice": "Preço total",
            "discount": "Desconto",
            "addition": "Adição",
            "externalCode": "Código do e-PDV",
            "observations": "Observação do item"
        }
    ],
    "subTotal": "Total do pedido(Sem taxa de entrega)",
    "totalPrice": "Total do pedido(Com taxa de entrega)",
    "deliveryFee": "Taxa de entrega",
    "deliveryAddress": {
        "formattedAddress": "Endereço completo de entrega",
        "country": "Pais",
        "state": "Estado",
        "city": "Cidade",
        "coordinates": {
            "latitude": "Latitude do endereço",
            "longitude": "Longitude do endereço"
        },
        "neighborhood": "Bairro",
        "streetName": "Endereço(Tipo logradouro + Logradouro)",
        "streetNumber": "Numero",
        "postalCode": "CEP",
        "reference": "Referencia",
        "complement": "Complemento do endereço"
    },
    "deliveryDateTime": "Timestamp do pedido",
    "preparationTimeInSeconds": "Tempo de preparo do pedido em segundos"
}`

var trackingOK = `{
	"date": 0,
	"deliveryTime": "2020-06-29T15:24:30.405Z",
	"eta": 10,
	"etaToDestination": 0,
	"etaToOrigin": 0,
	"latitude": 0,
	"longitude": 0,
	"orderId": "string",
	"trackDate": "2020-06-29T15:24:30.406Z"
  }`

func TestGetDetails_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v3.0/orders/reference_id", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, orderDetails)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	od, err := ordersService.GetDetails("reference_id")
	assert.Nil(t, err)
	assert.Equal(t, "REFERENCIA", od.ID)
}

func TestGetDetails_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.GetDetails("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestGetDetails_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.GetDetails("reference")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestGetDetails_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v3.0/orders/reference_id", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	od, err := ordersService.GetDetails("reference_id")
	assert.NotNil(t, err)
	assert.Equal(t, OrderDetails{}, od)
	assert.Contains(t, err.Error(), "could not retrieve details")
}

func TestGetDetails_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.GetDetails("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetIntegrateStatus_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1.0/orders/reference_id/statuses/integration", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetIntegrateStatus("reference_id")
	assert.Nil(t, err)
}

func TestSetIntegrateStatus_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetIntegrateStatus("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestSetIntegrateStatus_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetIntegrateStatus("reference")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetIntegrateStatus_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1.0/orders/reference_id/statuses/integration", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetIntegrateStatus("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not be integrated")
}

func TestSetIntegrateStatus_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetIntegrateStatus("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetConfirmStatus_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1.0/orders/reference_id/statuses/confirmation", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetConfirmStatus("reference_id")
	assert.Nil(t, err)
}

func TestSetConfirmStatus_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetConfirmStatus("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestSetConfirmStatus_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetConfirmStatus("reference")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetConfirmStatus_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1.0/orders/reference_id/statuses/confirmation", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetConfirmStatus("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not be confirmed")
}

func TestSetConfirmStatus_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetConfirmStatus("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetDispatchStatus_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1.0/orders/reference_id/statuses/dispatch", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetDispatchStatus("reference_id")
	assert.Nil(t, err)
}

func TestSetDispatchStatus_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetDispatchStatus("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestSetDispatchStatus_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetDispatchStatus("reference")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetDispatchStatus_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1.0/orders/reference_id/statuses/dispatch", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetDispatchStatus("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not be dispatched")
}

func TestSetDispatchStatus_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetDispatchStatus("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetReadyToDeliverStatus_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/statuses/readyToDeliver", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetReadyToDeliverStatus("reference_id")
	assert.Nil(t, err)
}

func TestSetReadyToDeliverStatus_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetReadyToDeliverStatus("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestSetReadyToDeliverStatus_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetReadyToDeliverStatus("reference")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetReadyToDeliverStatus_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/statuses/readyToDeliver", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetReadyToDeliverStatus("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), " could not be set as 'ready to deliver'")
}

func TestSetReadyToDeliverStatus_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetReadyToDeliverStatus("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetCancelStatus_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v3.0/orders/reference_id/statuses/cancellationRequested", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetCancelStatus("reference_id", "501")
	assert.Nil(t, err)
}

func TestSetCancelStatus_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetCancelStatus("", "501")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestSetCancelStatus_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetCancelStatus("reference", "501")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestSetCancelStatus_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v3.0/orders/reference_id/statuses/cancellationRequested", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetCancelStatus("reference_id", "501")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), " could not be set as 'cancelled' code")
}

func TestSetCancelStatus_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.SetCancelStatus("reference_id", "501")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestClientCancellationStatus_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/statuses/consumerCancellationAccepted", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.ClientCancellationStatus("reference_id", true)
	assert.Nil(t, err)
}

func TestClientCancellationStatus_OK_NotAccepted(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/statuses/consumerCancellationDenied", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusAccepted)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.ClientCancellationStatus("reference_id", false)
	assert.Nil(t, err)
}

func TestClientCancellationStatus_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.ClientCancellationStatus("", true)
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestClientCancellationStatus_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.ClientCancellationStatus("reference", true)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestClientCancellationStatus_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/statuses/consumerCancellationAccepted", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodPost)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.ClientCancellationStatus("reference_id", true)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), " could not set 'client cancellation' status")
}

func TestClientCancellationStatus_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	err := ordersService.ClientCancellationStatus("reference_id", true)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestTracking_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/tracking", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprintf(w, trackingOK)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	tr, err := ordersService.Tracking("reference_id")
	assert.Nil(t, err)
	assert.Equal(t, 10, tr.Eta)
}

func TestTracking_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.Tracking("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestTracking_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.Tracking("reference")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestTracking_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/tracking", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.Tracking("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), " could not get tracking information")
}

func TestTracking_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.Tracking("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestDeliveryInformation_OK(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/delivery-information", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprintf(w, trackingOK)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	tr, err := ordersService.DeliveryInformation("reference_id")
	assert.Nil(t, err)
	assert.Equal(t, 10, tr.Eta)
}

func TestDeliveryInformation_NoRefereceID(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.DeliveryInformation("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func TestDeliveryInformation_ValidateErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(errors.New("some err"))
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, "ts.URL")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.DeliveryInformation("reference")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func TestDeliveryInformation_StatusBadRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v2.0/orders/reference_id/delivery-information", r.URL.Path)
			assert.Equal(t, "Bearer token", r.Header["Authorization"][0])
			assert.Equal(t, r.Method, http.MethodGet)
			w.WriteHeader(http.StatusBadRequest)
		}),
	)
	defer ts.Close()
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	adapter := httpadapter.New(http.DefaultClient, ts.URL)
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.DeliveryInformation("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could get delivery information")
}

func TestDeliveryInformation_DoReqErr(t *testing.T) {
	am := auth.AuthMock{}
	am.On("Validate").Once().Return(nil)
	am.On("GetToken").Once().Return("token")
	httpmock := &mocks.HttpClientMock{}
	httpmock.On("Do", mock.Anything).Once().Return(nil, errors.New("some err"))
	adapter := httpadapter.New(httpmock, "")
	ordersService := New(adapter, &am)
	assert.NotNil(t, ordersService)
	_, err := ordersService.DeliveryInformation("reference_id")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "some")
}

func Test_verifyCancel_OK(t *testing.T) {
	err := verifyCancel("reference", "501")
	assert.Nil(t, err)
}

func Test_verifyCancel_NoReferenceID(t *testing.T) {
	err := verifyCancel("", "501")
	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderReferenceNotSpecified, err)
}

func Test_verifyCancel_NoCode(t *testing.T) {
	err := verifyCancel("reference", "")
	assert.NotNil(t, err)
	assert.Equal(t, ErrCancelCodeNotSpecified, err)
}

func Test_verifyCancel_InvalidCode(t *testing.T) {
	err := verifyCancel("reference", "12344")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "is invalid, verify docs")
}

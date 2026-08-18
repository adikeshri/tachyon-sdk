import pytest
import requests
import responses

from .conftest import BASE_URL
from tachyon_sdk import (
    Tachyon,
    TachyonAuthenticationError,
    TachyonAuthorizationError,
    TachyonConflictError,
    TachyonConnectionError,
    TachyonNotFoundError,
    TachyonRequestError,
    TachyonServerError,
    TachyonTimeoutError,
)

STATUS_TO_ERROR = [
    (400, "invalid_query", TachyonRequestError),
    (401, "unauthorized", TachyonAuthenticationError),
    (403, "forbidden", TachyonAuthorizationError),
    (404, "collection_not_found", TachyonNotFoundError),
    (409, "collection_exists", TachyonConflictError),
    (500, "internal_error", TachyonServerError),
]


@pytest.mark.parametrize("status,code,error_cls", STATUS_TO_ERROR)
@responses.activate
def test_maps_http_status_to_error_class(client, status, code, error_cls):
    responses.add(
        responses.GET,
        f"{BASE_URL}/collections/products",
        json={"error": {"code": code, "message": f"boom: {code}"}},
        status=status,
    )

    with pytest.raises(error_cls) as excinfo:
        client.collections.retrieve("products")

    assert excinfo.value.code == code
    assert excinfo.value.status == status
    assert excinfo.value.message == f"boom: {code}"


@responses.activate
def test_wraps_network_failures_in_connection_error(client):
    responses.add(responses.GET, f"{BASE_URL}/health", body=requests.exceptions.ConnectionError("boom"))

    with pytest.raises(TachyonConnectionError):
        client.health()


@responses.activate
def test_raises_timeout_error_once_configured_timeout_elapses():
    responses.add(responses.GET, f"{BASE_URL}/health", body=requests.exceptions.Timeout("boom"))
    client = Tachyon(url=BASE_URL, timeout=0.01)

    with pytest.raises(TachyonTimeoutError):
        client.health()


@responses.activate
def test_falls_back_to_generic_error_code_when_body_is_not_the_documented_shape(client):
    responses.add(responses.GET, f"{BASE_URL}/health", body="upstream exploded", status=500)

    with pytest.raises(TachyonServerError) as excinfo:
        client.health()

    assert excinfo.value.code == "internal_error"
    assert excinfo.value.status == 500

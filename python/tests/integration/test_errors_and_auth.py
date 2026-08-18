import pytest

from .support import anonymous_client, search_only_client, unique_name
from tachyon_sdk import Tachyon, TachyonAuthenticationError, TachyonAuthorizationError, TachyonConnectionError

pytestmark = pytest.mark.integration


def test_no_api_key_is_rejected_when_auth_is_enabled():
    with pytest.raises(TachyonAuthenticationError) as excinfo:
        anonymous_client().collections.list()
    assert excinfo.value.code == "unauthorized"
    assert excinfo.value.status == 401


def test_search_only_key_can_read():
    result = search_only_client().collections.list()
    assert isinstance(result, list)


def test_search_only_key_cannot_write():
    with pytest.raises(TachyonAuthorizationError) as excinfo:
        search_only_client().collections.create({"name": unique_name("coll"), "fields": []})
    assert excinfo.value.code == "forbidden"
    assert excinfo.value.status == 403


def test_real_network_failure_raises_connection_error():
    # Nothing listens on port 1 (a reserved, unused TCP port), so this exercises
    # the real requests failure path rather than a mocked one.
    unreachable = Tachyon(url="http://127.0.0.1:1", timeout=2.0)
    with pytest.raises(TachyonConnectionError):
        unreachable.health()


def test_admin_key_works_end_to_end(client):
    result = client.collections.list()
    assert isinstance(result, list)

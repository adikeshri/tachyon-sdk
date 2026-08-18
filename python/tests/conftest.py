import pytest

from tachyon_sdk import Tachyon

BASE_URL = "http://localhost:8108"


@pytest.fixture
def client() -> Tachyon:
    return Tachyon(url=BASE_URL, api_key="admin-key")

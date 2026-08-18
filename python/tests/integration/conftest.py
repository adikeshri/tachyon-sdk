import pytest

from .support import admin_client
from tachyon_sdk import Tachyon


@pytest.fixture
def client() -> Tachyon:
    return admin_client()

import os
import uuid
from contextlib import contextmanager
from typing import Any, Dict, Iterator, Optional

from tachyon_sdk import Tachyon

BASE_URL = os.environ.get("TACHYON_URL", "http://localhost:8108")
ADMIN_KEY = os.environ.get("TACHYON_ADMIN_KEY", "admin-key")
SEARCH_KEY = os.environ.get("TACHYON_SEARCH_KEY", "search-key")


def admin_client() -> Tachyon:
    return Tachyon(url=BASE_URL, api_key=ADMIN_KEY)


def search_only_client() -> Tachyon:
    return Tachyon(url=BASE_URL, api_key=SEARCH_KEY)


def anonymous_client() -> Tachyon:
    return Tachyon(url=BASE_URL)


def unique_name(prefix: str) -> str:
    return f"{prefix}-{uuid.uuid4()}"


@contextmanager
def collection(client: Tachyon, schema: Dict[str, Any], name: Optional[str] = None) -> Iterator[str]:
    """
    Creates a collection for the duration of the `with` block, then deletes
    it — even if the block raises — so integration tests never leak
    collections into the shared server between runs.
    """
    resolved_name = name or unique_name("coll")
    client.collections.create({**schema, "name": resolved_name})
    try:
        yield resolved_name
    finally:
        try:
            client.collections.delete(resolved_name)
        except Exception:
            pass

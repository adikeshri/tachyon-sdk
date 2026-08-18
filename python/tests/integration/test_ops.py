import pytest

from .support import admin_client, anonymous_client

pytestmark = pytest.mark.integration


def test_health_does_not_require_an_api_key():
    health = anonymous_client().health()
    assert health["ok"] is True
    assert isinstance(health["version"], str)
    assert health["num_collections"] >= 0


def test_metrics_exposes_prometheus_text():
    metrics = admin_client().metrics()
    assert "tachyon_uptime_seconds" in metrics
    assert "tachyon_collections" in metrics

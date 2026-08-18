from .analytics import AnalyticsApi
from .client import Tachyon
from .collection import Collection
from .collections import CollectionsApi
from .documents import DocumentsApi
from .errors import (
    TachyonAuthenticationError,
    TachyonAuthorizationError,
    TachyonConflictError,
    TachyonConnectionError,
    TachyonError,
    TachyonNotFoundError,
    TachyonRequestError,
    TachyonServerError,
    TachyonTimeoutError,
)
from .types import (
    AnalyticsLatencyResponse,
    AnalyticsQueriesResponse,
    AnalyticsQuery,
    CollectionInfo,
    CollectionSchema,
    DocumentIndexResult,
    DocumentsIndexResponse,
    FieldSchema,
    FieldType,
    HealthResponse,
    MatchMode,
    SearchHit,
    SearchResponse,
    SuggestResponse,
    Suggestion,
    TachyonDocument,
    TypoToleranceConfig,
)

__version__ = "1.0.0"

__all__ = [
    "Tachyon",
    "Collection",
    "CollectionsApi",
    "DocumentsApi",
    "AnalyticsApi",
    "TachyonError",
    "TachyonRequestError",
    "TachyonAuthenticationError",
    "TachyonAuthorizationError",
    "TachyonNotFoundError",
    "TachyonConflictError",
    "TachyonServerError",
    "TachyonConnectionError",
    "TachyonTimeoutError",
    "FieldType",
    "FieldSchema",
    "TypoToleranceConfig",
    "CollectionSchema",
    "CollectionInfo",
    "TachyonDocument",
    "DocumentIndexResult",
    "DocumentsIndexResponse",
    "MatchMode",
    "SearchHit",
    "SearchResponse",
    "Suggestion",
    "SuggestResponse",
    "AnalyticsQuery",
    "AnalyticsQueriesResponse",
    "AnalyticsLatencyResponse",
    "HealthResponse",
]

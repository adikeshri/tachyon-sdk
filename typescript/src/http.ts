import { TachyonConnectionError, TachyonTimeoutError, errorFromResponse } from './errors.js';

export interface HttpClientOptions {
  /** Base URL of the Tachyon server, e.g. `http://localhost:8108`. */
  url: string;
  /** Sent as `X-TACHYON-API-KEY` on every request. */
  apiKey?: string;
  /** Per-request timeout in milliseconds. Default 15000. */
  timeoutMs?: number;
  /** Extra headers merged into every request. */
  headers?: Record<string, string>;
  /** Override the `fetch` implementation (mainly for testing). */
  fetch?: typeof fetch;
}

interface RequestOptions {
  query?: Record<string, unknown>;
  body?: unknown;
}

interface RawResponse {
  status: number;
  text: string;
}

/** Thin JSON-over-HTTP client shared by every resource in the SDK. */
export class HttpClient {
  private readonly baseUrl: string;
  private readonly apiKey: string | undefined;
  private readonly timeoutMs: number;
  private readonly extraHeaders: Record<string, string>;
  private readonly fetchImpl: typeof fetch;

  constructor(options: HttpClientOptions) {
    if (!options.url) {
      throw new Error('Tachyon client requires a `url`.');
    }
    this.baseUrl = options.url.replace(/\/+$/, '');
    this.apiKey = options.apiKey;
    this.timeoutMs = options.timeoutMs ?? 15_000;
    this.extraHeaders = options.headers ?? {};

    const fetchImpl = options.fetch ?? globalThis.fetch;
    if (!fetchImpl) {
      throw new Error(
        'No `fetch` implementation is available in this runtime. Upgrade to Node 18+ ' +
          'or pass one explicitly via the `fetch` client option.',
      );
    }
    this.fetchImpl = fetchImpl;
  }

  async request<T>(method: string, path: string, opts: RequestOptions = {}): Promise<T> {
    const { status, text } = await this.send(method, path, opts);
    if (status === 204 || text.length === 0) {
      return undefined as T;
    }

    let payload: unknown;
    try {
      payload = JSON.parse(text);
    } catch {
      if (status >= 400) {
        throw errorFromResponse(status, text, 'internal_error');
      }
      throw new TachyonConnectionError(`Tachyon returned a non-JSON response: ${truncate(text)}`);
    }

    if (status >= 400) {
      throw errorFromResponse(status, extractMessage(payload, text), extractCode(payload));
    }

    return payload as T;
  }

  async requestText(method: string, path: string, opts: RequestOptions = {}): Promise<string> {
    const { status, text } = await this.send(method, path, opts);
    if (status >= 400) {
      let payload: unknown;
      try {
        payload = JSON.parse(text);
      } catch {
        throw errorFromResponse(status, text || `HTTP ${status}`, 'internal_error');
      }
      throw errorFromResponse(status, extractMessage(payload, text), extractCode(payload));
    }
    return text;
  }

  private async send(method: string, path: string, opts: RequestOptions): Promise<RawResponse> {
    const url = this.buildUrl(path, opts.query);
    const headers: Record<string, string> = { Accept: 'application/json', ...this.extraHeaders };
    if (this.apiKey) {
      headers['X-TACHYON-API-KEY'] = this.apiKey;
    }

    let body: string | undefined;
    if (opts.body !== undefined) {
      headers['Content-Type'] = 'application/json';
      body = JSON.stringify(opts.body);
    }

    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), this.timeoutMs);
    try {
      const response = await this.fetchImpl(url, { method, headers, body, signal: controller.signal });
      const text = await response.text();
      return { status: response.status, text };
    } catch (err) {
      if (err instanceof Error && err.name === 'AbortError') {
        throw new TachyonTimeoutError(`Request to ${url} timed out after ${this.timeoutMs}ms`);
      }
      throw new TachyonConnectionError(
        `Failed to reach Tachyon at ${url}: ${err instanceof Error ? err.message : String(err)}`,
        err,
      );
    } finally {
      clearTimeout(timer);
    }
  }

  private buildUrl(path: string, query?: Record<string, unknown>): string {
    const url = new URL(this.baseUrl + path);
    if (query) {
      for (const [key, value] of Object.entries(query)) {
        if (value === undefined || value === null) continue;
        url.searchParams.set(key, String(value));
      }
    }
    return url.toString();
  }
}

function extractMessage(payload: unknown, fallback: string): string {
  if (payload && typeof payload === 'object' && 'error' in payload) {
    const error = (payload as { error?: unknown }).error;
    if (error && typeof error === 'object' && 'message' in error && typeof (error as { message?: unknown }).message === 'string') {
      return (error as { message: string }).message;
    }
  }
  return fallback;
}

function extractCode(payload: unknown): string {
  if (payload && typeof payload === 'object' && 'error' in payload) {
    const error = (payload as { error?: unknown }).error;
    if (error && typeof error === 'object' && 'code' in error && typeof (error as { code?: unknown }).code === 'string') {
      return (error as { code: string }).code;
    }
  }
  return 'internal_error';
}

function truncate(text: string, max = 200): string {
  return text.length > max ? `${text.slice(0, max)}…` : text;
}

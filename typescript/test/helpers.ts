export interface MockCall {
  url: string;
  method: string;
  headers: Record<string, string>;
  body?: string;
}

export interface MockResult {
  status: number;
  body?: unknown;
  text?: string;
}

/**
 * Builds a fake `fetch` that records every call and replies according to
 * `handler`, so tests never touch the network.
 */
export function createMockFetch(handler: (call: MockCall) => MockResult) {
  const calls: MockCall[] = [];

  const fetchImpl = (async (input: string | URL, init?: RequestInit) => {
    const url = typeof input === 'string' ? input : input.toString();
    const headers: Record<string, string> = {};
    if (init?.headers) {
      new Headers(init.headers).forEach((value, key) => {
        headers[key] = value;
      });
    }

    const call: MockCall = { url, method: init?.method ?? 'GET', headers, body: init?.body as string | undefined };
    calls.push(call);

    const result = handler(call);
    const text = result.text ?? (result.body !== undefined ? JSON.stringify(result.body) : '');
    const noBody = result.status === 204 || result.status === 205 || result.status === 304;
    return new Response(noBody ? null : text, {
      status: result.status,
      headers: { 'Content-Type': 'application/json' },
    });
  }) as typeof fetch;

  return { fetchImpl, calls };
}

import createClient from 'openapi-fetch';
import type { paths } from './generated';
import { notifySessionExpired } from './session-expiry';
import type { ApiError } from './types';

export class ApiRequestError extends Error {
  readonly code: string;
  readonly requestId: string;
  readonly fields: NonNullable<ApiError['error']['fields']>;

  constructor(error?: ApiError['error']) {
    super(error?.message ?? 'The request could not be completed.');
    this.name = 'ApiRequestError';
    this.code = error?.code ?? 'request_failed';
    this.requestId = error?.requestId ?? '';
    this.fields = error?.fields ?? [];
  }
}

export function apiFormErrors(error: unknown): { general?: string; fields: Record<string, string> } {
  if (!(error instanceof ApiRequestError)) return { general: error instanceof Error ? error.message : undefined, fields: {} };
  return { general: error.message, fields: Object.fromEntries(error.fields.map((field) => [field.field, field.message])) };
}

export const apiFetch: typeof fetch = async (input, init) => {
  const response = await fetch(input, init);
  const url = input instanceof Request ? input.url : input.toString();
  notifySessionExpired(response, new URL(url, window.location.origin).pathname);
  return response;
};

export const client = createClient<paths>({
  baseUrl: new URL('/api', window.location.origin).toString().replace(/\/$/, ''),
  credentials: 'same-origin',
  fetch: apiFetch,
});

export function requestError(error: unknown): ApiRequestError {
  const response = error as Partial<ApiError> | undefined;
  return new ApiRequestError(response?.error);
}

export async function unwrap<T>(request: Promise<{ data?: T; error?: unknown; response: Response }>): Promise<T> {
  const { data, error, response } = await request;
  if (!response.ok || error !== undefined) throw requestError(error);
  return data as T;
}

export class CsrfStore {
  #token?: string;

  async value(): Promise<string> {
    if (!this.#token) this.#token = (await unwrap(client.GET('/v1/auth/csrf'))).token;
    return this.#token!;
  }

  clear(): void { this.#token = undefined; }
}

export const csrf = new CsrfStore();
export const csrfHeaders = async () => ({ header: { 'X-CSRF-Token': await csrf.value() } });

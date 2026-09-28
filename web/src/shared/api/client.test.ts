import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

const response = (status: number, body?: unknown) => new Response(body === undefined ? undefined : JSON.stringify(body), {
  status,
  headers: body === undefined ? undefined : { 'content-type': 'application/json' },
});

describe('same-origin API client', () => {
  beforeEach(() => vi.resetModules());
  afterEach(() => vi.unstubAllGlobals());

  it('gets a CSRF token in memory before login, then renews it after login', async () => {
    const fetch = vi.fn()
      .mockResolvedValueOnce(response(200, { token: 'csrf-before-login' }))
      .mockResolvedValueOnce(response(204))
      .mockResolvedValueOnce(response(200, { token: 'csrf-after-login' }))
      .mockResolvedValueOnce(response(201, { id: 'card-1' }));
    vi.stubGlobal('fetch', fetch);
    const { api } = await import('@shared/api/client');

    await api.login('alex', 'correct horse battery staple');
    await api.createCard({ title: 'First card', description: '' });

    expect(fetch).toHaveBeenCalledTimes(4);
    const csrfRequest = fetch.mock.calls[0][0] as Request;
    const loginRequest = fetch.mock.calls[1][0] as Request;
    const createRequest = fetch.mock.calls[3][0] as Request;
    expect(csrfRequest.url).toBe('http://localhost/api/v1/auth/csrf');
    expect(loginRequest.url).toBe('http://localhost/api/v1/auth/login');
    expect(loginRequest.credentials).toBe('same-origin');
    expect(loginRequest.headers.get('X-CSRF-Token')).toBe('csrf-before-login');
    expect(loginRequest.headers.has('Cookie')).toBe(false);
    expect(createRequest.headers.get('X-CSRF-Token')).toBe('csrf-after-login');
  });

  it('sends complete versioned updates and a guarded delete with one CSRF token', async () => {
    const fetch = vi.fn()
      .mockResolvedValueOnce(response(200, { token: 'csrf-card' }))
      .mockResolvedValueOnce(response(200, { id: 'card-1', version: 2 }))
      .mockResolvedValueOnce(response(204));
    vi.stubGlobal('fetch', fetch);
    const { api } = await import('@shared/api/client');
    const change = { title: 'First card', description: 'A note', status: 'doing' as const, version: 1 };

    await expect(api.updateCard('card-1', change)).resolves.toMatchObject({ version: 2 });
    await api.deleteCard('card-1', 2);

    const updateRequest = fetch.mock.calls[1][0] as Request;
    const deleteRequest = fetch.mock.calls[2][0] as Request;
    expect(updateRequest.method).toBe('PATCH');
    expect(await updateRequest.json()).toEqual(change);
    expect(deleteRequest.url).toBe('http://localhost/api/v1/cards/card-1?version=2');
    expect(deleteRequest.method).toBe('DELETE');
    expect(deleteRequest.headers.get('X-CSRF-Token')).toBe('csrf-card');
    expect(deleteRequest.headers.get('Content-Type')).toBeNull();
  });

  it('preserves API general and field errors for accessible form mapping', async () => {
    const fetch = vi.fn()
      .mockResolvedValueOnce(response(200, { token: 'csrf-form' }))
      .mockResolvedValueOnce(response(422, { error: { code: 'validation_failed', message: 'Card details were rejected.', requestId: 'request-1', fields: [{ field: 'title', message: 'Title is required.' }] } }));
    vi.stubGlobal('fetch', fetch);
    const { ApiRequestError, api, apiFormErrors } = await import('@shared/api/client');

    const error = await api.createCard({ title: '', description: '' }).catch((reason: unknown) => reason);
    expect(error).toBeInstanceOf(ApiRequestError);
    expect(apiFormErrors(error)).toEqual({ general: 'Card details were rejected.', fields: { title: 'Title is required.' } });
  });

  it('announces a missing session cookie on protected requests without treating anonymous login as expiry', async () => {
    const fetch = vi.fn()
      .mockResolvedValueOnce(response(401, { error: { code: 'unauthorized', message: 'Authentication is required.', requestId: 'r-1' } }))
      .mockResolvedValueOnce(response(200, { token: 'csrf-login' }))
      .mockResolvedValueOnce(response(401, { error: { code: 'invalid_credentials', message: 'Invalid credentials.', requestId: 'r-2' } }))
      .mockResolvedValueOnce(response(401, { error: { code: 'unauthorized', message: 'Authentication is required.', requestId: 'r-3' } }));
    vi.stubGlobal('fetch', fetch);
    const listener = vi.fn();
    const { SESSION_EXPIRED_EVENT } = await import('@shared/api/session-expiry');
    window.addEventListener(SESSION_EXPIRED_EVENT, listener);
    const { api } = await import('@shared/api/client');

    await expect(api.getCurrentUser()).rejects.toMatchObject({ code: 'unauthorized' });
    await expect(api.login('alex', 'incorrect-password')).rejects.toMatchObject({ code: 'invalid_credentials' });
    expect(listener).not.toHaveBeenCalled();
    await expect(api.listCards()).rejects.toMatchObject({ code: 'unauthorized' });
    expect(listener).toHaveBeenCalledOnce();
    window.removeEventListener(SESSION_EXPIRED_EVENT, listener);
  });

  it('announces an expired browser session for centralized login redirection', async () => {
    const fetch = vi.fn().mockResolvedValueOnce(response(419, { error: { code: 'session_expired', message: 'Sign in again.', requestId: 'request-expired' } }));
    vi.stubGlobal('fetch', fetch);
    const listener = vi.fn();
    const { SESSION_EXPIRED_EVENT } = await import('@shared/api/session-expiry');
    window.addEventListener(SESSION_EXPIRED_EVENT, listener);
    const { api } = await import('@shared/api/client');

    await expect(api.listCards()).rejects.toMatchObject({ code: 'session_expired' });
    expect(listener).toHaveBeenCalledOnce();
    window.removeEventListener(SESSION_EXPIRED_EVENT, listener);
  });
});

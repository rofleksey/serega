import type { Card, CreateCardRequest, UpdateCardRequest, User } from './types';
import { client, csrf, csrfHeaders, unwrap } from './transport';

export * from './types';
export { ApiRequestError, apiFormErrors } from './transport';

export const api = {
  getCurrentUser: () => unwrap<User>(client.GET('/v1/auth/me')),
  login: async (username: string, password: string) => {
    await unwrap(client.POST('/v1/auth/login', { body: { username, password }, params: await csrfHeaders() }));
    csrf.clear();
  },
  logout: async () => {
    await unwrap(client.POST('/v1/auth/logout', { params: await csrfHeaders() }));
    csrf.clear();
  },
  listCards: (signal?: AbortSignal) => unwrap<{ cards: Card[] }>(client.GET('/v1/cards', { signal })),
  createCard: async (body: CreateCardRequest) =>
    unwrap<Card>(client.POST('/v1/cards', { body, params: await csrfHeaders() })),
  updateCard: async (cardId: string, body: UpdateCardRequest) =>
    unwrap<Card>(client.PATCH('/v1/cards/{cardId}', { body, params: { path: { cardId }, ...(await csrfHeaders()) } })),
  deleteCard: async (cardId: string, version: number) =>
    unwrap(client.DELETE('/v1/cards/{cardId}', { params: { path: { cardId }, query: { version }, ...(await csrfHeaders()) } })),
};

import type { components } from './generated';

export type ApiError = components['schemas']['ErrorResponse'];
export type User = components['schemas']['User'];
export type Card = components['schemas']['Card'];
export type CardStatus = Card['status'];
export type CreateCardRequest = components['schemas']['CreateCardRequest'];
export type UpdateCardRequest = components['schemas']['UpdateCardRequest'];

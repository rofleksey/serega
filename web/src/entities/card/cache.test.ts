import { QueryClient, QueryObserver } from '@tanstack/react-query';
import { describe, expect, it } from 'vitest';
import type { Card } from '@shared/api/client';
import { captureCardCache, reconcileSavedCard } from './cache';

const card: Card = { id: 'card-b', title: 'First card', description: '', status: 'todo', version: 1, createdBy: { id: 'user-1', username: 'alex' }, updatedBy: { id: 'user-1', username: 'alex' }, createdAt: '2026-09-28T10:00:00Z', updatedAt: '2026-09-28T10:00:00Z' };

describe('saved card reconciliation', () => {
  it('inserts missing cards in server creation order with IDs breaking timestamp ties', async () => {
    const client = new QueryClient();
    const newest = { ...card, id: 'newest', createdAt: '2026-09-29T10:00:00Z' };
    const oldest = { ...card, id: 'oldest', createdAt: '2026-09-27T10:00:00Z' };
    const tiedBefore = { ...card, id: 'card-a' };
    const tiedAfter = { ...card, id: 'card-c' };
    client.setQueryData(['cards'], { cards: [newest, tiedBefore, tiedAfter, oldest], extra: 'retained' });

    await reconcileSavedCard(client, card, captureCardCache(client));

    expect(client.getQueryData(['cards'])).toEqual({ cards: [newest, tiedBefore, card, tiedAfter, oldest], extra: 'retained' });
  });

  it('retains a newer cached version when an earlier mutation response arrives late', async () => {
    const client = new QueryClient();
    const newer = { ...card, title: 'Latest teammate update', status: 'done' as const, version: 3 };
    client.setQueryData(['cards'], { cards: [newer] });

    const effective = await reconcileSavedCard(client, { ...card, status: 'doing', version: 2 }, captureCardCache(client));

    expect(effective).toEqual(newer);
    expect(client.getQueryData(['cards'])).toEqual({ cards: [newer] });
  });

  it('ignores the result if session cleanup removes the query during refresh', async () => {
    const client = new QueryClient();
    let finishRefresh!: (result: { cards: Card[] }) => void;
    const observer = new QueryObserver(client, {
      queryKey: ['cards'],
      initialData: { cards: [card] },
      staleTime: Infinity,
      queryFn: () => new Promise<{ cards: Card[] }>((resolve) => { finishRefresh = resolve; }),
    });
    const unsubscribe = observer.subscribe(() => {});
    const result = reconcileSavedCard(client, { ...card, version: 2 }, captureCardCache(client));
    client.clear();
    finishRefresh({ cards: [{ ...card, version: 2 }] });

    expect(await result).toBeUndefined();
    expect(client.getQueryData(['cards'])).toBeUndefined();
    unsubscribe();
  });
});

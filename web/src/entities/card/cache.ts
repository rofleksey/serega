import type { QueryClient } from '@tanstack/react-query';
import type { api, Card } from '@shared/api/client';

type CardList = Awaited<ReturnType<typeof api.listCards>>;

export function captureCardCache(client: QueryClient) {
  return client.getQueryCache().find({ queryKey: ['cards'], exact: true });
}

export async function reconcileSavedCard(client: QueryClient, saved: Card, scope: ReturnType<typeof captureCardCache>) {
  // Session cleanup removes this query. A late mutation must not recreate it.
  if (!scope || captureCardCache(client) !== scope) return;
  client.setQueryData<CardList>(['cards'], (current) => {
    const cards = current?.cards ?? [];
    const existing = cards.find((card) => card.id === saved.id);
    if (existing && existing.version > saved.version) return current;
    // Match the API's creation order even if newer cards arrived while saving.
    const updated = existing
      ? cards.map((card) => card.id === saved.id ? saved : card)
      : [...cards, saved].sort((left, right) => Date.parse(right.createdAt) - Date.parse(left.createdAt) || (left.id < right.id ? -1 : left.id > right.id ? 1 : 0));
    return { ...current, cards: updated };
  });
  await client.invalidateQueries({ queryKey: ['cards'] });
  if (captureCardCache(client) !== scope) return;
  return client.getQueryData<CardList>(['cards'])?.cards.find((card) => card.id === saved.id) ?? saved;
}

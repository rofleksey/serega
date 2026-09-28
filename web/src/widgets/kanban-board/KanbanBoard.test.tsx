import { onlineManager, QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { KanbanBoard } from '@widgets/kanban-board/KanbanBoard';
import { ApiRequestError, type Card } from '@shared/api/client';

const api = vi.hoisted(() => ({ listCards: vi.fn(), updateCard: vi.fn(), deleteCard: vi.fn() }));
vi.mock('@shared/api/client', async (original) => ({ ...await original<typeof import('@shared/api/client')>(), api }));
const card: Card = { id: 'card-1', title: 'Write the guide', description: 'Make it useful.', status: 'todo', version: 1, createdBy: { id: 'user-1', username: 'alex' }, updatedBy: { id: 'user-1', username: 'alex' }, createdAt: '2026-09-28T10:00:00Z', updatedAt: '2026-09-28T10:00:00Z' };
function show(cachedCards?: Card[]) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  if (cachedCards) client.setQueryData(['cards'], { cards: cachedCards });
  render(<QueryClientProvider client={client}><KanbanBoard /></QueryClientProvider>);
  return client;
}
afterEach(() => { cleanup(); onlineManager.setOnline(true); });
beforeEach(() => { vi.resetAllMocks(); api.listCards.mockResolvedValue({ cards: [card] }); });

describe('shared kanban board', () => {
  it('renders each column with cards and attribution', async () => {
    show();
    const article = await screen.findByRole('article', { name: 'Write the guide' });
    expect(within(screen.getByRole('region', { name: 'To do' })).getByRole('article')).toBe(article);
    expect(screen.getByRole('region', { name: 'In progress' })).toBeInTheDocument();
    expect(screen.getByRole('region', { name: 'Done' })).toBeInTheDocument();
    expect(within(article).getByText(/Updated by alex/)).toBeInTheDocument();
  });

  it('moves a card with its version and refreshes the board', async () => {
    show();
    await screen.findByRole('article');
    api.updateCard.mockResolvedValue({ ...card, status: 'doing', version: 2 });
    api.listCards.mockResolvedValue({ cards: [{ ...card, status: 'doing', version: 2 }] });
    fireEvent.change(screen.getByLabelText('Status for Write the guide'), { target: { value: 'doing' } });
    await waitFor(() => expect(api.updateCard).toHaveBeenCalledWith('card-1', { title: card.title, description: card.description, status: 'doing', version: 1 }));
    await waitFor(() => expect(within(screen.getByRole('region', { name: 'In progress' })).getByRole('article')).toBeInTheDocument());
  });

  it('requires confirmation for deletion, supports cancel, and uses the confirmed version', async () => {
    show();
    await screen.findByRole('article');
    fireEvent.click(screen.getByRole('button', { name: 'Actions for Write the guide' }));
    fireEvent.click(screen.getByRole('menuitem', { name: 'Delete card' }));
    expect(screen.getByRole('dialog')).toHaveTextContent('Write the guide');
    expect(api.deleteCard).not.toHaveBeenCalled();
    fireEvent.click(screen.getByRole('button', { name: 'Cancel' }));
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
    fireEvent.click(screen.getByRole('button', { name: 'Actions for Write the guide' }));
    fireEvent.click(screen.getByRole('menuitem', { name: 'Delete card' }));
    api.deleteCard.mockResolvedValue(undefined);
    api.listCards.mockResolvedValue({ cards: [] });
    fireEvent.click(screen.getByRole('button', { name: 'Delete card' }));
    await waitFor(() => expect(api.deleteCard).toHaveBeenCalledWith('card-1', 1));
    expect(await screen.findByText('A fresh start')).toBeInTheDocument();
  });

  it('reports a conflicting move and refreshes without overwriting teammate changes', async () => {
    show();
    await screen.findByRole('article');
    api.updateCard.mockRejectedValue(new ApiRequestError({ code: 'card_conflict', message: 'Card changed.', requestId: 'r-1' }));
    api.listCards.mockResolvedValue({ cards: [{ ...card, title: 'Teammate title', version: 2 }] });
    fireEvent.change(screen.getByLabelText('Status for Write the guide'), { target: { value: 'doing' } });
    expect(await screen.findByText(/changed while you were moving/)).toBeInTheDocument();
    expect(await screen.findByRole('article', { name: 'Teammate title' })).toBeInTheDocument();
    expect(api.updateCard).toHaveBeenCalledTimes(1);
  });

  it('preserves an open edit draft when shared board updates arrive', async () => {
    const client = show();
    fireEvent.click(await screen.findByRole('button', { name: 'Edit Write the guide' }));
    fireEvent.change(screen.getByLabelText(/Title/), { target: { value: 'My unsaved draft' } });
    api.listCards.mockResolvedValue({ cards: [{ ...card, title: 'Teammate title', version: 2 }] });
    await act(async () => { await client.invalidateQueries({ queryKey: ['cards'] }); });
    expect(screen.getByLabelText(/Title/)).toHaveValue('My unsaved draft');
    expect(screen.getByLabelText(/Description/)).toHaveValue('Make it useful.');
  });

  it('labels paused offline updates while retaining the last shared board', async () => {
    onlineManager.setOnline(false);
    show([card]);
    expect(await screen.findByText('Offline')).toBeInTheDocument();
    expect(screen.getByRole('article', { name: 'Write the guide' })).toBeInTheDocument();
    expect(screen.queryByText('Board up to date')).not.toBeInTheDocument();
    expect(api.listCards).not.toHaveBeenCalled();
  });

  it('offers a retry when loading fails and recovers', async () => {
    api.listCards.mockRejectedValueOnce(new Error('Connection lost'));
    show();
    expect(await screen.findByText('Could not load the board')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Try again' }));
    expect(await screen.findByRole('article', { name: 'Write the guide' })).toBeInTheDocument();
  });
});

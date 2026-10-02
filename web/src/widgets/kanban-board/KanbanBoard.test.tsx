import { onlineManager, QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { KanbanBoard } from '@widgets/kanban-board/KanbanBoard';
import { ApiRequestError, type Card } from '@shared/api/client';

const api = vi.hoisted(() => ({ listCards: vi.fn(), createCard: vi.fn(), updateCard: vi.fn(), deleteCard: vi.fn() }));
vi.mock('@shared/api/client', async (original) => ({ ...await original<typeof import('@shared/api/client')>(), api }));
const card: Card = { id: 'card-1', title: 'Write the guide', description: 'Make it useful.', status: 'todo', version: 1, createdBy: { id: 'user-1', username: 'alex' }, updatedBy: { id: 'user-1', username: 'alex' }, createdAt: '2026-09-28T10:00:00Z', updatedAt: '2026-09-28T10:00:00Z' };
function show(cachedCards?: Card[]) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  if (cachedCards) client.setQueryData(['cards'], { cards: cachedCards });
  render(<QueryClientProvider client={client}><KanbanBoard /></QueryClientProvider>);
  return client;
}
function mobileViewport() {
  vi.stubGlobal('matchMedia', (query: string) => ({ matches: true, media: query, addEventListener: vi.fn(), removeEventListener: vi.fn(), addListener: vi.fn(), removeListener: vi.fn() }));
  vi.stubGlobal('scrollTo', vi.fn());
}
afterEach(() => { cleanup(); vi.unstubAllGlobals(); onlineManager.setOnline(true); });
beforeEach(() => { vi.resetAllMocks(); api.listCards.mockResolvedValue({ cards: [card] }); });

describe('shared kanban board', () => {
  it('renders each column with cards and attribution', async () => {
    show();
    const article = await screen.findByRole('article', { name: 'Write the guide' });
    expect(within(screen.getByRole('region', { name: 'To do' })).getByRole('article')).toBe(article);
    expect(screen.getByRole('region', { name: 'In progress' })).toBeInTheDocument();
    expect(screen.getByRole('region', { name: 'Done' })).toBeInTheDocument();
    expect(within(article).getByText(/^alex ·/)).toBeInTheDocument();
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
    expect(await screen.findByText('No cards in To do.')).toBeInTheDocument();
    expect(await screen.findByRole('button', { name: 'New card' })).toBeInTheDocument();
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
    expect(await screen.findByText(/Offline/)).toBeInTheDocument();
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

  it('exposes every mobile column through counted tabs while showing one column at a time', async () => {
    mobileViewport();
    api.listCards.mockResolvedValue({ cards: [card, { ...card, id: 'card-2', title: 'Build the app', status: 'doing' }, { ...card, id: 'card-3', title: 'Review the change', status: 'done' }] });
    show();

    expect(await screen.findByRole('tab', { name: 'To do (1)', selected: true })).toBeInTheDocument();
    expect(screen.getAllByRole('tabpanel')).toHaveLength(1);
    expect(screen.getByRole('article', { name: 'Write the guide' })).toBeInTheDocument();
    expect(screen.queryByRole('article', { name: 'Build the app' })).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('tab', { name: 'In progress (1)' }));
    expect(screen.getByRole('article', { name: 'Build the app' })).toBeInTheDocument();
    fireEvent.click(screen.getByRole('tab', { name: 'Done (1)' }));
    expect(screen.getByRole('article', { name: 'Review the change' })).toBeInTheDocument();
    expect(screen.queryByRole('article', { name: 'Write the guide' })).not.toBeInTheDocument();
  });

  it('follows a moved card to its mobile column after the versioned move succeeds', async () => {
    mobileViewport();
    show();
    await screen.findByRole('article');
    const moved: Card = { ...card, status: 'doing', version: 2 };
    api.updateCard.mockResolvedValue(moved);
    api.listCards.mockResolvedValue({ cards: [moved] });

    fireEvent.change(screen.getByRole('combobox', { name: 'Status for Write the guide' }), { target: { value: 'doing' } });
    expect(await screen.findByRole('tab', { name: 'In progress (1)', selected: true })).toBeInTheDocument();
    expect(within(screen.getByRole('tabpanel')).getByRole('article', { name: 'Write the guide' })).toBeInTheDocument();
    expect(api.updateCard).toHaveBeenCalledWith('card-1', { title: card.title, description: card.description, status: 'doing', version: 1 });
    expect(window.scrollTo).toHaveBeenCalledExactlyOnceWith({ top: 0, behavior: 'auto' });
  });

  it('keeps the mobile column on cancel and follows a newly created card back to To do', async () => {
    mobileViewport();
    show();
    fireEvent.click(await screen.findByRole('tab', { name: 'Done (0)' }));
    fireEvent.click(screen.getByRole('button', { name: 'New card' }));
    fireEvent.click(screen.getByRole('button', { name: 'Cancel' }));
    expect(screen.getByRole('tab', { name: 'Done (0)', selected: true })).toBeInTheDocument();
    expect(api.createCard).not.toHaveBeenCalled();

    fireEvent.click(screen.getByRole('button', { name: 'New card' }));
    fireEvent.change(screen.getByLabelText(/Title/), { target: { value: 'A new task' } });
    const created: Card = { ...card, id: 'card-2', title: 'A new task', description: '' };
    api.createCard.mockResolvedValue(created);
    api.listCards.mockResolvedValue({ cards: [card, created] });
    fireEvent.click(screen.getByRole('button', { name: 'Create card' }));

    expect(await screen.findByRole('tab', { name: 'To do (2)', selected: true })).toBeInTheDocument();
    expect(screen.getByRole('article', { name: 'A new task' })).toBeInTheDocument();
    expect(api.createCard).toHaveBeenCalledWith({ title: 'A new task', description: '' });
  });

  it('follows the saved status after editing a card on mobile', async () => {
    mobileViewport();
    show();
    fireEvent.click(await screen.findByRole('button', { name: 'Edit Write the guide' }));
    fireEvent.mouseDown(screen.getByRole('combobox', { name: 'Status' }));
    fireEvent.click(screen.getByRole('option', { name: 'Done' }));
    const saved: Card = { ...card, status: 'done', version: 2 };
    api.updateCard.mockResolvedValue(saved);
    api.listCards.mockResolvedValue({ cards: [saved] });
    fireEvent.click(screen.getByRole('button', { name: 'Save changes' }));

    expect(await screen.findByRole('tab', { name: 'Done (1)', selected: true })).toBeInTheDocument();
    expect(screen.getByRole('article', { name: 'Write the guide' })).toBeInTheDocument();
    expect(api.updateCard).toHaveBeenCalledWith('card-1', { title: card.title, description: card.description, status: 'done', version: 1 });
  });

  it('resets mobile scroll for a column change but leaves the selected column and scroll alone on refresh', async () => {
    mobileViewport();
    const client = show();
    fireEvent.click(await screen.findByRole('tab', { name: 'Done (0)' }));
    expect(window.scrollTo).toHaveBeenCalledExactlyOnceWith({ top: 0, behavior: 'auto' });

    api.listCards.mockResolvedValue({ cards: [card, { ...card, id: 'card-2', title: 'Someone else’s task' }] });
    await act(async () => { await client.invalidateQueries({ queryKey: ['cards'] }); });
    expect(screen.getByRole('tab', { name: 'Done (0)', selected: true })).toBeInTheDocument();
    expect(await screen.findByRole('tab', { name: 'To do (2)' })).toBeInTheDocument();
    expect(window.scrollTo).toHaveBeenCalledTimes(1);
  });

  it('retains a successfully moved card in its destination when the board refresh fails', async () => {
    mobileViewport();
    show();
    await screen.findByRole('article');
    api.updateCard.mockResolvedValue({ ...card, status: 'doing', version: 2 });
    api.listCards.mockRejectedValue(new Error('Refresh unavailable'));
    fireEvent.change(screen.getByRole('combobox', { name: 'Status for Write the guide' }), { target: { value: 'doing' } });

    expect(await screen.findByRole('tab', { name: 'In progress (1)', selected: true })).toBeInTheDocument();
    expect(screen.getByRole('article', { name: 'Write the guide' })).toBeInTheDocument();
    expect(screen.getByRole('combobox', { name: 'Status for Write the guide' })).toHaveValue('doing');
    expect(screen.getByText('Could not refresh the board')).toBeInTheDocument();
    expect(screen.getByText('Refresh unavailable')).toBeInTheDocument();
  });

  it('retains a successfully saved edit after moving to its destination when refresh fails', async () => {
    mobileViewport();
    show();
    fireEvent.click(await screen.findByRole('button', { name: 'Edit Write the guide' }));
    fireEvent.change(screen.getByLabelText(/Title/), { target: { value: 'Finished guide' } });
    fireEvent.mouseDown(screen.getByRole('combobox', { name: 'Status' }));
    fireEvent.click(screen.getByRole('option', { name: 'Done' }));
    api.updateCard.mockResolvedValue({ ...card, title: 'Finished guide', status: 'done', version: 2 });
    api.listCards.mockRejectedValue(new Error('Refresh unavailable'));
    fireEvent.click(screen.getByRole('button', { name: 'Save changes' }));

    expect(await screen.findByRole('tab', { name: 'Done (1)', selected: true })).toBeInTheDocument();
    expect(screen.getByRole('article', { name: 'Finished guide' })).toBeInTheDocument();
    expect(screen.getByText('Could not refresh the board')).toBeInTheDocument();
  });

  it('retains a newly created card and existing cards when the board refresh fails', async () => {
    mobileViewport();
    show();
    fireEvent.click(await screen.findByRole('tab', { name: 'Done (0)' }));
    fireEvent.click(screen.getByRole('button', { name: 'New card' }));
    fireEvent.change(screen.getByLabelText(/Title/), { target: { value: 'New task' } });
    api.createCard.mockResolvedValue({ ...card, id: 'card-2', title: 'New task', description: '', createdAt: '2026-10-02T10:00:00Z' });
    api.listCards.mockRejectedValue(new Error('Refresh unavailable'));
    fireEvent.click(screen.getByRole('button', { name: 'Create card' }));

    expect(await screen.findByRole('tab', { name: 'To do (2)', selected: true })).toBeInTheDocument();
    const cards = screen.getAllByRole('article');
    expect(cards[0]).toHaveAccessibleName('New task');
    expect(cards[1]).toHaveAccessibleName('Write the guide');
    expect(screen.getByText('Could not refresh the board')).toBeInTheDocument();
  });

  it.each(['move', 'edit'] as const)('follows the newer polled version after a delayed %s response and failed refresh', async (action) => {
    mobileViewport();
    const client = show();
    await screen.findByRole('article');
    let finishSave!: (saved: Card) => void;
    api.updateCard.mockImplementation(() => new Promise<Card>((resolve) => { finishSave = resolve; }));
    if (action === 'move') {
      fireEvent.change(screen.getByRole('combobox', { name: 'Status for Write the guide' }), { target: { value: 'doing' } });
    } else {
      fireEvent.click(screen.getByRole('button', { name: 'Edit Write the guide' }));
      fireEvent.mouseDown(screen.getByRole('combobox', { name: 'Status' }));
      fireEvent.click(screen.getByRole('option', { name: 'In progress' }));
      fireEvent.click(screen.getByRole('button', { name: 'Save changes' }));
    }
    await waitFor(() => expect(api.updateCard).toHaveBeenCalledOnce());
    const latest: Card = { ...card, title: 'Latest teammate update', status: 'done', version: 3 };
    api.listCards.mockResolvedValue({ cards: [latest] });
    await act(async () => { await client.invalidateQueries({ queryKey: ['cards'] }); });
    api.listCards.mockRejectedValue(new Error('Refresh unavailable'));
    await act(async () => { finishSave({ ...card, status: 'doing', version: 2 }); });

    expect(await screen.findByRole('tab', { name: 'Done (1)', selected: true })).toBeInTheDocument();
    expect(screen.getByRole('article', { name: 'Latest teammate update' })).toBeInTheDocument();
    expect(screen.getByRole('combobox', { name: 'Status for Latest teammate update' })).toHaveValue('done');
    expect(screen.getByText('Could not refresh the board')).toBeInTheDocument();
    expect(client.getQueryData(['cards'])).toEqual({ cards: [latest] });
  });
});

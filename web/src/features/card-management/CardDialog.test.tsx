import { onlineManager, QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { createQueryClient } from '@app/providers/query-client';
import { CardDialog } from '@features/card-management/CardDialog';
import { ApiRequestError, type Card } from '@shared/api/client';

const api = vi.hoisted(() => ({ createCard: vi.fn(), updateCard: vi.fn(), listCards: vi.fn() }));
vi.mock('@shared/api/client', async (original) => ({ ...await original<typeof import('@shared/api/client')>(), api }));
const card: Card = { id: 'card-1', title: 'First card', description: 'Original note', status: 'todo', version: 1, createdBy: { id: 'user-1', username: 'alex' }, updatedBy: { id: 'user-1', username: 'alex' }, createdAt: '2026-09-28T10:00:00Z', updatedAt: '2026-09-28T10:00:00Z' };
const close = vi.fn();

function show(existing?: Card, client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })) {
  render(<QueryClientProvider client={client}><CardDialog card={existing} onClose={close} /></QueryClientProvider>);
}
afterEach(() => { cleanup(); onlineManager.setOnline(true); });
beforeEach(() => { vi.resetAllMocks(); });

describe('card editing', () => {
  it('validates a required title and submits a trimmed new card', async () => {
    api.createCard.mockResolvedValue(card);
    show();
    fireEvent.change(screen.getByLabelText(/Title/), { target: { value: '   ' } });
    fireEvent.click(screen.getByRole('button', { name: 'Create card' }));
    expect(await screen.findByText('Title is required.')).toBeInTheDocument();
    expect(api.createCard).not.toHaveBeenCalled();
    fireEvent.change(screen.getByLabelText(/Title/), { target: { value: '  First card  ' } });
    fireEvent.change(screen.getByLabelText(/Description/), { target: { value: 'A note\nwith whitespace' } });
    fireEvent.click(screen.getByRole('button', { name: 'Create card' }));
    await waitFor(() => expect(api.createCard).toHaveBeenCalledWith({ title: 'First card', description: 'A note\nwith whitespace' }));
    await waitFor(() => expect(close).toHaveBeenCalledOnce());
  });

  it('reports offline writes immediately and keeps the draft editable instead of queuing it', async () => {
    onlineManager.setOnline(false);
    api.createCard.mockRejectedValue(new TypeError('Failed to fetch'));
    show(undefined, createQueryClient());
    fireEvent.change(screen.getByLabelText(/Title/), { target: { value: 'My offline draft' } });
    fireEvent.click(screen.getByRole('button', { name: 'Create card' }));
    expect(await screen.findByText('Failed to fetch')).toBeInTheDocument();
    expect(api.createCard).toHaveBeenCalledOnce();
    expect(screen.getByLabelText(/Title/)).toHaveValue('My offline draft');
    expect(screen.getByLabelText(/Title/)).toBeEnabled();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeEnabled();
    expect(screen.getByRole('button', { name: 'Create card' })).toBeEnabled();
    expect(close).not.toHaveBeenCalled();
  });

  it('keeps drafts after server failure and shows mapped field errors', async () => {
    api.updateCard.mockRejectedValue(new ApiRequestError({ code: 'validation_failed', message: 'Review the title.', requestId: 'r-1', fields: [{ field: 'title', message: 'Title rejected.' }] }));
    show(card);
    fireEvent.change(screen.getByLabelText(/Title/), { target: { value: 'My draft' } });
    fireEvent.click(screen.getByRole('button', { name: 'Save changes' }));
    expect(await screen.findByText('Title rejected.')).toBeInTheDocument();
    expect(screen.getByLabelText(/Title/)).toHaveValue('My draft');
    expect(close).not.toHaveBeenCalled();
  });

  it('retains the draft through a version conflict and requires review before saving again', async () => {
    api.updateCard.mockRejectedValueOnce(new ApiRequestError({ code: 'card_conflict', message: 'The card changed.', requestId: 'r-1' })).mockResolvedValue({ ...card, version: 3 });
    api.listCards.mockResolvedValue({ cards: [{ ...card, title: 'Teammate update', version: 2, updatedBy: { id: 'user-2', username: 'sam' } }] });
    show(card);
    fireEvent.change(screen.getByLabelText(/Title/), { target: { value: 'My draft' } });
    fireEvent.click(screen.getByRole('button', { name: 'Save changes' }));
    expect(await screen.findByText(/Someone changed this card/)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save changes' })).toBeDisabled();
    expect(screen.getByLabelText(/Title/)).toHaveValue('My draft');
    fireEvent.click(screen.getByRole('button', { name: 'Review latest' }));
    expect(await screen.findByText('Teammate update')).toBeInTheDocument();
    expect(screen.getByLabelText(/Title/)).toHaveValue('My draft');
    fireEvent.click(screen.getByRole('button', { name: 'Save changes' }));
    await waitFor(() => expect(api.updateCard).toHaveBeenLastCalledWith('card-1', { title: 'My draft', description: 'Original note', status: 'todo', version: 2 }));
  });

  it('keeps a deleted card draft available without allowing a stale save', async () => {
    api.updateCard.mockRejectedValue(new ApiRequestError({ code: 'not_found', message: 'Card not found.', requestId: 'r-1' }));
    show(card);
    fireEvent.change(screen.getByLabelText(/Title/), { target: { value: 'My draft' } });
    fireEvent.click(screen.getByRole('button', { name: 'Save changes' }));
    expect(await screen.findByText(/This card was deleted/)).toBeInTheDocument();
    expect(screen.getByLabelText(/Title/)).toHaveValue('My draft');
    expect(screen.getByRole('button', { name: 'Save changes' })).toBeDisabled();
  });
});

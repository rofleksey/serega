import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { UserProvider } from '@entities/session/auth';
import { AppLayout } from '@app/layout/AppLayout';
import { LoginPage } from '@pages/login/LoginPage';
import { ApiRequestError } from '@shared/api/client';
import { SESSION_EXPIRED_EVENT } from '@shared/api/session-expiry';

const api = vi.hoisted(() => ({ getCurrentUser: vi.fn(), login: vi.fn(), logout: vi.fn() }));
vi.mock('@shared/api/client', async (original) => ({ ...await original<typeof import('@shared/api/client')>(), api }));
afterEach(() => { cleanup(); vi.resetAllMocks(); });

function providers(children: React.ReactNode, entries = ['/'], client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })) {
  return <QueryClientProvider client={client}><MemoryRouter initialEntries={entries}><UserProvider>{children}</UserProvider></MemoryRouter></QueryClientProvider>;
}

function routes() {
  return <Routes><Route path="login" element={<LoginPage />} /><Route element={<AppLayout />}><Route index element={<div>Board content</div>} /></Route></Routes>;
}

describe('responsive shell and login', () => {
  it('identifies the shared board and signed-in user', async () => {
    api.getCurrentUser.mockResolvedValue({ id: 'user-1', username: 'alex' });
    render(providers(routes()));
    expect(await screen.findByText('alex')).toBeInTheDocument();
    expect(screen.getByText('Shared board')).toBeInTheDocument();
    expect(screen.getByRole('main')).toHaveTextContent('Board content');
  });

  it('shows expired-session feedback and toggles password visibility', () => {
    api.getCurrentUser.mockRejectedValue(new Error('unauthorized'));
    render(providers(routes(), [{ pathname: '/login', state: { sessionExpired: true } }] as never));
    expect(screen.getByRole('heading', { name: 'Serega', level: 1 })).toBeInTheDocument();
    expect(screen.getByText(/session expired/i)).toBeInTheDocument();
    const password = screen.getByLabelText(/Password/);
    expect(password).toHaveAttribute('type', 'password');
    fireEvent.click(screen.getByRole('button', { name: 'Show password' }));
    expect(password).toHaveAttribute('type', 'text');
  });

  it('logs in and loads the current user before entering the board', async () => {
    api.getCurrentUser.mockRejectedValueOnce(new Error('unauthorized')).mockResolvedValue({ id: 'user-1', username: 'alex' });
    api.login.mockResolvedValue(undefined);
    render(providers(routes(), ['/login']));
    await waitFor(() => expect(api.getCurrentUser).toHaveBeenCalledOnce());
    fireEvent.change(screen.getByLabelText(/Username/), { target: { value: 'alex' } });
    fireEvent.change(screen.getByLabelText(/Password/), { target: { value: 'test-password' } });
    fireEvent.click(screen.getByRole('button', { name: 'Sign in' }));
    expect(await screen.findByText('Board content')).toBeInTheDocument();
    expect(api.login).toHaveBeenCalledWith('alex', 'test-password');
  });

  it('shows a failed session refresh after login instead of entering an unauthenticated board', async () => {
    api.getCurrentUser.mockRejectedValue(new Error('Session could not be loaded.'));
    api.login.mockResolvedValue(undefined);
    render(providers(routes(), ['/login']));
    await waitFor(() => expect(api.getCurrentUser).toHaveBeenCalledOnce());
    fireEvent.change(screen.getByLabelText(/Username/), { target: { value: 'alex' } });
    fireEvent.change(screen.getByLabelText(/Password/), { target: { value: 'test-password' } });
    fireEvent.click(screen.getByRole('button', { name: 'Sign in' }));
    expect(await screen.findByText('Session could not be loaded.')).toBeInTheDocument();
    expect(screen.queryByText('Board content')).not.toBeInTheDocument();
  });

  it('shows login failures without discarding the username', async () => {
    api.getCurrentUser.mockRejectedValue(new Error('unauthorized'));
    api.login.mockRejectedValue(new Error('Invalid username or password.'));
    render(providers(routes(), ['/login']));
    fireEvent.change(screen.getByLabelText(/Username/), { target: { value: 'alex' } });
    fireEvent.change(screen.getByLabelText(/Password/), { target: { value: 'incorrect-password' } });
    fireEvent.click(screen.getByRole('button', { name: 'Sign in' }));
    expect(await screen.findByText('Invalid username or password.')).toBeInTheDocument();
    expect(screen.getByLabelText(/Username/)).toHaveValue('alex');
  });

  it('clears all retained board data on logout', async () => {
    api.getCurrentUser.mockResolvedValue({ id: 'user-1', username: 'alex' });
    api.logout.mockResolvedValue(undefined);
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    client.setQueryData(['cards'], { cards: [{ id: 'card-1' }] });
    render(providers(routes(), ['/'], client));
    await screen.findByText('alex');
    fireEvent.click(screen.getByRole('button', { name: 'Sign out' }));
    expect(await screen.findByRole('button', { name: 'Sign in' })).toBeInTheDocument();
    expect(api.logout).toHaveBeenCalledOnce();
    expect(client.getQueryData(['cards'])).toBeUndefined();
  });

  it('keeps the board and reports a failed sign-out', async () => {
    api.getCurrentUser.mockResolvedValue({ id: 'user-1', username: 'alex' });
    api.logout.mockRejectedValue(new Error('Logout unavailable'));
    render(providers(routes()));
    fireEvent.click(await screen.findByRole('button', { name: 'Sign out' }));
    expect(await screen.findByText('Sign out failed')).toBeInTheDocument();
    expect(screen.getByText('Logout unavailable')).toBeInTheDocument();
    expect(screen.getByText('Board content')).toBeInTheDocument();
  });

  it('discards a previously loaded user when session refresh finds a removed cookie', async () => {
    api.getCurrentUser.mockResolvedValue({ id: 'user-1', username: 'alex' });
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    client.setQueryData(['cards'], { cards: [{ id: 'card-1' }] });
    render(providers(routes(), ['/'], client));
    await screen.findByText('alex');
    api.getCurrentUser.mockRejectedValue(new ApiRequestError({ code: 'unauthorized', message: 'Authentication is required.', requestId: 'r-1' }));
    await act(async () => { await client.invalidateQueries({ queryKey: ['current-user'] }); });
    expect(await screen.findByText(/session expired/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Sign in' })).toBeInTheDocument();
    expect(client.getQueryData(['cards'])).toBeUndefined();
  });

  it('redirects to login and clears the cache on session expiry', async () => {
    api.getCurrentUser.mockResolvedValue({ id: 'user-1', username: 'alex' });
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    client.setQueryData(['cards'], { cards: [{ id: 'card-1' }] });
    render(providers(routes(), ['/'], client));
    await screen.findByText('alex');
    act(() => window.dispatchEvent(new Event(SESSION_EXPIRED_EVENT)));
    expect(await screen.findByText(/session expired/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Sign in' })).toBeInTheDocument();
    expect(client.getQueryData(['cards'])).toBeUndefined();
  });
});

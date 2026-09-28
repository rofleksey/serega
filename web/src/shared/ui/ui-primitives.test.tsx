import { useState } from 'react';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { Button } from '@mui/material';
import { EmptyState, LoadingState, RecoverableError } from '@shared/ui/async-states';
import { ConfirmationDialog } from '@shared/ui/confirmation-dialog';
import { FeedbackProvider, useFeedback } from '@app/providers/feedback';
import { OverflowActions } from '@shared/ui/overflow-actions';
import { PageHeader } from '@shared/ui/page-header';
import { StatusIndicator } from '@shared/ui/status-indicator';

afterEach(() => cleanup());

describe('shared UI primitives', () => {
  it('renders page hierarchy and text-backed status', () => {
    render(<MemoryRouter><PageHeader title="Board" description="Shared tasks" crumbs={[{ label: 'Home', to: '/' }, { label: 'Board' }]} /></MemoryRouter>);
    expect(screen.getByRole('heading', { name: 'Board' })).toBeInTheDocument();
    expect(screen.getByRole('navigation', { name: 'Breadcrumb' })).toBeInTheDocument();
    render(<StatusIndicator label="Online" tone="success" detail="Inventory current" />);
    expect(screen.getByText('Online')).toBeInTheDocument();
    expect(screen.getByText('Inventory current')).toBeInTheDocument();
  });

  it('exposes loading, empty, and retryable error states accessibly', () => {
    const retry = vi.fn();
    const { rerender } = render(<LoadingState label="Loading board" />);
    expect(screen.getByRole('status', { name: 'Loading board' })).toBeInTheDocument();
    rerender(<EmptyState title="No cards" description="Add a card to begin." action={<Button>Add card</Button>} />);
    expect(screen.getByText('No cards')).toBeInTheDocument();
    rerender(<RecoverableError error={new Error('Connection lost')} onRetry={retry} />);
    fireEvent.click(screen.getByRole('button', { name: 'Try again' }));
    expect(retry).toHaveBeenCalledOnce();
  });

  it('keeps confirmation targets explicit and cancelable', () => {
    const confirm = vi.fn();
    const cancel = vi.fn();
    render(<ConfirmationDialog open title="Delete card?" target="Write the guide" effect="Removes this card for everyone." confirmLabel="Delete card" safetyNote="This action cannot be undone." onCancel={cancel} onConfirm={confirm} />);
    expect(screen.getByText(/Write the guide/)).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Delete card' }));
    expect(confirm).toHaveBeenCalledOnce();
    fireEvent.click(screen.getByRole('button', { name: 'Cancel' }));
    expect(cancel).toHaveBeenCalledOnce();
  });

  it('labels overflow actions and invokes the selected action', () => {
    const select = vi.fn();
    render(<OverflowActions label="Card actions" actions={[{ label: 'Edit card', onSelect: select }]} />);
    fireEvent.click(screen.getByRole('button', { name: 'Card actions' }));
    fireEvent.click(screen.getByRole('menuitem', { name: 'Edit card' }));
    expect(select).toHaveBeenCalledOnce();
  });

  it('announces transient success and error feedback', () => {
    function Harness() {
      const feedback = useFeedback();
      return <><Button onClick={() => feedback.showSuccess('Saved')}>Success</Button><Button onClick={() => feedback.showError('Failed')}>Error</Button></>;
    }
    render(<FeedbackProvider><Harness /></FeedbackProvider>);
    fireEvent.click(screen.getByRole('button', { name: 'Success' }));
    expect(screen.getByText('Saved')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Error' }));
    expect(screen.getByText('Failed')).toBeInTheDocument();
  });

  it('allows callers to control an open confirmation dialog', async () => {
    function Harness() {
      const [open, setOpen] = useState(true);
      return <ConfirmationDialog open={open} title="Remove?" target="card" effect="Removes the card." confirmLabel="Remove" onCancel={() => setOpen(false)} onConfirm={() => setOpen(false)} />;
    }
    render(<Harness />);
    fireEvent.click(screen.getByRole('button', { name: 'Cancel' }));
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
  });
});

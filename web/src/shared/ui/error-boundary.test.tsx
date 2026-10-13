import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { ErrorBoundary } from './error-boundary';

afterEach(() => {
  cleanup();
  vi.restoreAllMocks();
});
function Broken({ fail }: { fail: boolean }) {
  if (fail) throw new Error('private diagnostic that must not be displayed');
  return <p>Recovered section</p>;
}
it('contains rendering failures without replacing neighboring content or exposing errors', () => {
  vi.spyOn(console, 'error').mockImplementation(() => {});
  const view = render(
    <>
      <nav>Navigation still works</nav>
      <ErrorBoundary resetKey="one">
        <Broken fail />
      </ErrorBoundary>
    </>,
  );
  expect(screen.getByRole('alert')).toHaveTextContent('Something went wrong');
  expect(screen.getByRole('navigation')).toHaveTextContent('Navigation still works');
  expect(screen.queryByText(/private diagnostic/)).not.toBeInTheDocument();
  view.rerender(
    <>
      <nav>Navigation still works</nav>
      <ErrorBoundary resetKey="one">
        <Broken fail={false} />
      </ErrorBoundary>
    </>,
  );
  fireEvent.click(screen.getByRole('button', { name: 'Try again' }));
  expect(screen.getByText('Recovered section')).toBeInTheDocument();
});
it('recovers automatically after navigation to another route', () => {
  vi.spyOn(console, 'error').mockImplementation(() => {});
  const view = render(
    <ErrorBoundary resetKey="one">
      <Broken fail />
    </ErrorBoundary>,
  );
  view.rerender(
    <ErrorBoundary resetKey="two">
      <Broken fail={false} />
    </ErrorBoundary>,
  );
  expect(screen.getByText('Recovered section')).toBeInTheDocument();
});

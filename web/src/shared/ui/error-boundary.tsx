import { Component, type ReactNode } from 'react';
import { Box, Button, Stack, Typography } from '@mui/material';

type Props = { children: ReactNode; section?: string; resetKey?: string; onRetry?: () => void };
type State = { failed: boolean };

export function UnexpectedError({ section, onRetry }: { section?: string; onRetry?: () => void }) {
  return (
    <Box role="alert" sx={{ p: 3, maxWidth: 560, mx: 'auto' }}>
      <Stack gap={1.5}>
        <Typography component="h1" variant="h6">
          Something went wrong
        </Typography>
        <Typography color="text.secondary">
          {section ? `${section} could not be displayed.` : 'This page could not be displayed.'}
        </Typography>
        <Stack direction="row" gap={1}>
          {onRetry && (
            <Button variant="contained" onClick={onRetry}>
              Try again
            </Button>
          )}
          <Button onClick={() => window.location.reload()}>Reload page</Button>
          <Button component="a" href="/">
            Home
          </Button>
        </Stack>
      </Stack>
    </Box>
  );
}

// Keep unexpected render failures away from unrelated sections and navigation.
// Request failures still use their query/form recovery paths; never show raw errors.
export class ErrorBoundary extends Component<Props, State> {
  state: State = { failed: false };

  static getDerivedStateFromError(): State {
    return { failed: true };
  }

  componentDidUpdate(previous: Props) {
    if (this.state.failed && previous.resetKey !== this.props.resetKey) {
      this.setState({ failed: false });
    }
  }

  render() {
    if (this.state.failed)
      return (
        <UnexpectedError
          section={this.props.section}
          onRetry={() => {
            this.props.onRetry?.();
            this.setState({ failed: false });
          }}
        />
      );
    return this.props.children;
  }
}

import { createTheme } from '@mui/material/styles';

export const theme = createTheme({
  palette: {
    mode: 'dark',
    background: { default: '#090b0f', paper: '#11151b' },
    primary: { main: '#6ea8fe' },
    success: { main: '#66c98b' },
    warning: { main: '#e8b85b' },
    error: { main: '#ef7d7d' },
    divider: 'rgba(255, 255, 255, 0.09)',
  },
  shape: { borderRadius: 10 },
  typography: {
    fontFamily: 'Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
    h4: { fontSize: 'clamp(1.65rem, 3vw, 2rem)', lineHeight: 1.18, fontWeight: 750 },
    h5: { fontWeight: 750 },
    h6: { fontWeight: 700 },
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        ':root': { colorScheme: 'dark' },
        body: { minWidth: 0, overflowX: 'hidden' },
        '::selection': { backgroundColor: 'rgba(110, 168, 254, 0.32)' },
        '@media (prefers-reduced-motion: reduce)': {
          '*, *::before, *::after': { scrollBehavior: 'auto !important', transitionDuration: '0.01ms !important', animationDuration: '0.01ms !important', animationIterationCount: '1 !important' },
        },
      },
    },
    MuiButton: { styleOverrides: { root: { minHeight: 44, textTransform: 'none', fontWeight: 700 } } },
    MuiIconButton: { styleOverrides: { root: { minWidth: 44, minHeight: 44 } } },
    MuiLink: { defaultProps: { underline: 'hover' } },
    MuiPaper: { styleOverrides: { root: { backgroundImage: 'none' }, outlined: { borderColor: 'rgba(255, 255, 255, 0.09)' } } },
    MuiCard: { styleOverrides: { root: { backgroundImage: 'none' } } },
    MuiDialog: { styleOverrides: { paper: { backgroundImage: 'none' } } },
    MuiDialogActions: { styleOverrides: { root: { flexWrap: 'wrap', gap: 8 } } },
    MuiSnackbar: { styleOverrides: { root: { bottom: 'calc(16px + env(safe-area-inset-bottom))' } } },
  },
});

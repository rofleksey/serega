import { createTheme } from '@mui/material/styles';

// Compact pointer controls must not shrink the touch targets on phones or tablets.
const touchControls = '@media (pointer: coarse), (max-width: 899.95px)';

export const theme = createTheme({
  palette: {
    mode: 'dark',
    background: { default: '#101215', paper: '#171a20' },
    primary: { main: '#8ab4f8' },
    text: { primary: '#e4e8ef', secondary: '#a3adbb' },
    success: { main: '#66c98b' },
    warning: { main: '#e8b85b' },
    error: { main: '#ef7d7d' },
    divider: 'rgba(255, 255, 255, 0.09)',
  },
  shape: { borderRadius: 6 },
  typography: {
    fontFamily: 'Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
    h4: { fontSize: '1.375rem', lineHeight: 1.3, fontWeight: 650 },
    h5: { fontSize: '1.0625rem', lineHeight: 1.4, fontWeight: 600 },
    h6: { fontSize: '1rem', lineHeight: 1.4, fontWeight: 600 },
    body1: { fontSize: '.875rem', lineHeight: 1.5 },
    body2: { fontSize: '.8125rem', lineHeight: 1.5 },
    button: { fontSize: '.8125rem', fontWeight: 600, textTransform: 'none' },
    caption: { fontSize: '.75rem', lineHeight: 1.5 },
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        ':root': { colorScheme: 'dark' },
        body: { minWidth: 0, overflowX: 'hidden' },
        'button, a, input, select, textarea': { WebkitTapHighlightColor: 'transparent' },
        '::selection': { backgroundColor: 'rgba(110, 168, 254, 0.32)' },
        '@media (prefers-reduced-motion: reduce)': {
          '*, *::before, *::after': { scrollBehavior: 'auto !important', transitionDuration: '0.01ms !important', animationDuration: '0.01ms !important', animationIterationCount: '1 !important' },
        },
      },
    },
    MuiButton: { defaultProps: { size: 'small', disableElevation: true }, styleOverrides: { root: { minHeight: 32, textTransform: 'none', fontWeight: 600, [touchControls]: { minHeight: 44 } } } },
    MuiIconButton: { defaultProps: { size: 'small' }, styleOverrides: { root: { minWidth: 32, minHeight: 32, color: '#a3adbb', '& .MuiSvgIcon-root': { fontSize: 20 }, [touchControls]: { minWidth: 44, minHeight: 44 } } } },
    MuiTextField: { defaultProps: { size: 'small' } },
    MuiFormControl: { defaultProps: { size: 'small' } },
    MuiInputBase: { styleOverrides: { input: { [touchControls]: { fontSize: 16 } } } },
    MuiOutlinedInput: { styleOverrides: { root: { [touchControls]: { minHeight: 44 } } } },
    MuiFormControlLabel: { styleOverrides: { root: { [touchControls]: { minHeight: 44 } } } },
    MuiTabs: { styleOverrides: { root: { minHeight: 36, [touchControls]: { minHeight: 44 } }, indicator: { height: 2 } } },
    MuiTab: { styleOverrides: { root: { minHeight: 36, minWidth: 0, padding: '8px 12px', textTransform: 'none', fontWeight: 600, [touchControls]: { minHeight: 44 } } } },
    MuiMenuItem: { styleOverrides: { root: { minHeight: 36, [touchControls]: { minHeight: 44 } } } },
    MuiChip: { styleOverrides: { sizeSmall: { height: 22, fontSize: '.75rem' } } },
    MuiLink: { defaultProps: { underline: 'hover' } },
    MuiPaper: { styleOverrides: { root: { backgroundImage: 'none' }, outlined: { borderColor: 'rgba(255, 255, 255, 0.09)' } } },
    MuiCard: { styleOverrides: { root: { backgroundImage: 'none' } } },
    MuiDialog: { styleOverrides: { paper: { backgroundImage: 'none', '@media (max-width: 599.95px)': { margin: 12, maxWidth: 'calc(100% - 24px)' } }, paperFullScreen: { margin: 0, maxWidth: '100%', '@media (max-width: 599.95px)': { margin: 0, maxWidth: '100%' } } } },
    MuiDialogTitle: { styleOverrides: { root: { padding: '16px 20px', fontSize: '1.0625rem', overflowWrap: 'anywhere' } } },
    MuiDialogContent: { styleOverrides: { root: { padding: '12px 20px 20px', overflowWrap: 'anywhere' } } },
    MuiDialogActions: { styleOverrides: { root: { flexWrap: 'wrap', gap: 8 } } },
    MuiSnackbar: { styleOverrides: { root: { bottom: 'calc(16px + env(safe-area-inset-bottom))' } } },
  },
});

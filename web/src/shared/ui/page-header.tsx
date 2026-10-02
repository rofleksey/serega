import { Link as RouterLink } from 'react-router-dom';
import { Breadcrumbs, Link, Stack, Typography } from '@mui/material';

export type PageCrumb = { label: string; to?: string };

export function PageHeader({ title, description, crumbs = [], actions }: { title: string; description?: string; crumbs?: PageCrumb[]; actions?: React.ReactNode }) {
  return <Stack spacing={.75} sx={{ minWidth: 0 }}>
      {crumbs.length > 0 && <Breadcrumbs aria-label="Breadcrumb" sx={{ color: 'text.secondary', fontSize: '.75rem', '& .MuiBreadcrumbs-ol': { flexWrap: 'wrap' }, '& .MuiBreadcrumbs-li': { minWidth: 0, maxWidth: '100%' } }}>
        {crumbs.map((crumb, index) => crumb.to
          ? <Link key={`${crumb.label}-${index}`} component={RouterLink} to={crumb.to} color="inherit" sx={{ overflowWrap: 'anywhere' }}>{crumb.label}</Link>
          : <Typography key={`${crumb.label}-${index}`} color="text.primary" variant="caption" sx={{ overflowWrap: 'anywhere' }}>{crumb.label}</Typography>)}
      </Breadcrumbs>}
    <Stack direction="row" spacing={1} useFlexGap flexWrap="wrap" alignItems="center" justifyContent="space-between">
      <Typography component="h1" variant="h4" sx={{ minWidth: 0, overflowWrap: 'anywhere' }}>{title}</Typography>
      {actions && <Stack direction="row" spacing={.75} useFlexGap flexWrap="wrap" alignItems="center">{actions}</Stack>}
    </Stack>
    {description && <Typography variant="body2" color="text.secondary" sx={{ maxWidth: 760 }}>{description}</Typography>}
  </Stack>;
}

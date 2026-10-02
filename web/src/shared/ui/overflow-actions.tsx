import { useState } from 'react';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import { IconButton, ListItemIcon, ListItemText, Menu, MenuItem, Tooltip } from '@mui/material';

export type OverflowAction = { label: string; icon?: React.ReactNode; disabled?: boolean; danger?: boolean; onSelect: () => void };

export function OverflowActions({ label = 'More actions', actions }: { label?: string; actions: OverflowAction[] }) {
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);
  const close = () => setAnchor(null);
  return <>
    <Tooltip title={label}><IconButton aria-label={label} aria-haspopup="menu" aria-expanded={Boolean(anchor)} onClick={(event) => setAnchor(event.currentTarget)}><MoreVertIcon /></IconButton></Tooltip>
    <Menu anchorEl={anchor} open={Boolean(anchor)} onClose={close} MenuListProps={{ 'aria-label': label }}>
      {actions.map((action) => <MenuItem key={action.label} disabled={action.disabled} onClick={() => { close(); action.onSelect(); }} sx={{ color: action.danger ? 'error.main' : undefined }}>
        {action.icon && <ListItemIcon sx={{ color: 'inherit' }}>{action.icon}</ListItemIcon>}
        <ListItemText>{action.label}</ListItemText>
      </MenuItem>)}
    </Menu>
  </>;
}

export const SESSION_EXPIRED_STATUS = 419;
export const SESSION_EXPIRED_EVENT = 'serega:session-expired';

const anonymousPaths = new Set(['/api/v1/auth/me', '/api/v1/auth/login', '/api/v1/auth/csrf']);

export function notifySessionExpired(response: Response, requestPath: string): void {
  // A removed/expired browser cookie produces 401 rather than 419. Initial
  // session discovery and invalid login credentials are expected anonymous calls.
  if (response.status === SESSION_EXPIRED_STATUS || (response.status === 401 && !anonymousPaths.has(requestPath))) {
    window.dispatchEvent(new Event(SESSION_EXPIRED_EVENT));
  }
}

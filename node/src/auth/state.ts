
export const AuthSessionCookie = 'auth-session';

export interface AuthState {
  isAuth: boolean,
  sessionID?: string,
  userID?: string,
}
import { Router } from 'express';
import { ICreateUser, createUser } from '../auth/create.js';
import { ILoginData, login } from '../auth/login.js';
import { AuthSessionCookie } from '../auth/state.js';
import { logout } from '../auth/logout.js';

export const AuthRouter = Router();

AuthRouter.post('/create', async (req, res) => {
  const data: ICreateUser = req.body;
  const result = await createUser(req['db'], data);
  if (!result) {
    res.status(500);
    res.send('Internal server error');
    return;
  }

  res.cookie(AuthSessionCookie, result.sessionID, { httpOnly: true, path: '/' });
  res.status(200);
  res.json({
    user_id: result.userID,
  });
});

AuthRouter.post('/login', async (req, res) => {
  const data: ILoginData = req.body;
  const result = await login(req['db'], data);
  if (!result) {
    res.status(500);
    res.send('Internal server error');
    return;
  }

  res.cookie(AuthSessionCookie, result, { httpOnly: true, path: '/' });
  res.status(200);
  res.send('Done');
});

AuthRouter.post('/logout', async (req, res) => {
  const cookie = req.cookies[AuthSessionCookie];
  if (!cookie) {
    res.status(400);
    res.send('Not logged in');
    return;
  }

  await logout(req['db'], cookie);
  res.cookie(AuthSessionCookie, '', { httpOnly: true, path: '/', maxAge: -1 });
  res.status(200);
  res.send('Done');
});
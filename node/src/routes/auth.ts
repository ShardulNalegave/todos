import { Router } from 'express';
import { ICreateUser, createUser } from '../auth/create.js';
import { ILoginData, login } from '../auth/login.js';
import { AuthSessionCookie, AuthState } from '../auth/state.js';
import { logout } from '../auth/logout.js';
import { PrismaClient } from '@prisma/client';

// Router containing all auth related routes
export const AuthRouter = Router();

// GET - Return currently logged in user data
AuthRouter.get('/user', async (req, res) => {
  const db: PrismaClient = req['db'];
  const state: AuthState = req['authState'];
  if (!state.isAuth) {
    res.status(401);
    res.json({
      message: 'User not logged in',
    });
    return;
  }

  const user = await db.user.findUnique({
    where: { id: state.userID },
  });

  res.status(200);
  res.json({
    id: user.id,
    name: user.name,
    email: user.email,
  });
});

// POST - Create new user with given data
AuthRouter.post('/create', async (req, res) => {
  const data: ICreateUser = req.body;
  const result = await createUser(req['db'], data);
  if (!result) {
    res.status(500);
    res.json({
      message: 'Couldn\'t create user',
    });
    return;
  }

  res.cookie(AuthSessionCookie, result.sessionID, { httpOnly: true, path: '/' });
  res.status(200);
  res.json({
    user_id: result.userID,
  });
});

// POST - Log in with provided credentials
AuthRouter.post('/login', async (req, res) => {
  const data: ILoginData = req.body;
  const result = await login(req['db'], data);
  if (!result) {
    res.status(400);
    res.json({
      message: 'Couldn\'t log in',
    });
    return;
  }

  res.cookie(AuthSessionCookie, result, { httpOnly: true, path: '/' });
  res.status(200);
  res.json({
    message: 'Done',
  });
});

// POST - Logs current user out
AuthRouter.post('/logout', async (req, res) => {
  const cookie = req.cookies[AuthSessionCookie];
  if (!cookie) {
    res.status(400);
    res.json({
      message: 'Not logged in',
    });
    return;
  }

  await logout(req['db'], cookie);
  res.cookie(AuthSessionCookie, '', { httpOnly: true, path: '/', maxAge: -1 });
  res.status(200);
  res.json({
    message: 'Done',
  });
});
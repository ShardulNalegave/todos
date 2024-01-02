
import { Request, Response } from 'express';
import { AuthSessionCookie, AuthState } from '../auth/state.js';
import { PrismaClient } from '@prisma/client';

export function AuthMiddleware() {
  return async (req: Request, res: Response, next) => {
    const cookie = req.cookies[AuthSessionCookie];
    if (!cookie || cookie === null) {
      const state: AuthState = { isAuth: false };
      req['authState'] = state;
      return await next();
    }

    const db: PrismaClient = req['db'];
    const session = await db.session.findUnique({
      where: {
        id: cookie,
      }
    });

    if (!session || session === null) {
      const state: AuthState = { isAuth: false };
      req['authState'] = state;
      return await next();
    }

    const state: AuthState = {
      isAuth: true,
      sessionID: session.id,
      userID: session.user_id,
    };
    req['authState'] = state;
    await next();
  };
}
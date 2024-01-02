import { PrismaClient } from '@prisma/client';

export function DatabaseMiddleware(prisma: PrismaClient) {
  return async (req, res, next) => {
    req['db'] = prisma;
    await next();
  };
}
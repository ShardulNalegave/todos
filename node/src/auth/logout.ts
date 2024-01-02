import { PrismaClient } from '@prisma/client';

export async function logout(db: PrismaClient, id: string) {
  await db.session.delete({
    where: {
      id,
    }
  });
}
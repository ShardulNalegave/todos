
import bcrypt from 'bcrypt';
import { PrismaClient } from '@prisma/client';

export interface ILoginData {
  email: string,
  password: string,
}

export async function login(db: PrismaClient, data: ILoginData): Promise<string | null> {
  const user = await db.user.findFirst({
    where: {
      email: data.email,
    },
  });

  if (!user) {
    return null;
  }

  if (!bcrypt.compareSync(data.password, user.password_hash)) {
    return null;
  }

  const session = await db.session.create({
    data: {
      user_id: user.id,
    },
  });

  return session.id;
}
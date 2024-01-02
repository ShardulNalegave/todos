
import bcrypt from 'bcrypt';
import { PrismaClient } from '@prisma/client';

export interface ICreateUser {
  name: string,
  email: string,
  password: string,
}

export async function createUser(db: PrismaClient, data: ICreateUser): Promise<{ sessionID: string, userID: string } | null> {
  let user = await db.user.findFirst({
    where: {
      email: data.email,
    }
  });

  if (user) {
    return null;
  }

  const hashedPass = bcrypt.hashSync(data.password, 10);
  user = await db.user.create({
    data: {
      name: data.name,
      email: data.email,
      password_hash: hashedPass,
    },
  });

  const session = await db.session.create({
    data: { user_id: user.id },
  });

  return {
    sessionID: session.id,
    userID: user.id,
  };
}
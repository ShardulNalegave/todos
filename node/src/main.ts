
import express from 'express';
import cors from 'cors';
import bodyParser from 'body-parser';
import cookieParser from 'cookie-parser';
import { PrismaClient } from '@prisma/client';
import { DatabaseMiddleware } from './middleware/db.js';
import { AuthMiddleware } from './middleware/auth.js';
import { AuthRouter } from './routes/auth.js';
import { TodosRouter } from './routes/todos.js';

const PORT = process.env.PORT || '5000';

const prisma = new PrismaClient();

const app = express();
app.use(
  cors(),
  bodyParser.json(),
  cookieParser(),
  DatabaseMiddleware(prisma),
  AuthMiddleware(),
);

app.use('/auth', AuthRouter);
app.use('/todos', TodosRouter);

console.log(`Listening at :${PORT}`);
app.listen(PORT, async () => {
  await prisma.$disconnect();
});
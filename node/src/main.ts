
import express from 'express';
import cors from 'cors';
import bodyParser from 'body-parser';
import { PrismaClient } from '@prisma/client';

const PORT = process.env.PORT || '5000';

const prisma = new PrismaClient();

const app = express();
app.use(bodyParser.json(), cors());

app.get('/', (req, res) => {
  res.send('Hello, World!');
});

app.get('/todos', async (req, res) => {
  const todos = await prisma.todo.findMany();
  res.json(todos);
});

console.log(`Listening at :${PORT}`);
app.listen(PORT, async () => {
  await prisma.$disconnect();
});
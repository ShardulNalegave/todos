import { PrismaClient } from '@prisma/client';
import { Router } from 'express';
import { AuthState } from '../auth/state.js';

export const TodosRouter = Router();

TodosRouter.post('/', async (req, res) => {
  const data: ICreateTodo = req.body;
  const db: PrismaClient = req['db'];
  const state: AuthState = req['authState'];
  if (!state.isAuth) {
    res.status(401);
    res.send('User unauthenticated');
    return;
  }

  const todo = await db.todo.create({
    data: {
      created_by: state.userID,
      ...data,
    },
  });
  res.status(200);
  res.json(todo);
});

TodosRouter.get('/', async (req, res) => {
  const db: PrismaClient = req['db'];
  const state: AuthState = req['authState'];
  if (!state.isAuth) {
    res.status(401);
    res.send('User unauthenticated');
    return;
  }

  const todos = await db.todo.findMany({
    where: {
      created_by: state.userID,
    },
  });
  res.status(200);
  res.json(todos);
});

TodosRouter.get('/:id', async (req, res) => {
  const id = req.params.id;
  const db: PrismaClient = req['db'];
  const state: AuthState = req['authState'];
  if (!state.isAuth) {
    res.status(401);
    res.send('User unauthenticated');
    return;
  }

  const todos = await db.todo.findUnique({
    where: {
      id,
      created_by: state.userID,
    },
  });
  res.status(200);
  res.json(todos);
});

TodosRouter.delete('/:id', async (req, res) => {
  const id = req.params.id;
  const db: PrismaClient = req['db'];
  const state: AuthState = req['authState'];
  if (!state.isAuth) {
    res.status(401);
    res.send('User unauthenticated');
    return;
  }

  const todo = await db.todo.findUnique({
    where: { id },
  });
  if (todo.created_by != state.userID) {
    res.status(401);
    res.send('Cannot delete other\'s todos');
    return;
  }

  await db.todo.delete({
    where: { id },
  });
  res.status(200);
  res.send('Done');
});

TodosRouter.put('/:id', async (req, res) => {
  const id = req.params.id;
  const data: IUpdateTodo = req.body;
  const db: PrismaClient = req['db'];
  const state: AuthState = req['authState'];
  if (!state.isAuth) {
    res.status(401);
    res.send('User unauthenticated');
    return;
  }

  const todo = await db.todo.findUnique({
    where: { id },
  });
  if (todo.created_by != state.userID) {
    res.status(401);
    res.send('Cannot edit other\'s todos');
    return;
  }

  await db.todo.update({
    where: {
      id,
      created_by: state.userID,
    },
    data,
  });
  res.status(200);
  res.json(await db.todo.findUnique({ where: { id } }));
});

interface IUpdateTodo {
  content: string,
  completed: boolean,
}

interface ICreateTodo {
  content: string,
  completed: boolean,
}
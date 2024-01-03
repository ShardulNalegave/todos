import { PrismaClient } from '@prisma/client';
import { Router } from 'express';
import { AuthState } from '../auth/state.js';

// Router containing all todos related routes
export const TodosRouter = Router();

// POST - Creates a new Todo
TodosRouter.post('/', async (req, res) => {
  const data: ICreateTodo = req.body;
  const db: PrismaClient = req['db'];
  const state: AuthState = req['authState'];
  if (!state.isAuth) {
    res.status(401);
    res.json({
      message: 'User unauthenticated',
    });
    return;
  }

  const todo = await db.todo.create({
    data: {
      created_by: state.userID,
      completed: false,
      ...data,
    },
  });
  res.status(200);
  res.json(todo);
});

// GET - Returns all Todos by current user
TodosRouter.get('/', async (req, res) => {
  const db: PrismaClient = req['db'];
  const state: AuthState = req['authState'];
  if (!state.isAuth) {
    res.status(401);
    res.json({
      message: 'User unauthenticated',
    });
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

// GET - Returns Todo with given ID created by current user
TodosRouter.get('/:id', async (req, res) => {
  const id = req.params.id;
  const db: PrismaClient = req['db'];
  const state: AuthState = req['authState'];
  if (!state.isAuth) {
    res.status(401);
    res.json({
      message: 'User unauthenticated',
    });
    return;
  }

  const todo = await db.todo.findUnique({
    where: {
      id,
      created_by: state.userID,
    },
  });

  if (!todo) {
    res.status(404);
    res.json({
      message: 'No such Todo',
    });
  }

  res.status(200);
  res.json(todo);
});

// DELETE - Deletes Todo with given ID if created by current user
TodosRouter.delete('/:id', async (req, res) => {
  const id = req.params.id;
  const db: PrismaClient = req['db'];
  const state: AuthState = req['authState'];
  if (!state.isAuth) {
    res.status(401);
    res.json({
      message: 'User unauthenticated',
    });
    return;
  }

  const todo = await db.todo.findUnique({
    where: { id },
  });

  if (!todo) {
    res.status(404);
    res.json({
      message: 'No such Todo',
    });
  }

  if (todo.created_by != state.userID) {
    res.status(401);
    res.json({
      message: 'Cannot delete other\'s todos',
    });
    return;
  }

  await db.todo.delete({
    where: { id },
  });
  res.status(200);
  res.json({
    message: 'Done',
  });
});

// PUT - Updates Todo with given ID if created by current user
TodosRouter.put('/:id', async (req, res) => {
  const id = req.params.id;
  const data: IUpdateTodo = req.body;
  const db: PrismaClient = req['db'];
  const state: AuthState = req['authState'];
  if (!state.isAuth) {
    res.status(401);
    res.json({
      message: 'User unauthenticated',
    });
    return;
  }

  const todo = await db.todo.findUnique({
    where: { id },
  });

  if (!todo) {
    res.status(404);
    res.json({
      message: 'No such Todo',
    });
  }

  if (todo.created_by != state.userID) {
    res.status(401);
    res.json({
      message: 'Cannot edit other\'s todos',
    });
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
}
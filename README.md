
# Todos
I'm implementing the same exact backend for a todo application in three languages:-
- [ ] Go
- [ ] Rust
- [ ] Typescript

My goal is to find out how the developer experience is when using them and if and why they are best suited for writing backend servers.
The backend being implemented has quite a lot going on:-
- It uses **Sqlite** for storing data
- **GET, POST, PUT, DELETE** routes
- Implement **cookie-based authentication** from scratch in all languages

## Routes

- `/todos` - Method `GET`
  - Lists all Todos
- `/todos` - Method `POST`
  - Create new Todo
- `/todos/:id` - Method `GET`
  - Get Todo with given ID
- `/todos/:id` - Method `PUT`
  - Update Todo with given ID
- `/todos/:id` - Method `DELETE`
  - Delete Todo with given ID
- `/auth/create` - Method `POST`
  - Create new User and login
- `/auth/login` - Method `POST`
  - Login with given email and password
- `/auth/logout` - Method `POST`
  - Logout
- `/auth/user` - Method `GET`
  - Returns currently logged in user data
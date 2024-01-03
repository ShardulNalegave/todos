
// ===== Imports =====
use axum::{extract::{State, Path}, extract::Json, http::StatusCode, Extension};
use sea_orm::{EntityTrait, QueryFilter, ColumnTrait, Set, ActiveModelTrait};
use serde::{Serialize, Deserialize};
use serde_json::{Value, json};
use crate::{context::Context, auth::AuthState};
// ===================

pub async fn todos(
  State(context): State<Context>,
  Extension(auth_state): Extension<AuthState>,
) -> (StatusCode, Json<Value>) {
  match auth_state {
    AuthState::Unauthenticated => (
      StatusCode::UNAUTHORIZED,
      Json(json!({ "message": "Not logged in" })),
    ),
    AuthState::Authenticated(auth_data) => {
      let todos_res = entity::todo::Entity::find()
        .filter(entity::todo::Column::CreatedBy.eq(auth_data.user_id))
        .into_json()
        .all(&context.db).await;

      match todos_res {
        Err(_) => (
          StatusCode::INTERNAL_SERVER_ERROR,
          Json(json!({ "message": "Couldn't read todos" })),
        ),
        Ok(todos) => (
          StatusCode::OK,
          Json(serde_json::to_value(todos).unwrap()),
        ),
      }
    },
  }
}

pub async fn todo(
  State(context): State<Context>,
  Path(todo_id): Path<String>,
  Extension(auth_state): Extension<AuthState>,
) -> (StatusCode, Json<Value>) {
  match auth_state {
    AuthState::Unauthenticated => (
      StatusCode::UNAUTHORIZED,
      Json(json!({ "message": "Not logged in" })),
    ),
    AuthState::Authenticated(auth_data) => {
      let todos_res = entity::todo::Entity::find_by_id(todo_id)
        .filter(entity::todo::Column::CreatedBy.eq(auth_data.user_id))
        .into_json()
        .one(&context.db).await;

      match todos_res {
        Err(err) => (
          StatusCode::INTERNAL_SERVER_ERROR,
          Json(json!({ "message": err.to_string() })),
        ),
        Ok(todo_option) => match todo_option {
          None => (StatusCode::NOT_FOUND, Json(json!({ "message": "Todo Not Found" }))),
          Some(todo) => (
            StatusCode::OK,
            Json(serde_json::to_value(todo).unwrap()),
          ),
        },
      }
    },
  }
} 

pub async fn delete_todo(
  State(context): State<Context>,
  Path(todo_id): Path<String>,
  Extension(auth_state): Extension<AuthState>,
) -> (StatusCode, Json<Value>) {
  match auth_state {
    AuthState::Unauthenticated => (
      StatusCode::UNAUTHORIZED,
      Json(json!({ "message": "Not logged in" })),
    ),
    AuthState::Authenticated(auth_data) => {
      let todo = entity::todo::ActiveModel {
        id: Set(todo_id),
        created_by: Set(auth_data.user_id),
        ..Default::default()
      };

      match entity::todo::Entity::delete(todo).exec(&context.db).await {
        Err(_) => (
          StatusCode::INTERNAL_SERVER_ERROR,
          Json(json!({ "message": "Could not delete todo" })),
        ),
        Ok(_) => (
          StatusCode::OK,
          Json(json!({ "message": "Done" })),
        ),
      }
    },
  }
}

pub async fn create_todo(
  State(context): State<Context>,
  Extension(auth_state): Extension<AuthState>,
  Json(payload): Json<CreateTodoPayload>,
) -> (StatusCode, Json<Value>) {
  match auth_state {
    AuthState::Unauthenticated => (
      StatusCode::UNAUTHORIZED,
      Json(json!({ "message": "Not logged in" })),
    ),
    AuthState::Authenticated(auth_data) => {
      let todo_id = uuid::Uuid::new_v4().to_string();
      let todo = entity::todo::ActiveModel {
        id: Set(todo_id.clone()),
        content: Set(payload.content),
        completed: Set(payload.completed),
        created_by: Set(auth_data.user_id.clone()),
      };

      match entity::todo::Entity::insert(todo).exec(&context.db).await {
        Err(_) => (
          StatusCode::INTERNAL_SERVER_ERROR,
          Json(json!({ "message": "Could not create todo" })),
        ),
        Ok(_) => {
          let todos_res = entity::todo::Entity::find_by_id(todo_id)
            .filter(entity::todo::Column::CreatedBy.eq(auth_data.user_id))
            .into_json()
            .one(&context.db).await;

          match todos_res {
            Err(err) => (
              StatusCode::INTERNAL_SERVER_ERROR,
              Json(json!({ "message": err.to_string() })),
            ),
            Ok(todo_option) => match todo_option {
              None => (StatusCode::INTERNAL_SERVER_ERROR, Json(json!({ "message": "Something went wrong" }))),
              Some(todo) => (
                StatusCode::OK,
                Json(serde_json::to_value(todo).unwrap()),
              ),
            },
          }
        },
      }
    },
  }
}

pub async fn update_todo(
  State(context): State<Context>,
  Path(todo_id): Path<String>,
  Extension(auth_state): Extension<AuthState>,
  Json(payload): Json<UpdateTodoPayload>,
) -> (StatusCode, Json<Value>) {
  match auth_state {
    AuthState::Unauthenticated => (
      StatusCode::UNAUTHORIZED,
      Json(json!({ "message": "Not logged in" })),
    ),
    AuthState::Authenticated(auth_data) => {
      let todo_opt = entity::todo::Entity::find_by_id(todo_id.clone())
      .filter(entity::todo::Column::CreatedBy.eq(auth_data.user_id.clone()))
      .one(&context.db).await.unwrap();

      match todo_opt {
        None => (
          StatusCode::INTERNAL_SERVER_ERROR,
          Json(json!({ "message": "Todo doesn't exist" })),
        ),
        Some(todo) => {
          let mut todo: entity::todo::ActiveModel = todo.into();
          todo.content = Set(payload.content);
          todo.completed = Set(payload.completed);
          match todo.save(&context.db).await {
            Err(err) => (
              StatusCode::INTERNAL_SERVER_ERROR,
              Json(json!({ "message": err.to_string() })),
            ),
            Ok(_) => {
              let todos_res = entity::todo::Entity::find_by_id(todo_id)
                .filter(entity::todo::Column::CreatedBy.eq(auth_data.user_id))
                .into_json()
                .one(&context.db).await;

              match todos_res {
                Err(err) => (
                  StatusCode::INTERNAL_SERVER_ERROR,
                  Json(json!({ "message": err.to_string() })),
                ),
                Ok(todo_option) => match todo_option {
                  None => (StatusCode::INTERNAL_SERVER_ERROR, Json(json!({ "message": "Something went wrong" }))),
                  Some(todo) => (
                    StatusCode::OK,
                    Json(serde_json::to_value(todo).unwrap()),
                  ),
                },
              }
            },
          }
        },
      }
    },
  }
}

#[derive(Serialize, Deserialize)]
pub struct CreateTodoPayload {
  pub content: String,
  pub completed: bool,
}

#[derive(Serialize, Deserialize)]
pub struct UpdateTodoPayload {
  pub content: String,
  pub completed: bool,
}
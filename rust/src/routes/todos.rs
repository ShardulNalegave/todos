
// ===== Imports =====
use axum::{extract::{State, Path}, extract::Json, http::StatusCode, Extension};
use sea_orm::{EntityTrait, QueryFilter, ColumnTrait};
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
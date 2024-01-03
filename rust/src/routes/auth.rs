
// ===== Imports ======
use axum::{extract::State, extract::Json, http::StatusCode, Extension};
use sea_orm::EntityTrait;
use serde_json::{Value, json};
use tower_cookies::{Cookies, Cookie};
use crate::{context::Context, auth::{self, AUTH_SESSION_COOKIE, AuthState}};
// ====================

// GET - Returns currently logged in user data
pub async fn current_user(
  State(context): State<Context>,
  Extension(auth_state): Extension<auth::AuthState>,
) -> (StatusCode, Json<Value>) {
  match auth_state {
    AuthState::Unauthenticated => (
      StatusCode::UNAUTHORIZED,
      Json(json!({
        "message": "User not logged in",
      })),
    ),
    AuthState::Authenticated(auth_data) => {
      let user = entity::user::Entity::find_by_id(auth_data.user_id)
        .one(&context.db).await
        .expect("Couldn't read data")
        .expect("No such user");

      (StatusCode::OK, Json(json!({
        "id": user.id,
        "name": user.name,
        "email": user.email,
      })))
    },
  }
}

// POST - Logs in with provided credentials
pub async fn login(
  State(context): State<Context>,
  cookies: Cookies,
  Json(payload): Json<auth::LoginPayload>,
) -> (StatusCode, Json<Value>) {
  match auth::login(&context.db, payload).await {
    Err(_) => (
      StatusCode::BAD_REQUEST,
      Json(json!({
        "message": "Couldn't log in",
      })),
    ),
    Ok((session_id, _user_id)) => {
      cookies.add(
        Cookie::build((AUTH_SESSION_COOKIE, session_id))
          .path("/")
          .http_only(true)
          .build()
      );
      (
        StatusCode::OK,
        Json(json!({
          "message": "Done",
        })),
      )
    },
  }
}

// POST - Creates user with given credentials and logs them in
pub async fn create_user(
  State(context): State<Context>,
  cookies: Cookies,
  Json(payload): Json<auth::CreateUserPayload>,
) -> (StatusCode, Json<Value>) {
  match auth::create_user(&context.db, payload).await {
    Ok((session_id, user_id)) => {
      cookies.add(
        Cookie::build((AUTH_SESSION_COOKIE, session_id))
          .path("/")
          .http_only(true)
          .build()
      );
      (
        StatusCode::OK,
        Json(json!({
          "user_id": user_id,
        })),
      )
    },
    Err(_) => (
      StatusCode::BAD_REQUEST,
      Json(json!({
        "message": "Couldn't create user",
      })),
    ),
  }
}

// POST - Logs the user out
pub async fn logout(
  State(context): State<Context>,
  cookies: Cookies,
  Extension(auth_state): Extension<auth::AuthState>,
) -> (StatusCode, Json<Value>) {
  match auth_state {
    AuthState::Unauthenticated => (
      StatusCode::UNAUTHORIZED,
      Json(json!({ "message": "Not logged in" })),
    ),
    AuthState::Authenticated(auth_data) => match auth::logout(&context.db, auth_data.session_id.clone()).await {
      Err(_) => (
        StatusCode::INTERNAL_SERVER_ERROR,
        Json(json!({
          "message": "Could not log out"
        })),
      ),
      Ok(_) =>  {
        cookies.remove(
          Cookie::build(AUTH_SESSION_COOKIE)
            .path("/")
            .http_only(true)
            .build(),
        );
        (
          StatusCode::OK,
          Json(json!({
            "message": "Done"
          })),
        )
      },
    },
  }
}

// ===== Imports ======
use axum::{extract::State, extract::Json, http::StatusCode, Extension};
use serde_json::{Value, json};
use tower_cookies::{Cookies, Cookie};
use crate::{context::Context, auth::{self, AUTH_SESSION_COOKIE, AuthState}};
// ====================

pub async fn login(
  State(context): State<Context>,
  cookies: Cookies,
  Json(payload): Json<auth::LoginPayload>,
) -> (StatusCode, Json<Value>) {
  match auth::login(&context.db, payload).await {
    Err(err) => (
      StatusCode::INTERNAL_SERVER_ERROR,
      Json(json!({
        "message": err.to_string(),
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

pub async fn create_user(
  State(context): State<Context>,
  cookies: Cookies,
  Json(payload): Json<auth::CreateUserPayload>,
) -> (StatusCode, Json<Value>) {
  match auth::create_user(&context.db, payload).await {
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
    Err(err) => (
      StatusCode::INTERNAL_SERVER_ERROR,
      Json(json!({
        "message": err.to_string(),
      })),
    ),
  }
}

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
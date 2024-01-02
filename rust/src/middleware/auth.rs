
// ===== Imports =====
use axum::{
  extract::{Request, State},
  middleware::Next,
  response::Response,
};
use sea_orm::EntityTrait;
use tower_cookies::Cookies;
use crate::{context::Context, auth::{AUTH_SESSION_COOKIE, AuthState}};
// ===================


pub async fn auth_middleware(
  State(mut state): State<Context>,
  cookies: Cookies,
  request: Request,
  next: Next,
) -> Response {
  let session_cookie = cookies.get(AUTH_SESSION_COOKIE);
  match session_cookie {
    Some(session_cookie) => match entity::user::Entity::find_by_id(session_cookie.value()).one(&state.db).await {
      Ok(user) => match user {
        Some(user) => {
          state.auth_state = Some(AuthState {
            session_id: session_cookie.value().to_owned(),
            user_id: user.id,
          });
        },
        None => state.auth_state = None,
      },
      Err(_) => state.auth_state = None,
    },
    None => state.auth_state = None,
  }

  next.run(request).await
}
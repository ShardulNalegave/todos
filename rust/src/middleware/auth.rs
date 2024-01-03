
// ===== Imports =====
use axum::{
  extract::{Request, State},
  middleware::Next,
  response::Response,
};
use sea_orm::EntityTrait;
use tower_cookies::Cookies;
use crate::{context::Context, auth::{AUTH_SESSION_COOKIE, AuthState, AuthData}};
// ===================

pub async fn auth_middleware(
  State(state): State<Context>,
  cookies: Cookies,
  mut request: Request,
  next: Next,
) -> Response {
  let session_cookie = match cookies.get(AUTH_SESSION_COOKIE) {
    None => {
      request.extensions_mut().insert(AuthState::Unauthenticated);
      return next.run(request).await
    },
    Some(session_cookie) => session_cookie,
  };

  match entity::session::Entity::find_by_id(session_cookie.value()).one(&state.db).await {
    Err(_) => { request.extensions_mut().insert(AuthState::Unauthenticated); },
    Ok(session) => match session {
      None => { request.extensions_mut().insert(AuthState::Unauthenticated); },
      Some(session) => {
        request.extensions_mut().insert(AuthState::Authenticated(AuthData{
          session_id: session_cookie.value().to_owned(),
          user_id: session.user_id,
        }));
      },
    },
  }

  next.run(request).await
}
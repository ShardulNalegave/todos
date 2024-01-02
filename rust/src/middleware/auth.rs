
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
  State(state): State<Context>,
  cookies: Cookies,
  mut request: Request,
  next: Next,
) -> Response {
  let session_cookie = cookies.get(AUTH_SESSION_COOKIE);
  match session_cookie {
    Some(session_cookie) => match entity::session::Entity::find_by_id(session_cookie.value()).one(&state.db).await {
      Ok(session) => match session {
        Some(session) => {
          request.extensions_mut().insert(AuthState {
            session_id: session_cookie.value().to_owned(),
            user_id: session.user_id,
          });
        },
        None => {},
      },
      Err(_) => {},
    },
    None => {},
  }

  next.run(request).await
}

#[derive(Debug, Clone)]
pub enum AuthState {
  Unauthenticated,
  Authenticated(AuthData),
}

#[derive(Debug, Clone)]
pub struct AuthData {
  pub session_id: String,
  pub user_id: String,
}
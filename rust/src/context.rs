
// ===== Imports =====
use sea_orm::DatabaseConnection;
use crate::auth::AuthState;
// ===================

#[derive(Clone)]
pub struct Context {
  pub db: DatabaseConnection,
  pub auth_state: Option<AuthState>,
}
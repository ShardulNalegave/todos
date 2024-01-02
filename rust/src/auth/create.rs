
// ===== Imports =====
use anyhow::{Result, anyhow};
use sea_orm::{DatabaseConnection, EntityTrait, QueryFilter, Set, ColumnTrait};
use serde::{Serialize, Deserialize};
// ===================

#[derive(Serialize, Deserialize)]
pub struct CreateUserPayload {
  name: String,
  email: String,
  password: String,
}

pub async fn create_user(db: &DatabaseConnection, user: CreateUserPayload) -> Result<(String, String)> {
  let users = entity::user::Entity::find()
    .filter(entity::user::Column::Email.eq(user.email.clone()))
    .all(db).await?;
  if !users.is_empty() {
    return Err(anyhow!("User with given email already exists"));
  }

  let user_id = uuid::Uuid::new_v4().to_string();
  let password_hash = bcrypt::hash(user.password, 10)?;
  let user_doc = entity::user::ActiveModel {
    id: Set(user_id.clone()),
    name: Set(user.name),
    email: Set(user.email.clone()),
    password_hash: Set(password_hash),
  };
  entity::user::Entity::insert(user_doc).exec(db).await?;

  let session_id = uuid::Uuid::new_v4().to_string();
  let session_doc = entity::session::ActiveModel {
    id: Set(session_id.clone()),
    user_id: Set(user_id.clone()),
  };
  entity::session::Entity::insert(session_doc).exec(db).await?;

  Ok((session_id, user_id))
}
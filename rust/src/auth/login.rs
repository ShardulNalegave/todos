use anyhow::{Result, anyhow};
use sea_orm::{DatabaseConnection, EntityTrait, QueryFilter, Set, ColumnTrait};
use serde::{Serialize, Deserialize};


#[derive(Serialize, Deserialize)]
pub struct LoginPayload {
  email: String,
  password: String,
}

pub async fn login(db: &DatabaseConnection, user: LoginPayload) -> Result<(String, String)> {
  let user = entity::user::Entity::find()
    .filter(entity::user::Column::Email.eq(user.email.clone()))
    .one(db).await?;
  if user.is_none() {
    return Err(anyhow!("User does not exist"));
  }
  let user = user.unwrap();

  let session_id = uuid::Uuid::new_v4().to_string();
  let session_doc = entity::session::ActiveModel {
    id: Set(session_id.clone()),
    user_id: Set(user.id.clone()),
  };
  entity::session::Entity::insert(session_doc).exec(db).await?;

  Ok((session_id, user.id.clone()))
}
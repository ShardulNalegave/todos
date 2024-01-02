
use anyhow::Result;
// ===== Imports =====
use sea_orm::{DatabaseConnection, EntityTrait};
// ===================

pub async fn logout(db: &DatabaseConnection, id: String) -> Result<()> {
  entity::user::Entity::delete_by_id(id).exec(db).await?;
  Ok(())
}
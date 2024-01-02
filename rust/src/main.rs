
// ===== Imports =====
use anyhow::Result;
use migration::{Migrator, MigratorTrait};
use sea_orm::{DatabaseConnection, Database};
use entity::prelude::*;
// ===================

#[tokio::main]
async fn main() -> Result<()> {
  let db: DatabaseConnection = Database::connect("sqlite://todos.sqlite?mode=rwc").await?;
  Migrator::up(&db, None).await?;

  db.close().await?;
  Ok(())
}

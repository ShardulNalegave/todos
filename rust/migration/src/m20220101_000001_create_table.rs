use sea_orm_migration::prelude::*;

#[derive(DeriveMigrationName)]
pub struct Migration;

#[async_trait::async_trait]
impl MigrationTrait for Migration {
  async fn up(&self, manager: &SchemaManager) -> Result<(), DbErr> {
    let sessions = Table::create()
      .table(Session::Table)
      .if_not_exists()
      .col(ColumnDef::new(Session::ID).string().not_null().primary_key())
      .col(ColumnDef::new(Session::UserID).string().not_null())
      .to_owned();

    let users = Table::create()
      .table(User::Table)
      .if_not_exists()
      .col(ColumnDef::new(User::ID).string().not_null().primary_key())
      .col(ColumnDef::new(User::Name).string().not_null())
      .col(ColumnDef::new(User::Email).string().not_null())
      .col(ColumnDef::new(User::PasswordHash).string().not_null())
      .to_owned();

    let todos = Table::create()
      .table(Todo::Table)
      .if_not_exists()
      .col(ColumnDef::new(Todo::ID).string().not_null().primary_key())
      .col(ColumnDef::new(Todo::Content).string().not_null())
      .col(ColumnDef::new(Todo::Completed).boolean().not_null())
      .col(ColumnDef::new(Todo::CreatedBy).string().not_null())
      .to_owned();

    manager.create_table(sessions).await?;
    manager.create_table(users).await?;
    manager.create_table(todos).await?;

    Ok(())
  }

  async fn down(&self, manager: &SchemaManager) -> Result<(), DbErr> {
    manager.drop_table(Table::drop().table(Session::Table).to_owned()).await?;
    manager.drop_table(Table::drop().table(User::Table).to_owned()).await?;
    manager.drop_table(Table::drop().table(Todo::Table).to_owned()).await?;

    Ok(())
  }
}

#[derive(DeriveIden)]
enum Session {
  Table,
  ID,
  UserID,
}

#[derive(DeriveIden)]
enum User {
  Table,
  ID,
  Name,
  Email,
  PasswordHash,
}

#[derive(DeriveIden)]
enum Todo {
  Table,
  ID,
  Content,
  Completed,
  CreatedBy,
}
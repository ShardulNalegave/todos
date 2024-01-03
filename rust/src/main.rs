
pub mod auth;
pub mod context;
pub mod middleware;
pub mod routes;

// ===== Imports =====
use anyhow::Result;
use axum::{
  Router,
  routing::{post, get},
};
use migration::{Migrator, MigratorTrait};
use sea_orm::{DatabaseConnection, Database};
use tower_cookies::CookieManagerLayer;
// ===================

#[tokio::main]
async fn main() -> Result<()> {
  let port = std::env::var("PORT").unwrap_or("5000".to_owned());

  let db: DatabaseConnection = Database::connect("sqlite://todos.sqlite?mode=rwc").await?;
  Migrator::up(&db, None).await?;

  let ctx = context::Context { db };

  let app = Router::new()
    .route("/auth/logout", post(routes::auth::logout))
    .route("/auth/create", post(routes::auth::create_user))
    .route("/auth/login", post(routes::auth::login))
    .route("/todos", get(routes::todos::todos).post(routes::todos::create_todo))
    .route("/todos/:id", get(routes::todos::todo).delete(routes::todos::delete_todo).put(routes::todos::update_todo))
    .route_layer(axum::middleware::from_fn_with_state(ctx.clone(), middleware::auth_middleware))
    .layer(CookieManagerLayer::new())
    .with_state(ctx);

  println!("Listening at :{}", port);
  let listener = tokio::net::TcpListener::bind(format!("0.0.0.0:{}", port)).await?;
  axum::serve(listener, app).await?;

  Ok(())
}


mod login;
pub use login::{login, LoginPayload};

mod create;
pub use create::{create_user, CreateUserPayload};

mod logout;
pub use logout::logout;

mod state;
pub use state::AuthState;

pub const AUTH_SESSION_COOKIE: &str = "auth-session";
mod domain;
mod adapter;
mod app;
mod config;

use std::sync::Arc;
use axum::{
    routing::{get, post},
    Router,
};
use tower_http::services::ServeDir;
use tower_sessions::{MemoryStore, SessionManagerLayer};
use crate::adapter::web::handlers::{AppState, handle_index_page, handle_about_page, handle_login_page, handle_login_post, handle_logout, handle_prize_page, handle_admin_page, handle_user_status_switch};
use crate::adapter::mongodb::{prize_repository::PrizeMongoRepository, user_repository::UserMongoRepository};
use crate::app::{prize_service::PrizeService, user_service::UserService, auth_service::AuthService};
use crate::config::Config;
use mongodb::Client;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    tracing_subscriber::fmt::init();
    
    let cfg = Config::load();
    
    let client = Client::with_uri_str(&cfg.mongodb_url).await?;
    let db = client.database("test");
    
    let prize_repo = Arc::new(PrizeMongoRepository { db: db.clone() });
    let user_repo = Arc::new(UserMongoRepository { db: db.clone() });
    
    let prize_service = PrizeService::new(prize_repo);
    let user_service = UserService::new(user_repo.clone());
    let auth_service = AuthService::new(user_repo);
    
    let state = Arc::new(AppState {
        prize_service,
        user_service,
        auth_service,
    });
    
    let session_store = MemoryStore::default();
    let session_layer = SessionManagerLayer::new(session_store)
        .with_secure(false); // Set to true in production
    
    let app = Router::new()
        .route("/", get(handle_index_page))
        .route("/about", get(handle_about_page))
        .route("/login", get(handle_login_page).post(handle_login_post))
        .route("/logout", post(handle_logout))
        .route("/prize", get(handle_prize_page))
        .route("/admin", get(handle_admin_page))
        .route("/api/switch/:id", post(handle_user_status_switch))
        .nest_service("/static", ServeDir::new("internal/adapter/web/static"))
        .layer(session_layer)
        .with_state(state);
    
    let addr = format!("0.0.0.0:{}", cfg.port);
    let listener = tokio::net::TcpListener::bind(&addr).await?;
    println!("Listening on {}", addr);
    axum::serve(listener, app).await?;
    
    Ok(())
}

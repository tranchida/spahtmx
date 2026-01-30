use axum::{
    extract::{State, Path, Query},
    response::{IntoResponse, Redirect},
    Form,
    http::HeaderMap,
};
use serde::Deserialize;
use std::sync::Arc;
use crate::app::{user_service::UserService, prize_service::PrizeService, auth_service::AuthService};
use crate::domain::models::{User, Prize};
use askama::Template;
use tower_sessions::Session;

pub struct AppState {
    pub user_service: UserService,
    pub prize_service: PrizeService,
    pub auth_service: AuthService,
}

#[derive(Template)]
#[template(path = "index.html")]
pub struct IndexTemplate {
    pub page: String,
    pub user: Option<User>,
}

#[derive(Template)]
#[template(path = "about.html")]
pub struct AboutTemplate {
    pub page: String,
    pub user: Option<User>,
}

#[derive(Template)]
#[template(path = "login.html")]
pub struct LoginTemplate {
    pub page: String,
    pub user: Option<User>,
    pub error_msg: Option<String>,
}

#[derive(Template)]
#[template(path = "prize.html")]
pub struct PrizeTemplate {
    pub page: String,
    pub user: Option<User>,
    pub prizes: Vec<Prize>,
    pub categories: Vec<String>,
    pub years: Vec<String>,
    pub selected_category: String,
    pub selected_year: String,
}

#[derive(Template)]
#[template(path = "prize_list.html")]
pub struct PrizeListTemplate {
    pub prizes: Vec<Prize>,
}

#[derive(Template)]
#[template(path = "admin.html")]
pub struct AdminTemplate {
    pub page: String,
    pub user: Option<User>,
    pub users: Vec<User>,
    pub user_count: String,
    pub page_view: String,
}

#[derive(Template)]
#[template(path = "user_list.html")]
pub struct UserListTemplate {
    pub users: Vec<User>,
}

pub async fn handle_index_page(
    State(state): State<Arc<AppState>>,
    session: Session,
) -> impl IntoResponse {
    let user = session.get::<User>("user").await.unwrap_or(None);
    IndexTemplate {
        page: "/".to_string(),
        user,
    }
}

pub async fn handle_about_page(
    State(state): State<Arc<AppState>>,
    session: Session,
) -> impl IntoResponse {
    let user = session.get::<User>("user").await.unwrap_or(None);
    AboutTemplate {
        page: "/about".to_string(),
        user,
    }
}

pub async fn handle_login_page(
    State(state): State<Arc<AppState>>,
    session: Session,
) -> impl IntoResponse {
    let user = session.get::<User>("user").await.unwrap_or(None);
    LoginTemplate {
        page: "/login".to_string(),
        user,
        error_msg: None,
    }
}

#[derive(Deserialize)]
pub struct LoginForm {
    pub username: String,
    pub password: String,
}

pub async fn handle_login_post(
    State(state): State<Arc<AppState>>,
    session: Session,
    Form(form): Form<LoginForm>,
) -> impl IntoResponse {
    match state.auth_service.login(&form.username, &form.password).await {
        Ok(user) => {
            session.insert("user", user.clone()).await.unwrap();
            Redirect::to("/").into_response()
        }
        Err(_) => {
            LoginTemplate {
                page: "/login".to_string(),
                user: None,
                error_msg: Some("Nom d'utilisateur ou mot de passe incorrect".to_string()),
            }.into_response()
        }
    }
}

pub async fn handle_logout(
    session: Session,
) -> impl IntoResponse {
    session.remove::<User>("user").await.unwrap();
    Redirect::to("/").into_response()
}

#[derive(Deserialize)]
pub struct PrizeQuery {
    pub category: Option<String>,
    pub year: Option<String>,
}

pub async fn handle_prize_page(
    State(state): State<Arc<AppState>>,
    session: Session,
    Query(query): Query<PrizeQuery>,
    headers: HeaderMap,
) -> impl IntoResponse {
    let user = session.get::<User>("user").await.unwrap_or(None);
    let category = query.category.clone().unwrap_or_default();
    let year = query.year.clone().unwrap_or_default();

    let prizes = if !category.is_empty() && !year.is_empty() {
        state.prize_service.get_prizes_by_category_and_year(&category, &year).await.unwrap_or_default()
    } else if !category.is_empty() {
        state.prize_service.get_prizes_by_category(&category).await.unwrap_or_default()
    } else if !year.is_empty() {
        state.prize_service.get_prizes_by_year(&year).await.unwrap_or_default()
    } else {
        state.prize_service.get_prizes().await.unwrap_or_default()
    };

    let categories = state.prize_service.get_categories().await.unwrap_or_default();
    let years = state.prize_service.get_years().await.unwrap_or_default();

    if headers.contains_key("HX-Request") && !headers.contains_key("HX-Boosted") {
        return PrizeListTemplate { prizes }.into_response();
    }

    PrizeTemplate {
        page: "/prize".to_string(),
        user,
        prizes,
        categories,
        years,
        selected_category: category,
        selected_year: year,
    }.into_response()
}

pub async fn handle_admin_page(
    State(state): State<Arc<AppState>>,
    session: Session,
) -> impl IntoResponse {
    let user = session.get::<User>("user").await.unwrap_or(None);
    if user.is_none() {
        return Redirect::to("/login").into_response();
    }

    let users = state.user_service.get_users().await.unwrap_or_default();
    let user_count = state.user_service.get_user_count();
    let page_view = state.user_service.get_page_view();

    AdminTemplate {
        page: "/admin".to_string(),
        user,
        users,
        user_count,
        page_view,
    }.into_response()
}

pub async fn handle_user_status_switch(
    State(state): State<Arc<AppState>>,
    Path(id): Path<String>,
) -> impl IntoResponse {
    let _ = state.user_service.update_user_status(&id).await;
    let users = state.user_service.get_users().await.unwrap_or_default();
    UserListTemplate { users }.into_response()
}

use std::sync::Arc;
use crate::domain::repositories::UserRepository;
use crate::domain::models::User;
use bcrypt::{verify, hash, DEFAULT_COST};

pub struct AuthService {
    user_repo: Arc<dyn UserRepository>,
}

impl AuthService {
    pub fn new(user_repo: Arc<dyn UserRepository>) -> Self {
        Self { user_repo }
    }

    pub async fn login(&self, username: &str, password: &str) -> Result<User, Box<dyn std::error::Error + Send + Sync>> {
        let user = self.user_repo.get_by_username(username).await?;
        
        if verify(password, &user.password)? {
            Ok(user)
        } else {
            Err("Unauthorized".into())
        }
    }

    pub async fn get_user_by_username(&self, username: &str) -> Result<User, Box<dyn std::error::Error + Send + Sync>> {
        self.user_repo.get_by_username(username).await
    }

    pub fn hash_password(&self, password: &str) -> Result<String, Box<dyn std::error::Error + Send + Sync>> {
        Ok(hash(password, DEFAULT_COST)?)
    }
}

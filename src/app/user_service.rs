use std::sync::Arc;
use crate::domain::repositories::UserRepository;
use crate::domain::models::User;

pub struct UserService {
    repo: Arc<dyn UserRepository>,
}

impl UserService {
    pub fn new(repo: Arc<dyn UserRepository>) -> Self {
        Self { repo }
    }

    pub async fn get_users(&self) -> Result<Vec<User>, Box<dyn std::error::Error + Send + Sync>> {
        self.repo.get_users().await
    }

    pub async fn update_user_status(&self, id: &str) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
        if id.is_empty() {
            return Err("Invalid input".into());
        }
        let mut u = self.repo.get_user(id).await?;
        u.status = !u.status;
        self.repo.update_user(u).await
    }

    pub fn get_user_count(&self) -> String {
        "1234".to_string()
    }

    pub fn get_page_view(&self) -> String {
        "1212121".to_string()
    }
}

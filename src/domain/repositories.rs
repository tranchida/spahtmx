use async_trait::async_trait;
use crate::domain::models::{User, Prize};

#[async_trait]
pub trait UserRepository: Send + Sync {
    async fn get_users(&self) -> Result<Vec<User>, Box<dyn std::error::Error + Send + Sync>>;
    async fn get_user(&self, id: &str) -> Result<User, Box<dyn std::error::Error + Send + Sync>>;
    async fn get_by_username(&self, username: &str) -> Result<User, Box<dyn std::error::Error + Send + Sync>>;
    async fn create_user(&self, user: User) -> Result<(), Box<dyn std::error::Error + Send + Sync>>;
    async fn update_user(&self, user: User) -> Result<(), Box<dyn std::error::Error + Send + Sync>>;
}

#[async_trait]
pub trait PrizeRepository: Send + Sync {
    async fn get_prizes(&self) -> Result<Vec<Prize>, Box<dyn std::error::Error + Send + Sync>>;
    async fn get_prize(&self, id: &str) -> Result<Prize, Box<dyn std::error::Error + Send + Sync>>;
    async fn get_prizes_by_year(&self, year: &str) -> Result<Vec<Prize>, Box<dyn std::error::Error + Send + Sync>>;
    async fn get_prizes_by_category(&self, category: &str) -> Result<Vec<Prize>, Box<dyn std::error::Error + Send + Sync>>;
    async fn get_prizes_by_category_and_year(&self, category: &str, year: &str) -> Result<Vec<Prize>, Box<dyn std::error::Error + Send + Sync>>;
    async fn get_categories(&self) -> Result<Vec<String>, Box<dyn std::error::Error + Send + Sync>>;
    async fn get_years(&self) -> Result<Vec<String>, Box<dyn std::error::Error + Send + Sync>>;
}

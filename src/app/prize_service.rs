use std::sync::Arc;
use crate::domain::repositories::PrizeRepository;
use crate::domain::models::Prize;

pub struct PrizeService {
    repo: Arc<dyn PrizeRepository>,
}

impl PrizeService {
    pub fn new(repo: Arc<dyn PrizeRepository>) -> Self {
        Self { repo }
    }

    pub async fn get_prizes(&self) -> Result<Vec<Prize>, Box<dyn std::error::Error + Send + Sync>> {
        self.repo.get_prizes().await
    }

    pub async fn get_prize(&self, id: &str) -> Result<Prize, Box<dyn std::error::Error + Send + Sync>> {
        self.repo.get_prize(id).await
    }

    pub async fn get_prizes_by_year(&self, year: &str) -> Result<Vec<Prize>, Box<dyn std::error::Error + Send + Sync>> {
        self.repo.get_prizes_by_year(year).await
    }

    pub async fn get_prizes_by_category(&self, category: &str) -> Result<Vec<Prize>, Box<dyn std::error::Error + Send + Sync>> {
        self.repo.get_prizes_by_category(category).await
    }

    pub async fn get_prizes_by_category_and_year(&self, category: &str, year: &str) -> Result<Vec<Prize>, Box<dyn std::error::Error + Send + Sync>> {
        self.repo.get_prizes_by_category_and_year(category, year).await
    }

    pub async fn get_categories(&self) -> Result<Vec<String>, Box<dyn std::error::Error + Send + Sync>> {
        self.repo.get_categories().await
    }

    pub async fn get_years(&self) -> Result<Vec<String>, Box<dyn std::error::Error + Send + Sync>> {
        self.repo.get_years().await
    }
}

use async_trait::async_trait;
use mongodb::{Database, bson::{doc, oid::ObjectId}};
use futures_util::stream::StreamExt;
use crate::domain::models::{Prize, Laureate};
use crate::domain::repositories::PrizeRepository;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct PrizeMongo {
    #[serde(rename = "_id")]
    pub id: ObjectId,
    pub year: String,
    pub category: String,
    #[serde(rename = "overallMotivation", skip_serializing_if = "Option::is_none")]
    pub overall_motivation: Option<String>,
    pub laureates: Vec<LaureateMongo>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct LaureateMongo {
    pub firstname: Option<String>,
    pub surname: Option<String>,
    pub motivation: Option<String>,
    pub share: Option<String>,
}

impl From<PrizeMongo> for Prize {
    fn from(p: PrizeMongo) -> Self {
        Prize {
            id: p.id.to_hex(),
            year: p.year,
            category: p.category,
            overall_motivation: p.overall_motivation,
            laureates: Some(p.laureates.into_iter().map(Into::into).collect()),
        }
    }
}

impl From<LaureateMongo> for Laureate {
    fn from(l: LaureateMongo) -> Self {
        Laureate {
            firstname: l.firstname,
            surname: l.surname,
            motivation: l.motivation,
            share: l.share,
        }
    }
}

pub struct PrizeMongoRepository {
    pub db: Database,
}

#[async_trait]
impl PrizeRepository for PrizeMongoRepository {
    async fn get_prizes(&self) -> Result<Vec<Prize>, Box<dyn std::error::Error + Send + Sync>> {
        let collection = self.db.collection::<PrizeMongo>("prize");
        let mut cursor = collection.find(doc! {}).await?;
        let mut prizes = Vec::new();
        while let Some(prize) = cursor.next().await {
            prizes.push(prize?.into());
        }
        Ok(prizes)
    }

    async fn get_prize(&self, id: &str) -> Result<Prize, Box<dyn std::error::Error + Send + Sync>> {
        let obj_id = ObjectId::parse_str(id)?;
        let collection = self.db.collection::<PrizeMongo>("prize");
        let prize = collection.find_one(doc! { "_id": obj_id }).await?;
        prize.map(Into::into).ok_or_else(|| "Prize not found".into())
    }

    async fn get_prizes_by_year(&self, year: &str) -> Result<Vec<Prize>, Box<dyn std::error::Error + Send + Sync>> {
        let collection = self.db.collection::<PrizeMongo>("prize");
        let mut cursor = collection.find(doc! { "year": year }).await?;
        let mut prizes = Vec::new();
        while let Some(prize) = cursor.next().await {
            prizes.push(prize?.into());
        }
        Ok(prizes)
    }

    async fn get_prizes_by_category(&self, category: &str) -> Result<Vec<Prize>, Box<dyn std::error::Error + Send + Sync>> {
        let collection = self.db.collection::<PrizeMongo>("prize");
        let mut cursor = collection.find(doc! { "category": category }).await?;
        let mut prizes = Vec::new();
        while let Some(prize) = cursor.next().await {
            prizes.push(prize?.into());
        }
        Ok(prizes)
    }

    async fn get_prizes_by_category_and_year(&self, category: &str, year: &str) -> Result<Vec<Prize>, Box<dyn std::error::Error + Send + Sync>> {
        let collection = self.db.collection::<PrizeMongo>("prize");
        let mut cursor = collection.find(doc! { "category": category, "year": year }).await?;
        let mut prizes = Vec::new();
        while let Some(prize) = cursor.next().await {
            prizes.push(prize?.into());
        }
        Ok(prizes)
    }

    async fn get_categories(&self) -> Result<Vec<String>, Box<dyn std::error::Error + Send + Sync>> {
        let collection = self.db.collection::<PrizeMongo>("prize");
        let categories = collection.distinct("category", doc! {}).await?;
        Ok(categories.into_iter().map(|b| b.as_str().unwrap_or_default().to_string()).collect())
    }

    async fn get_years(&self) -> Result<Vec<String>, Box<dyn std::error::Error + Send + Sync>> {
        let collection = self.db.collection::<PrizeMongo>("prize");
        let years = collection.distinct("year", doc! {}).await?;
        Ok(years.into_iter().map(|b| b.as_str().unwrap_or_default().to_string()).collect())
    }
}
